-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- create the forms pk sequence
CREATE SEQUENCE IF NOT EXISTS form_ids START 1;

-- create form statuses enum
CREATE TYPE form_status AS ENUM (
  'published',
  'draft'
);

-- the forms table contains all forms
-- workspaces represent users, tenants, groupings, etc. in external systems
CREATE TABLE IF NOT EXISTS forms (
    id BIGINT PRIMARY KEY DEFAULT nextval('form_ids'),
    form_id BIGINT REFERENCES forms(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL,
    name text not null,
    fields jsonb default '{}' NOT NULL,
    status form_status default 'draft' NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);


COMMENT ON table forms IS 'Form contains all the data necesary to render a form';
COMMENT ON column forms.workspace_id IS 'a namespace for the form';
COMMENT ON column forms.fields IS 'all form fields are serialized to JSON, see types.FormFields for structure details';
