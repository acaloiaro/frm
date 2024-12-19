-- name: GetForm :one

SELECT *
FROM forms
WHERE workspace_id = @workspace_id
  AND id = @id;

-- name: GetDraft :one

SELECT *
FROM forms
WHERE workspace_id = @workspace_id
  AND id = @id
  AND status = 'draft';

-- name: ListForms :many

SELECT *
FROM forms
WHERE workspace_id = @workspace_id;

-- name: ListDrafts :many

SELECT *
FROM forms
WHERE workspace_id = @workspace_id
  AND form_id = @form_id
  AND status = 'draft';

-- name: SaveDraft :one

INSERT INTO forms (id, form_id, workspace_id, name, fields)
VALUES (coalesce(nullif(@id, 0), nextval('form_ids'))::bigint, @form_id, @workspace_id, @name, @fields) ON conflict(id) DO
UPDATE
SET updated_at = timezone('utc', now()),
    name = @name,
    fields = @fields RETURNING *;

-- name: PublishDraft :one
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
   WHERE forms.id = @id)
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
    status = 'published' RETURNING *;

-- update
 -- INSERT INTO forms(id, form_id, workspace_id, name, fields, status)
-- SELECT CASE
--            WHEN form_id IS NOT NULL THEN form_id
--            ELSE nextval('form_ids')
--        END AS id,
--        form_id,
--        workspace_id,
--        name,
--        fields,
--        'published'
-- FROM forms
-- WHERE forms.id = @id ON conflict(id) DO
--   UPDATE
--   SET updated_at = timezone('utc', now()),
--       form_id = NULL,
--       name = forms.name,
--       fields = forms.fields,
--       status = 'published' RETURNING *;
