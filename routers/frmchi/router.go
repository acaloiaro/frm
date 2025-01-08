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
	UrlParamFormID urlParam = "frm_form_id"
	// the name for the chi url parameter for field IDs
	UrlParamFieldID urlParam = "frm_field_id"
)

// MountBuilder mounts the frm form builder to the router at the given path
//
// router: The router on which frm mounts the builder.
// f: The frm instance
func MountBuilder(router chi.Router, f *frm.Frm) {
	r := chi.NewRouter()
	r.Use(Middlware(f))
	router.Mount(f.BuilderMountPoint, r)
	r.NotFound(handlers.StaticAssetHandler)
	r.With(addRequestContext).Post("/draft", handlers.NewDraft)
	r.Route(fmt.Sprintf("/forms/{%s}", UrlParamFormID), func(form chi.Router) {
		form = form.With(addRequestContext)
		form.Get("/", handlers.DraftEditor)
		form.Delete("/", handlers.DeleteForm)
		form.Post("/draft", handlers.NewDraft)
		form.Put("/publish", handlers.PublishDraft)
		form.Put("/fields/order", handlers.UpdateFieldOrder)
		form.Put("/settings", handlers.UpdateSettings)
		form.Post("/fields", handlers.NewField)
		form.Put("/fields", handlers.UpdateFields)
		form.Delete(fmt.Sprintf("/fields/{%s}", UrlParamFieldID), handlers.DeleteField)
		form.Get(fmt.Sprintf("/logic_configurator/{%s}/step3", UrlParamFieldID), handlers.LogicConfiguratorStep3)
	})
}

// MountCollector mounts the frm form collector to the router at the given path
//
// router: The router on which frm mounts the collector
// f: The frm instance
func MountCollector(router chi.Router, f *frm.Frm) {
	r := chi.NewRouter()
	r.Use(Middlware(f))
	router.Mount(f.CollectorMountPoint, r)
	r.NotFound(handlers.StaticAssetHandler)
	r.Route(fmt.Sprintf("/{%s}", UrlParamFormID), func(form chi.Router) {
		form = form.With(addRequestContext)
		form.Get("/", handlers.View)
	})
}

// Middlware adds all the context necessary for frm's handlers and path helpers to function
//
// Adds the mount point where frm is mounted to the request context
func Middlware(f *frm.Frm) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			mountPoint := chi.RouteContext(ctx).RoutePattern()
			var workspaceID string
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

			// remove chi wildcard patterns from the final path
			mountPoint = strings.ReplaceAll(mountPoint, "*", "")
			ctx = context.WithValue(ctx, internal.MountPointContextKey, mountPoint)

			// Add the frm instance to the request context, using the workspace ID extracted from the chi route context
			f.WorkspaceID = uuid.MustParse(workspaceID) // TODO don't use MustParse here, figure out what the failure scenario should look like, or switch frm to use string workspace IDs rather than UUIDs, so that parsing is not necessary and provide more flexiblity for users to namespace forms
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
				// populate the request context with the form id from the URL
				if urlParam == string(UrlParamFormID) {
					formID, err := strconv.ParseInt(chi.URLParam(r, string(UrlParamFormID)), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, handlers.FormIDContextKey, &formID)
				}

				// populate the request context with the field id from the URL
				if urlParam == string(UrlParamFieldID) {
					fieldID, err := uuid.Parse(chi.URLParam(r, string(UrlParamFieldID)))
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
