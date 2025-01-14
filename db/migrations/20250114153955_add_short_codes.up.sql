-- create the short_codes pk sequence
CREATE SEQUENCE IF NOT EXISTS short_code_ids START 1;

-- short_codes is a mapping of short codes to their corresponding forms, along with an optional subject id identifying
-- the subject that the short code belongs to 
CREATE TABLE IF NOT EXISTS short_codes (
    id BIGINT PRIMARY KEY DEFAULT nextval('short_code_ids'),
    workspace_id TEXT NOT NULL,
    form_id BIGINT REFERENCES forms(id) ON DELETE CASCADE,
    short_code TEXT NOT NULL,
    subject_id TEXT NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);

CREATE INDEX IF NOT EXISTS workspace_short_code_idx ON short_codes USING btree (workspace_id,short_code);
CREATE INDEX IF NOT EXISTS form_idx ON short_codes USING btree (form_id);
CREATE INDEX IF NOT EXISTS subject_idx ON short_codes USING btree (subject_id);
CREATE UNIQUE INDEX IF NOT EXISTS subject_form_idx ON short_codes USING btree (subject_id, form_id);

COMMENT ON table short_codes IS 'Short codes are short codes/names for URLs that identify the subject submitting a form'
