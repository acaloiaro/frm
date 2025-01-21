package frm

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/acaloiaro/frm/internal"
)

const (
	EventDraftCreated = "frmDraftCreated" // htmx event sent when new drafts are created
)

var ErrCannotDetermineWorkspace = errors.New("workspace cannot be determine without WorkspaceID or WorkspaceIDUrlParam")
var ErrNoInstanceAvailable = errors.New("no frm instance is available on the context")

// Frm is the primary API into frm
type Frm struct {
	BuilderMountPoint   string          // relative URL path where frm mounts the builder to your app's router
	CollectorMountPoint string          // relative URL path where frm mounts the collector to your app's router
	DBArgs              internal.DBArgs // database arguments
	WorkspaceID         string          // ID of the workspace that frm acts on behalf of
	WorkspaceIDUrlParam string          // name of the URL parameter that provides your workspace ID
}

// Args are arguments passed to Frm
type Args struct {
	BuilderMountPoint   string // path on the router to mount frm's builder
	CollectorMountPoint string // path on the router to mount frm's collector
	PostgresDisableSSL  bool   // disable ssl when connecting to postgres
	PostgresURL         string // postgres database URL
	PostgresSchema      string // postgres schema where frm stores data
	WorkspaceID         string // ID of the workspace for which frm is being initialized
	WorkspaceIDUrlParam string // named URL parameter that identifies the workspace, e.g. for route /{workspace_id}, the value would be "workspace_id"
}

// FormStatus is the status of a Form
//
// - Published forms are available to be used
//
// - Draft forms are in a draft state, yet to be published
type FormStatus = internal.FormStatus

const FormStatusPublished = internal.FormStatusPublished
const FormStatusDraft = internal.FormStatusDraft

// New initializes a new frm instance
//
// If the frm database has not yet been initialized, Init() should be called before mounting to a router
func New(args Args) (f *Frm, err error) {
	if args.WorkspaceID == "" && args.WorkspaceIDUrlParam == "" {
		return nil, ErrCannotDetermineWorkspace
	}

	f = &Frm{
		BuilderMountPoint:   strings.TrimSuffix(args.BuilderMountPoint, "/"),
		CollectorMountPoint: strings.TrimSuffix(args.CollectorMountPoint, "/"),
		WorkspaceID:         args.WorkspaceID,
		WorkspaceIDUrlParam: args.WorkspaceIDUrlParam,
		DBArgs: internal.DBArgs{
			URL:        args.PostgresURL,
			DisableSSL: args.PostgresDisableSSL,
			Schema:     args.PostgresSchema,
		},
	}
	return
}

// Init initializes the frm database if it hasn't been initialized
func (f *Frm) Init(ctx context.Context) (err error) {
	err = internal.InitializeDB(ctx, f.DBArgs)
	return
}

// GetForm retrieves forms by ID
func (f *Frm) GetForm(ctx context.Context, id int64) (form Form, err error) {
	var frm internal.Form
	frm, err = internal.Q(ctx, f.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          id,
	})
	if err != nil {
		return
	}

	form = (Form)(frm)
	return
}

type ListFormsArgs struct {
	Statuses []FormStatus
}

// ListForms lists all forms for the current workspace
func (f *Frm) ListForms(ctx context.Context, args ListFormsArgs) (forms []Form, err error) {
	var fs Forms
	statuses := []internal.FormStatus{}
	for _, s := range args.Statuses {
		statuses = append(statuses, (internal.FormStatus)(s))
	}
	fs, err = internal.Q(ctx, f.DBArgs).ListForms(ctx, internal.ListFormsParams{
		WorkspaceID: f.WorkspaceID,
		Statuses:    statuses,
	})
	if err != nil {
		return
	}

	for _, f := range fs {
		forms = append(forms, (Form)(f))
	}

	return
}

// Instance returns the frm instance from the request context
func Instance(ctx context.Context) (i *Frm, err error) {
	var ok bool
	i, ok = ctx.Value(internal.FrmContextKey).(*Frm)
	if !ok {
		return nil, ErrNoInstanceAvailable
	}
	return
}

type CreateShortCodeArgs struct {
	FormID    int64
	SubjectID string
}

// CreateShortCode creates short code for a given form and subject
func (f *Frm) CreateShortCode(ctx context.Context, args CreateShortCodeArgs) (sc ShortCode, err error) {
	var s internal.ShortCode
	s, err = internal.Q(ctx, f.DBArgs).SaveShortCode(ctx, internal.SaveShortCodeParams{
		WorkspaceID: f.WorkspaceID,
		FormID:      &args.FormID,
		ShortCode:   internal.GenCode(),
		SubjectID:   args.SubjectID,
	})
	return (ShortCode)(s), err
}

// # CollectorPath returns paths to frm collector endpoints
//
// It uses the collector's mount point on the router to generate collector paths
func CollectorPath(ctx context.Context, path string) string {
	base, ok := ctx.Value(internal.CollectorMountPointContextKey).(string)
	if !ok {
		return "/"
	}
	urlPath := filepath.Clean(fmt.Sprintf("%s/%s", base, path))
	return urlPath
}

// BuilderPathForm returns the builder URL path for the provided form ID
func BuilderPathForm(ctx context.Context, formID int64) string {
	base, ok := ctx.Value(internal.BuilderMountPointContextKey).(string)
	if !ok {
		return "/"
	}

	return fmt.Sprintf("%s/%d", base, formID)
}

// BuilderPathFormField returns the builder URL path for the provided form ID and field ID
func BuilderPathFormField(ctx context.Context, formID int64, fieldID string, args ...string) string {
	base, ok := ctx.Value(internal.BuilderMountPointContextKey).(string)
	if !ok {
		return "/"
	}

	additionalPath := ""
	if len(args) > 0 {
		additionalPath = args[0]
	}

	path := filepath.Clean(fmt.Sprintf("%s/%d/fields/%s/%s", base, formID, fieldID, additionalPath))
	return path
}

// CollectorPathForm returns the collector URL path for the provided form ID
func CollectorPathForm(ctx context.Context, formID int64, path ...string) string {
	base, ok := ctx.Value(internal.CollectorMountPointContextKey).(string)
	if !ok {
		return "/"
	}
	base = filepath.Clean(base)
	return fmt.Sprintf("%s/%d", base, formID)
}
