-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- create the forms pk sequence
CREATE SEQUENCE IF NOT EXISTS form_ids START 1;
CREATE SEQUENCE IF NOT EXISTS draft_form_ids START 1;

-- the forms table contains all forms
-- workspaces represent users, tenants, groupings, etc. in external systems
CREATE TABLE IF NOT EXISTS forms (
    id BIGINT PRIMARY KEY DEFAULT nextval('form_ids'),
    workspace_id UUID NOT NULL,
    name text not null,
    fields jsonb default '{}' NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);

CREATE TABLE IF NOT EXISTS draft_forms (
    id BIGINT PRIMARY KEY DEFAULT nextval('draft_form_ids'),
    form_id BIGINT REFERENCES forms(id),
    workspace_id UUID,
    name text not null,
    fields jsonb default '{}' NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);

COMMENT ON table forms IS 'Form contains all the data necesary to render a form';
COMMENT ON table draft_forms IS 'Draft forms allow forms to be edited in situ and act as a staging ground for Forms';
COMMENT ON column forms.workspace_id IS 'a namespace for the form';
COMMENT ON column forms.fields IS 'all form fields are serialized to JSON, see types.FormFields for structure details';
