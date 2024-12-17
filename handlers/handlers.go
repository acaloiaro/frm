package handlers

type contextKey string

const (
	// MountPointContextKey is the context key for retrieving frm's mount point from the request context
	MountPointContextKey contextKey = "mount_point_context_key"
)
