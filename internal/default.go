package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/acaloiaro/frm/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	// MountPointContextKey is the context key representing frm's mount point on the request context
	MountPointContextKey contextKey = "mount_point_context_key"
	// FrmContextKey is the context key representing the frm instance on the request context
	FrmContextKey contextKey = "frm_instance"
)

var pool *pgxpool.Pool

type Forms []Form

// getPool returns a database pool for the specified connection string
func getPool(ctx context.Context, databaseURL string) (p *pgxpool.Pool, err error) {
	if pool == nil {
		var poolConfig *pgxpool.Config
		poolConfig, err = pgxpool.ParseConfig(databaseURL)
		if err != nil {
			err = fmt.Errorf("invalid connection string: %v", err)
			return
		}

		err = InitializeDB(ctx, databaseURL)
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
func Q(ctx context.Context, postgresURL string) *Queries {
	p, err := getPool(ctx, postgresURL)
	if err != nil {
		slog.Error("database pool unavailable", "error", err)
		return nil // TODO return a no-op DBTX to avoid NPEs
	}
	return New(p)
}

// InitializeDB creates the application database if it doesn't exist and runs all migrations against it
//
// databaseUrl is the database URL string to connect to the database
func InitializeDB(ctx context.Context, postgresURL string) (err error) {
	postgresURL, err = pgConnectionString(postgresURL, true)
	if err != nil {
		return fmt.Errorf("invalid connection string: %v", err)
	}

	var cfg *pgx.ConnConfig
	cfg, err = pgx.ParseConfig(postgresURL)
	if err != nil {
		return
	}

	err = createIfNotExist(ctx, cfg)
	if err != nil {
		return
	}

	err = runMigrations(postgresURL)
	return
}

// TODO not all db crdentials have permission to do this. Fail gracefully when user lacks permission
func createIfNotExist(ctx context.Context, cfg *pgx.ConnConfig) (err error) {
	dbName := cfg.Database
	cfg.Database = ""
	cfg.RuntimeParams = nil

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return
	}

	var exists int
	err = conn.QueryRow(ctx, "SELECT 1 FROM pg_database WHERE datname=$1", dbName).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("unable to create database: %w", err)
	}

	if exists != 1 {
		// note: prepared statements are not supported by pgx for CREATE DATABASE queries
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("unable to create database: %w", err)
		}
	}
	return
}

// runMigrations runs all available migrations if steps are unspecified ( 0 ), and runs either up steps or down steps
func runMigrations(postgresURL string) (err error) {
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
		err = fmt.Errorf("unable to run migrations!!!: %v", err)
	} else {
		err = nil
	}

	return
}

func pgConnectionString(postgresURL string, disableSSL bool, options ...string) (connString string, err error) {
	var cfg *pgx.ConnConfig
	cfg, err = pgx.ParseConfig(postgresURL)
	if err != nil {
		return
	}

	options = append(options, "x-migrations-table=frm_migrations")

	if disableSSL {
		options = append(options, "sslmode=disable")
	}

	connString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, strings.Join(options, "&"))
	return
}
