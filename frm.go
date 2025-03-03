package frm

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/acaloiaro/frm/internal"
)

const (
	EventDraftCreated     = "frmDraftCreated" // htmx event sent when new drafts are created
	EventCloneCreated     = "frmCloneCreated" // htmx event sent when new clones are created
	DefaultCopyNameSuffix = "(COPY)"          // the default suffix added to forms when they're copied
)

var ErrCannotDetermineWorkspace = errors.New("workspace cannot be determine without WorkspaceID or WorkspaceIDUrlParam")
var ErrNoInstanceAvailable = errors.New("no frm instance is available on the context")

// Frm is the primary API into frm
type Frm struct {
	BuilderMountPoint   string                 // relative URL path where frm mounts the builder to your app's router
	CollectorMountPoint string                 // relative URL path where frm mounts the collector to your app's router
	CollectorFooter     string                 // footer shown at the bottom of the collector page
	DraftMaxAge         time.Duration          // the duration that form drafts may remain in the draft stage before removal
	DBArgs              internal.DBArgs        // database arguments
	Receiver            FormSubmissionReceiver // function that processes incoming form submissions
	WorkspaceID         string                 // ID of the workspace that frm acts on behalf of
	WorkspaceIDUrlParam string                 // name of the URL parameter that provides your workspace ID
}

// Args are arguments passed to Frm
type Args struct {
	BuilderMountPoint   string                 // path on the router to mount frm's builder
	CollectorMountPoint string                 // path on the router to mount frm's collector
	CollectorFooter     string                 // footer shown at the bottom of the collector page
	DraftMaxAge         time.Duration          // the duration that form drafts may remain in the draft state before removal
	PostgresDisableSSL  bool                   // disable ssl when connecting to postgres
	PostgresSchema      string                 // postgres schema where frm stores data
	PostgresURL         string                 // postgres database URL
	Reciever            FormSubmissionReceiver // function that processes incoming form submissions
	WorkspaceID         string                 // ID of the workspace for which frm is being initialized
	WorkspaceIDUrlParam string                 // named URL parameter that identifies the workspace, e.g. for route /{workspace_id}, the value would be "workspace_id"
}

// FormSubmissionReceiver processes form submissions
type FormSubmissionReceiver = func(ctx context.Context, submission FormSubmission) (err error)

// FormStatus is the status of a Form
//
// - Published forms are available to be used
//
// - Draft forms are in a draft state, yet to be published
//
// - Archived forms are not intended to be used
type FormStatus = internal.FormStatus

const FormStatusPublished = internal.FormStatusPublished
const FormStatusDraft = internal.FormStatusDraft
const FormStatusArchived = internal.FormStatusArchived

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
		CollectorFooter:     args.CollectorFooter,
		DraftMaxAge:         args.DraftMaxAge,
		DBArgs: internal.DBArgs{
			URL:        args.PostgresURL,
			DisableSSL: args.PostgresDisableSSL,
			Schema:     args.PostgresSchema,
		},
		Receiver:            args.Reciever,
		WorkspaceID:         args.WorkspaceID,
		WorkspaceIDUrlParam: args.WorkspaceIDUrlParam,
	}
	return
}

// Init initializes the frm database if it hasn't been initialized
func (f *Frm) Init(ctx context.Context) (err error) {
	err = internal.InitializeDB(ctx, f.DBArgs)
	if err != nil {
		return
	}

	go func() {
		err = internal.DraftMonitor(ctx, f.DBArgs, f.DraftMaxAge)
		if err != nil {
			return
		}
	}()
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

// CopyFormArgs are passed to frm.CopyForm()
type CopyFormArgs struct {
	ForgetParentForm bool   // forget the parent form from which the copied form is copied
	ID               int64  // the id of the form to copy
	NameSuffix       string // suffix to be added to the original form's name, e.g. "(COPY)"
}

// CopyForm copies existing Forms
//
// Returns the copied form
func (f *Frm) CopyForm(ctx context.Context, args CopyFormArgs) (form Form, err error) {
	var of internal.Form
	of, err = internal.Q(ctx, f.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          args.ID,
	})
	if err != nil {
		return
	}
	copiedFormName := of.Name
	if args.NameSuffix != "" {
		copiedFormName = fmt.Sprintf("%s %s", of.Name, args.NameSuffix)
	}
	p := &internal.SaveFormParams{
		WorkspaceID: of.WorkspaceID,
		FormID:      &of.ID,
		Name:        copiedFormName,
		Fields:      of.Fields,
		Status:      FormStatusDraft,
	}
	if args.ForgetParentForm {
		p.FormID = nil
	}
	nf, err := internal.Q(ctx, f.DBArgs).SaveForm(ctx, *p)
	if err != nil {
		return
	}

	form = (Form)(nf)
	return
}

// GetFormSubmission retrieves form submissions by ID
func (f *Frm) GetFormSubmission(ctx context.Context, submissionID int64) (sub FormSubmission, err error) {
	var s internal.FormSubmission
	s, err = internal.Q(ctx, f.DBArgs).GetFormSubmission(ctx, internal.GetFormSubmissionParams{
		WorkspaceID:  f.WorkspaceID,
		SubmissionID: submissionID,
	})
	if err != nil {
		return
	}

	sub = (FormSubmission)(s)
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

// Instance returns the frm instance from the request context (if available)
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
		ShortCode:   internal.GenShortCode(),
		SubjectID:   args.SubjectID,
	})
	return (ShortCode)(s), err
}

// BuilderPath returns paths to frm builder endpoints
//
// It uses the builer's mount point on the router to generate builder paths
func BuilderPath(ctx context.Context, path string) string {
	base, ok := ctx.Value(internal.BuilderMountPointContextKey).(string)
	if !ok {
		return "/"
	}

	return fmt.Sprintf("%s/%s", base, path)
}

// CollectorPath returns paths to frm collector endpoints
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

// CollectorPathShortCode returns the collector's URL path for the provided shortcode
func CollectorPathShortCode(ctx context.Context, shortCode string) string {
	base, ok := ctx.Value(internal.CollectorMountPointContextKey).(string)
	if !ok {
		return "/"
	}
	base = filepath.Clean(base)
	return fmt.Sprintf("%s/s/%s", base, shortCode)
}
