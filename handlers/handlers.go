package handlers

import (
	"context"
	"net/http"
)

type contextKey string

const (
	MountPointContextKey contextKey = "mount_point_context_key"
)

// AddMountPointContext adds the mount point where frm is mounted to the request context
func AddMountPointContext(mountPoint string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, MountPointContextKey, mountPoint)
			h.ServeHTTP(w, r.Clone(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
