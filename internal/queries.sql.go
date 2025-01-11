// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package internal

import (
	"context"

	"github.com/acaloiaro/frm/types"
	uuid "github.com/google/uuid"
)

const deleteForm = `-- name: DeleteForm :exec

DELETE
FROM forms
WHERE workspace_id = $1
  AND id = $2
`

type DeleteFormParams struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ID          int64     `json:"id"`
}

// DeleteForm
//
//	DELETE
//	FROM forms
//	WHERE workspace_id = $1
//	  AND id = $2
func (q *Queries) DeleteForm(ctx context.Context, arg DeleteFormParams) error {
	_, err := q.db.Exec(ctx, deleteForm, arg.WorkspaceID, arg.ID)
	return err
}

const getDraft = `-- name: GetDraft :one

SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
FROM forms
WHERE workspace_id = $1
  AND id = $2
  AND status = 'draft'
`

type GetDraftParams struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ID          int64     `json:"id"`
}

// GetDraft
//
//	SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
//	FROM forms
//	WHERE workspace_id = $1
//	  AND id = $2
//	  AND status = 'draft'
func (q *Queries) GetDraft(ctx context.Context, arg GetDraftParams) (Form, error) {
	row := q.db.QueryRow(ctx, getDraft, arg.WorkspaceID, arg.ID)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.FormID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getForm = `-- name: GetForm :one

SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
FROM forms
WHERE workspace_id = $1
  AND id = $2
`

type GetFormParams struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ID          int64     `json:"id"`
}

// GetForm
//
//	SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
//	FROM forms
//	WHERE workspace_id = $1
//	  AND id = $2
func (q *Queries) GetForm(ctx context.Context, arg GetFormParams) (Form, error) {
	row := q.db.QueryRow(ctx, getForm, arg.WorkspaceID, arg.ID)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.FormID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listDrafts = `-- name: ListDrafts :many

SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
FROM forms
WHERE workspace_id = $1
  AND form_id = $2
  AND status = 'draft'
`

type ListDraftsParams struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	FormID      *int64    `json:"form_id"`
}

// ListDrafts
//
//	SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
//	FROM forms
//	WHERE workspace_id = $1
//	  AND form_id = $2
//	  AND status = 'draft'
func (q *Queries) ListDrafts(ctx context.Context, arg ListDraftsParams) ([]Form, error) {
	rows, err := q.db.Query(ctx, listDrafts, arg.WorkspaceID, arg.FormID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Form
	for rows.Next() {
		var i Form
		if err := rows.Scan(
			&i.ID,
			&i.FormID,
			&i.WorkspaceID,
			&i.Name,
			&i.Fields,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listForms = `-- name: ListForms :many

SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
FROM forms
WHERE workspace_id = $1
  AND status = any(CASE
                       WHEN cardinality($2::form_status[]) > 0 THEN $2::form_status[]
                       ELSE enum_range(NULL::form_status)::form_status[]
                   END::form_status[])
`

type ListFormsParams struct {
	WorkspaceID uuid.UUID    `json:"workspace_id"`
	Statuses    []FormStatus `json:"statuses"`
}

// ListForms
//
//	SELECT id, form_id, workspace_id, name, fields, status, created_at, updated_at
//	FROM forms
//	WHERE workspace_id = $1
//	  AND status = any(CASE
//	                       WHEN cardinality($2::form_status[]) > 0 THEN $2::form_status[]
//	                       ELSE enum_range(NULL::form_status)::form_status[]
//	                   END::form_status[])
func (q *Queries) ListForms(ctx context.Context, arg ListFormsParams) ([]Form, error) {
	rows, err := q.db.Query(ctx, listForms, arg.WorkspaceID, arg.Statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Form
	for rows.Next() {
		var i Form
		if err := rows.Scan(
			&i.ID,
			&i.FormID,
			&i.WorkspaceID,
			&i.Name,
			&i.Fields,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const publishDraft = `-- name: PublishDraft :one
WITH draft AS
  (SELECT CASE
              WHEN form_id IS NOT NULL THEN form_id
              ELSE nextval('form_ids')
          END AS id,
          form_id,
          workspace_id,
          name,
          fields,
          'published'
   FROM forms
   WHERE forms.id = $1)
INSERT INTO forms(id, form_id, workspace_id, name, fields, status)
VALUES ((SELECT id FROM draft), NULL, (SELECT workspace_id FROM draft), (SELECT name FROM draft), (SELECT fields FROM draft), 'published') ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    form_id = NULL,
    workspace_id =
  (SELECT workspace_id
   FROM draft),
    name =
  (SELECT name
   FROM draft),
    fields =
  (SELECT fields
   FROM draft),
    status = 'published' RETURNING id, form_id, workspace_id, name, fields, status, created_at, updated_at
`

// PublishDraft
//
//	WITH draft AS
//	  (SELECT CASE
//	              WHEN form_id IS NOT NULL THEN form_id
//	              ELSE nextval('form_ids')
//	          END AS id,
//	          form_id,
//	          workspace_id,
//	          name,
//	          fields,
//	          'published'
//	   FROM forms
//	   WHERE forms.id = $1)
//	INSERT INTO forms(id, form_id, workspace_id, name, fields, status)
//	VALUES ((SELECT id FROM draft), NULL, (SELECT workspace_id FROM draft), (SELECT name FROM draft), (SELECT fields FROM draft), 'published') ON conflict(id) DO
//	UPDATE
//	SET updated_at = timezone('utc', now()),
//	    form_id = NULL,
//	    workspace_id =
//	  (SELECT workspace_id
//	   FROM draft),
//	    name =
//	  (SELECT name
//	   FROM draft),
//	    fields =
//	  (SELECT fields
//	   FROM draft),
//	    status = 'published' RETURNING id, form_id, workspace_id, name, fields, status, created_at, updated_at
func (q *Queries) PublishDraft(ctx context.Context, id int64) (Form, error) {
	row := q.db.QueryRow(ctx, publishDraft, id)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.FormID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const saveDraft = `-- name: SaveDraft :one

INSERT INTO forms (id, form_id, workspace_id, name, fields)
VALUES (coalesce(nullif($1, 0), nextval('form_ids'))::bigint, $2, $3, $4, $5) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    name = $4,
    fields = $5 RETURNING id, form_id, workspace_id, name, fields, status, created_at, updated_at
`

type SaveDraftParams struct {
	ID          interface{}      `json:"id"`
	FormID      *int64           `json:"form_id"`
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	Name        string           `json:"name"`
	Fields      types.FormFields `json:"fields"`
}

// SaveDraft
//
//	INSERT INTO forms (id, form_id, workspace_id, name, fields)
//	VALUES (coalesce(nullif($1, 0), nextval('form_ids'))::bigint, $2, $3, $4, $5) ON conflict(id) DO
//	UPDATE
//	SET updated_at = timezone('utc', now()),
//	    name = $4,
//	    fields = $5 RETURNING id, form_id, workspace_id, name, fields, status, created_at, updated_at
func (q *Queries) SaveDraft(ctx context.Context, arg SaveDraftParams) (Form, error) {
	row := q.db.QueryRow(ctx, saveDraft,
		arg.ID,
		arg.FormID,
		arg.WorkspaceID,
		arg.Name,
		arg.Fields,
	)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.FormID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const saveSubmission = `-- name: SaveSubmission :one

INSERT INTO form_submissions (id, form_id, workspace_id, fields, status)
VALUES (coalesce(nullif($1, 0), nextval('submission_ids'))::bigint, $2, $3, $4, $5) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    fields = $4,
    status = $5 RETURNING id, form_id, workspace_id, fields, status, created_at, updated_at
`

type SaveSubmissionParams struct {
	ID          interface{}           `json:"id"`
	FormID      *int64                `json:"form_id"`
	WorkspaceID uuid.UUID             `json:"workspace_id"`
	Fields      types.FormFieldValues `json:"fields"`
	Status      SubmissionStatus      `json:"status"`
}

// SaveSubmission
//
//	INSERT INTO form_submissions (id, form_id, workspace_id, fields, status)
//	VALUES (coalesce(nullif($1, 0), nextval('submission_ids'))::bigint, $2, $3, $4, $5) ON conflict(id) DO
//	UPDATE
//	SET updated_at = timezone('utc', now()),
//	    fields = $4,
//	    status = $5 RETURNING id, form_id, workspace_id, fields, status, created_at, updated_at
func (q *Queries) SaveSubmission(ctx context.Context, arg SaveSubmissionParams) (FormSubmission, error) {
	row := q.db.QueryRow(ctx, saveSubmission,
		arg.ID,
		arg.FormID,
		arg.WorkspaceID,
		arg.Fields,
		arg.Status,
	)
	var i FormSubmission
	err := row.Scan(
		&i.ID,
		&i.FormID,
		&i.WorkspaceID,
		&i.Fields,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
