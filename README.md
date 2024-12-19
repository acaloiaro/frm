# frm 

frm is an embeddable HTML form builder for Go. Its goal is to be embeddable within any `net/http`-comptabile http router.

**This is a work in progress and not production-ready**

## Concepts

### Workspaces 

frm's design is oriented around the concept of "workspaces". Workspaces may represent users or tenants within your application. As such, every interaction with frm must be workspace-aware.

### Draft forms

"Drafts" are created to edit existing forms. When drafts are saved, they replace the form they were drafted from. Edits are _always_ performed on drafts, such that edits do not go live until the user decides.

Drafts are cleaned up from the database periodically.

## Usage

### chi

frm mounts to a `chi.Router` instance and uses the `WorkspaceIDUrlParam` name to look up the workspace that requests belong to.

Example
```go
const chiUrlParamName = "frm_workspace_id"
f, err := frm.New(frm.Args{
	PostgresURL:         os.Getenv("POSTGRES_URL"),
	WorkspaceIDUrlParam: chiUrlParamName, // name of the chi URL parameter name
})
if err != nil {
	panic(err)
}
err = f.Init(context.Background())
if err != nil {
	panic(err)
}
frmchi.Mount(chiRouter, fmt.Sprintf("/frm/{%s}", chiUrlParamName), f)
```

This mounts frm to a router at `/frm/{frm_workspace_id}`.

## Inspiration

This project is heavily inspired by [opnform](https://github.com/jhumanj/opnform). OpnForm is very good, and if you don't have an _embeddable in Go_ requirement, then you should consider OpnForm instead. 
