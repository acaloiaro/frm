package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/acaloiaro/frm/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoopDatabase = errors.New("frm had a problem getting a database connection")
)

type MockRow struct{}
type NoopDBTX struct{}

func (m *MockRow) Scan(dest ...any) (err error) {
	return ErrNoopDatabase
}

func (n *NoopDBTX) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, ErrNoopDatabase
}

func (n *NoopDBTX) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, ErrNoopDatabase
}
func (n *NoopDBTX) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return &MockRow{}
}

type contextKey string

const (
	// BuilderMountPointContextKey is the context key representing frm's builder mount point on the request context
	BuilderMountPointContextKey contextKey = "builder_mount_point_context_key"
	// DefaultShortcodeLen is the default length for generating short codes
	DefaultShortcodeLen = 6
	// CollectorMountPointContextKey is the context key representing frm's collector mount point on the request context
	CollectorMountPointContextKey contextKey = "collector_mount_point_context_key"
	// FrmContextKey is the context key representing the frm instance on the request context
	FrmContextKey contextKey = "frm_instance"
)

var (
	FormIDContextKey    contextKey = "frm_form_id"
	FieldIDContextKey   contextKey = "frm_field_id"
	ShortCodeContextKey contextKey = "frm_short_code"
	pool                *pgxpool.Pool
	shortcodeCharset    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type Forms []Form

var enumTypes = []string{
	"form_status",
	"_form_status", // array of form statuses
	"submission_status",
	"_submission_status",
}

// getPool returns a database pool for the specified connection string
func getPool(ctx context.Context, args DBArgs) (p *pgxpool.Pool, err error) {
	if pool == nil {
		var poolConfig *pgxpool.Config
		var postgresURL string
		postgresURL, err = pgConnectionString(args)
		if err != nil {
			return
		}
		poolConfig, err = pgxpool.ParseConfig(postgresURL)
		if err != nil {
			err = fmt.Errorf("invalid connection string: %v", err)
			return
		}

		// this after conect hook allows pgx to correclty encode enum types as query params
		// reference: https://github.com/jackc/pgx/issues/1549#issuecomment-1467107173
		poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			for _, typ := range enumTypes {
				t, err := conn.LoadType(ctx, typ)
				if err != nil {
					return err
				}
				conn.TypeMap().RegisterType(t)
			}
			return nil
		}

		err = InitializeDB(ctx, args)
		if err != nil {
			err = fmt.Errorf("database failed to initialize: %v", err)
			return
		}

		p, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			err = fmt.Errorf("invalid connection string: %v", err)
			return
		}
	} else {
		p = pool
	}
	return
}

// Q returns a DBTX instance for querying the frm database
func Q(ctx context.Context, args DBArgs) *Queries {
	p, err := getPool(ctx, args)
	if err != nil {
		slog.Error("unable to get database connect", "error", err)
		return New(&NoopDBTX{})
	}
	return New(p)
}

type DBArgs struct {
	URL        string
	DisableSSL bool
	Schema     string
}

// InitializeDB creates the application database if it doesn't exist and runs all migrations against it
//
// databaseUrl is the database URL string to connect to the database
func InitializeDB(ctx context.Context, args DBArgs) (err error) {
	err = createIfNotExist(ctx, args)
	if err != nil {
		return
	}

	err = runMigrations(args)
	return
}

// TODO not all db credentials have permission create databases/schemas. Fail gracefully when user lacks permission
func createIfNotExist(ctx context.Context, args DBArgs) (err error) {
	var cfg *pgx.ConnConfig
	cfg, err = pgx.ParseConfig(args.URL)
	if err != nil {
		return
	}

	dbName := cfg.Database
	schemaName := "frm"
	if args.Schema != "" {
		schemaName = args.Schema
	}

	cfg.Database = "postgres"
	cfg.RuntimeParams = nil

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return
	}

	var dbExists int
	err = conn.QueryRow(ctx, "SELECT 1 FROM pg_database WHERE datname=$1", dbName).Scan(&dbExists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("unable to create database: %w", err)
	}

	if dbExists != 1 {
		// note: prepared statements are not supported by pgx for CREATE DATABASE queries
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("unable to create database: %w", err)
		}
	}

	cfg.Database = dbName
	conn, err = pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return
	}
	var schemaExists int
	err = conn.QueryRow(ctx, "SELECT 1 FROM information_schema.schemata where schema_name=$1", schemaName).Scan(&schemaExists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("unable to create schema: %w", err)
	}

	if schemaExists != 1 {
		// note: prepared statements are not supported by pgx for CREATE DATABASE queries
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE schema %s", schemaName))
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("unable to create schema: %w", err)
		}
	}

	return
}

// runMigrations runs all available migrations if steps are unspecified ( 0 ), and runs either up steps or down steps
func runMigrations(args DBArgs) (err error) {
	postgresURL, err := pgConnectionString(args)
	if err != nil {
		return fmt.Errorf("invalid connection string: %v", err)
	}
	ms, err := iofs.New(migrations.FS, ".")
	if err != nil {
		panic(fmt.Sprintf("unable to run migrations: %v", err))
	}

	m, err := migrate.NewWithSourceInstance("iofs", ms, postgresURL)
	if err != nil {
		panic(fmt.Sprintf("unable to run migrations: %v", err))
	}

	// We don't need the migration tooling to hold it's connections to the DB once it has been completed.
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		err = fmt.Errorf("unable to run frm migrations: %v", err)
	} else {
		err = nil
	}

	return
}

func pgConnectionString(args DBArgs) (connString string, err error) {
	var cfg *pgx.ConnConfig
	cfg, err = pgx.ParseConfig(args.URL)
	if err != nil {
		return
	}
	options := []string{}
	if args.DisableSSL {
		options = append(options, "sslmode=disable")
	}

	if args.Schema != "" {
		options = append(options, fmt.Sprintf("search_path=%s", args.Schema))
	} else {
		options = append(options, "search_path=frm")
	}

	connString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, strings.Join(options, "&"))
	return
}

// JSON returns the form's JSON-seralized string representation
func (f Form) JSON() string {
	b, err := json.Marshal(f)
	if err != nil {
		return ""
	}

	return string(b)
}

// GenShortCode generates new shortcodes
func GenShortCode() string {
	b := make([]rune, DefaultShortcodeLen)
	chsLen := len(shortcodeCharset)
	for i := range b {
		b[i] = shortcodeCharset[rand.Intn(chsLen)]
	}
	return string(b)
}

// FormSubmissionMap converts a FormSubmission to its map[string]any representation (for queueing webhooks)
func FormSubmissionMap(s FormSubmission) (m map[string]any) {
	m = map[string]any{}
	m["id"] = s.ID
	m["form_id"] = s.FormID
	m["workspace_id"] = s.WorkspaceID
	if s.SubjectID != nil {
		m["subject_id"] = *s.SubjectID
	}
	m["fields"] = s.Fields
	m["created_at"] = s.CreatedAt
	m["updated_at"] = s.UpdatedAt
	return
}
