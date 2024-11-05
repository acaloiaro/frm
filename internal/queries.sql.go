// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package internal

import (
	"context"

	"github.com/acaloiaro/frm/types"
)

const getForm = `-- name: GetForm :one

SELECT id, workspace_id, name, fields, created_at, updated_at
FROM forms
WHERE id = $1
`

// GetForm
//
//	SELECT id, workspace_id, name, fields, created_at, updated_at
//	FROM forms
//	WHERE id = $1
func (q *Queries) GetForm(ctx context.Context, id int64) (Form, error) {
	row := q.db.QueryRow(ctx, getForm, id)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listForm = `-- name: ListForm :many

SELECT id, workspace_id, name, fields, created_at, updated_at
FROM forms
`

// ListForm
//
//	SELECT id, workspace_id, name, fields, created_at, updated_at
//	FROM forms
func (q *Queries) ListForm(ctx context.Context) ([]Form, error) {
	rows, err := q.db.Query(ctx, listForm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Form
	for rows.Next() {
		var i Form
		if err := rows.Scan(
			&i.ID,
			&i.WorkspaceID,
			&i.Name,
			&i.Fields,
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

const saveForm = `-- name: SaveForm :one

INSERT INTO forms (id, name, fields)
VALUES (coalesce(nullif($1, 0), nextval('form_ids'))::bigint, $2, $3) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    name = $2,
    fields = $3 RETURNING id, workspace_id, name, fields, created_at, updated_at
`

type SaveFormParams struct {
	ID     interface{}      `json:"id"`
	Name   string           `json:"name"`
	Fields types.FormFields `json:"fields"`
}

// SaveForm
//
//	INSERT INTO forms (id, name, fields)
//	VALUES (coalesce(nullif($1, 0), nextval('form_ids'))::bigint, $2, $3) ON conflict(id) DO
//	UPDATE
//	SET updated_at = timezone('utc', now()),
//	    name = $2,
//	    fields = $3 RETURNING id, workspace_id, name, fields, created_at, updated_at
func (q *Queries) SaveForm(ctx context.Context, arg SaveFormParams) (Form, error) {
	row := q.db.QueryRow(ctx, saveForm, arg.ID, arg.Name, arg.Fields)
	var i Form
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Name,
		&i.Fields,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
