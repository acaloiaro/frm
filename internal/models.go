// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package internal

import (
	"time"

	"github.com/acaloiaro/frm/types"
	uuid "github.com/google/uuid"
)

// Form contains all the data necesary to render a form
type Form struct {
	ID int64 `json:"id"`
	// a namespace for the form
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	// all form fields are serialized to JSON, see types.FormFields for structure details
	Fields    types.FormFields `json:"fields"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
