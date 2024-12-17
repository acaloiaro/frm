-- name: GetForm :one

SELECT *
FROM forms
WHERE workspace_id = @workspace_id
  AND id = @id;

-- name: ListForms :many

SELECT *
FROM forms
WHERE workspace_id = @workspace_id;

-- name: SaveForm :one

INSERT INTO forms (id, workspace_id, name, fields)
VALUES (coalesce(nullif(@id, 0), nextval('form_ids'))::bigint, @workspace_id, @name, @fields) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    name = @name,
    fields = @fields RETURNING *;
