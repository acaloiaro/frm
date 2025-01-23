 -- create the submissions pk sequence
CREATE SEQUENCE IF NOT EXISTS submission_ids START 1;

-- create form statuses enum
CREATE TYPE submission_status AS ENUM (
  'complete',
  'partial'
);

-- submissions of forms/fields to the collector
CREATE TABLE IF NOT EXISTS form_submissions (
    id BIGINT PRIMARY KEY DEFAULT nextval('submission_ids'),
    form_id BIGINT REFERENCES forms(id) ON DELETE CASCADE NOT NULL ,
    workspace_id TEXT NOT NULL,
    subject_id TEXT DEFAULT NULL,
    fields jsonb default '{}' NOT NULL,
    status submission_status default 'partial' NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);


CREATE INDEX IF NOT EXISTS workspace_idx ON form_submissions USING btree (workspace_id);
CREATE INDEX IF NOT EXISTS form_idx ON form_submissions USING btree (form_id);

COMMENT ON table form_submissions IS 'Respondants submit forms/fields to the collector as form_submissions';
COMMENT ON column form_submissions.fields IS 'all form submissions are serialized to JSON, see types.FormFieldValue for structure details';
COMMENT ON column form_submissions.subject_id IS 'identifies the subject/respondent who filled out the form';
