-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- create the forms pk sequence
CREATE SEQUENCE form_ids START 1;

-- the forms table contains all forms
-- workspaces are to separate form ownership
CREATE TABLE IF NOT EXISTS forms (
    id BIGINT PRIMARY KEY DEFAULT nextval('form_ids'),
    workspace_id UUID NOT NULL DEFAULT uuid_generate_v4(), -- // workspace_id is used no namespace the form
    name text not null,
    fields jsonb default '{}' NOT NULL,
    created_at timestamptz not null default timezone('utc', now()),
    updated_at timestamptz not null default timezone('utc', now())
);

COMMENT ON table forms IS 'Form contains all the data necesary to render a form';
COMMENT ON column forms.workspace_id IS 'a namespace for the form';
COMMENT ON column forms.fields IS 'all form fields are serialized to JSON, see types.FormFields for structure details';
