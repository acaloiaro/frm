-- name: ListForm :many

SELECT *
FROM forms;

-- name: GetForm :one

SELECT *
FROM forms
WHERE id = @id;

-- name: SaveForm :one

INSERT INTO forms (id, name, fields)
VALUES (coalesce(nullif(@id, 0), nextval('form_ids'))::bigint, @name, @fields) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    name = @name,
    fields = @fields RETURNING *;
