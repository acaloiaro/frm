package frm

import (
	"context"

	"github.com/acaloiaro/frm/internal"
	"github.com/google/uuid"
)

type Frm struct {
	PostgresURL string    // the database URL where forms are stored
	WorkspaceID uuid.UUID // the ID of the workspace that this instance acts on behalf of
}

type Args struct {
	PostgresURL string
	WorkspaceID uuid.UUID
}

// New initializes a new frm instance
//
// If the frm database hasn't been initiailized, the database is initialized
func New(args Args) *Frm {
	return &Frm{
		PostgresURL: args.PostgresURL,
		WorkspaceID: args.WorkspaceID,
	}
}

// Init ensures that the frm database is initialized
func (f *Frm) Init(ctx context.Context) (err error) {
	err = internal.InitializeDB(ctx, f.PostgresURL)
	return
}

// GetForm retrieves a form by ID
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
