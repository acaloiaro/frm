// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package internal

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/acaloiaro/frm/types"
	uuid "github.com/google/uuid"
)

type FormStatus string

const (
	FormStatusPublished FormStatus = "published"
	FormStatusDraft     FormStatus = "draft"
)

func (e *FormStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FormStatus(s)
	case string:
		*e = FormStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for FormStatus: %T", src)
	}
	return nil
}

type NullFormStatus struct {
	FormStatus FormStatus `json:"form_status"`
	Valid      bool       `json:"valid"` // Valid is true if FormStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFormStatus) Scan(value interface{}) error {
	if value == nil {
		ns.FormStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FormStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFormStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FormStatus), nil
}

// Form contains all the data necesary to render a form
type Form struct {
	ID     int64  `json:"id"`
	FormID *int64 `json:"form_id"`
	// a namespace for the form
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	// all form fields are serialized to JSON, see types.FormFields for structure details
	Fields    types.FormFields `json:"fields"`
	Status    FormStatus       `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
