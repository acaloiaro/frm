package frmchi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/handlers"
	"github.com/acaloiaro/frm/internal"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type urlParam string

const (
	// the name for the chi url parameter for form IDs
	urlParamFormID urlParam = "frm_form_id"
	// the name for the chi url parameter for field IDs
	urlParamFieldID urlParam = "frm_field_id"
)

// Mount mounts frm to the router at the given path
func Mount(router chi.Router, mountPoint string, f *frm.Frm) {
	r := chi.NewRouter()
	r.Use(addFrmContext(f))
	router.Mount(mountPoint, r)
	rc := r.With(addRequestContext)
	rc.NotFound(handlers.StaticAssetHandler)
	rc.Route(fmt.Sprintf("/forms/{%s}", urlParamFormID), func(form chi.Router) {
		form.Get("/", handlers.DraftEditor)
		form.Delete("/", handlers.DeleteForm)
		form.Post("/draft", handlers.NewDraft)
		form.Put("/publish", handlers.PublishDraft)
		form.Put("/fields/order", handlers.UpdateFieldOrder)
		form.Put("/settings", handlers.UpdateSettings)
		form.Post("/fields", handlers.NewField)
		form.Put("/fields", handlers.UpdateFields)
		form.Delete(fmt.Sprintf("/fields/{%s}", urlParamFieldID), handlers.DeleteField)
		form.Get(fmt.Sprintf("/logic_configurator/{%s}/step3", urlParamFieldID), handlers.LogicConfiguratorStep3)
		form.NotFound(handlers.StaticAssetHandler)
	})
}

// addFrmContext adds all the context necessary for its handlers to function
//
// 1. Adds the mount point where frm is mounted to the request context
// 2. Add an frm instance to the request context, to be used by handlers
func addFrmContext(f *frm.Frm) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var workspaceID string
			mountPoint := chi.RouteContext(r.Context()).RoutePattern()
			if rctx := chi.RouteContext(ctx); rctx != nil {
				for i, urlParam := range rctx.URLParams.Keys {
					if strings.Contains(mountPoint, urlParam) {
						// routePatterns look like: /foo/{bar}/baz
						// Where {bar} is the chi pattern placeholder. These placeholders must be replaced with the actual value
						// that it holders the place for, so we can use the full, realized routePattern as frm's mountpoint
						mountPoint = strings.ReplaceAll(mountPoint, fmt.Sprintf("{%s}", urlParam), rctx.URLParams.Values[i])
					}

					// extract the workspace id
					if urlParam == f.WorkspaceIDUrlParam {
						workspaceID = rctx.URLParams.Values[i]
					}
				}
			}

			// remove extraneous chi wildcard patterns from the final path
			mountPoint = strings.ReplaceAll(mountPoint, "*", "")
			ctx = context.WithValue(ctx, internal.MountPointContextKey, mountPoint)

			// Add the frm instance to the request context, using the workspace ID extracted from the chi route context
			f.WorkspaceID = uuid.MustParse(workspaceID) // TODO don't use MustParse here, figure out what the failure scenario should look like
			ctx = context.WithValue(ctx, internal.FrmContextKey, f)

			h.ServeHTTP(w, r.Clone(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func addRequestContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if rctx := chi.RouteContext(ctx); rctx != nil {
			for _, urlParam := range rctx.URLParams.Keys {
				// populate the request context with the form id form the URL
				if urlParam == string(urlParamFormID) {
					formID, err := strconv.ParseInt(chi.URLParam(r, string(urlParamFormID)), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, handlers.FormIDContextKey, &formID)
				}

				// populate the request context with the field id form the URL
				if urlParam == string(urlParamFieldID) {
					fieldID, err := uuid.Parse(chi.URLParam(r, string(urlParamFieldID)))
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, handlers.FieldIDContextKey, &fieldID)
				}
			}

		}
		h.ServeHTTP(w, r.Clone(ctx))
	}

	return http.HandlerFunc(fn)
}
