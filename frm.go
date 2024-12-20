package frm

import (
	"context"
	"errors"

	"github.com/acaloiaro/frm/internal"
	"github.com/google/uuid"
)

const (
	EventDraftCreated = "frmDraftCreated" // htmx event sent when new drafts are created
)

var ErrCannotDetermineWorkspace = errors.New("workspace cannot be determine without WorkspaceID or WorkspaceIDUrlParam")

// Frm is the primary API into frm
type Frm struct {
	PostgresURL         string    // the database URL where forms are stored
	WorkspaceID         uuid.UUID // the ID of the workspace that the frm acts on behalf of
	WorkspaceIDUrlParam string    // the name of the URL parameter that provides your workspace ID
}

// Args are arguments passed to Frm
type Args struct {
	PostgresURL         string
	WorkspaceID         uuid.UUID
	WorkspaceIDUrlParam string
}

type FormStatus internal.FormStatus

// New initializes a new frm instance
//
// If the frm database hasn't been initiailized, the database is initialized
func New(args Args) (f *Frm, err error) {
	if args.WorkspaceID == uuid.Nil && args.WorkspaceIDUrlParam == "" {
		return nil, ErrCannotDetermineWorkspace
	}

	f = &Frm{
		PostgresURL:         args.PostgresURL,
		WorkspaceID:         args.WorkspaceID,
		WorkspaceIDUrlParam: args.WorkspaceIDUrlParam,
	}
	return
}

// Init initializes the frm database if it hasn't been initialized
func (f *Frm) Init(ctx context.Context) (err error) {
	err = internal.InitializeDB(ctx, f.PostgresURL)
	return
}

// GetForm retrieves forms by ID
func (f *Frm) GetForm(ctx context.Context, id int64) (form Form, err error) {
	var frm internal.Form
	frm, err = internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
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
	fs, err = internal.Q(ctx, f.PostgresURL).ListForms(ctx, internal.ListFormsParams{
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
