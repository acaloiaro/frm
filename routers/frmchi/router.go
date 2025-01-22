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
	// the name of the chi short code parameter for form short codes
	UrlParamShortCode urlParam = "frm_short_code"
)

// Mount mounts the frm form builder and collector to the router
//
// router: The router on which frm mounts the builder.
// f: The frm instance
func Mount(router chi.Router, f *frm.Frm) {
	builder := chi.NewRouter()
	builder.Use(Middlware(f))
	builder.NotFound(handlers.StaticAssetHandler)
	builder.With(addRequestContext).Post("/draft", handlers.NewDraft)
	builder.Route(fmt.Sprintf("/{%s}", UrlParamFormID), func(form chi.Router) {
		form = form.With(addRequestContext)
		form.Get("/", handlers.DraftEditor)
		form.Delete("/", handlers.DeleteForm)
		form.Post("/draft", handlers.NewDraft)
		form.Put("/publish", handlers.PublishDraft)
		form.Put("/fields/order", handlers.UpdateFieldOrder)
		form.Put("/settings", handlers.UpdateSettings)
		form.Post("/fields", handlers.NewField)
		form.Put("/fields", handlers.UpdateFields)
		form.Route(fmt.Sprintf("/fields/{%s}", UrlParamFieldID), func(fields chi.Router) {
			fields.Delete("/", handlers.DeleteField)
			fields.Get("/logic/choices", handlers.LogicConfiguratorChoices)
		})
	})

	collector := chi.NewRouter()
	collector.Use(Middlware(f))
	collector.NotFound(handlers.StaticAssetHandler)
	collector.With(addRequestContext).Route(fmt.Sprintf("/s/{%s}", UrlParamShortCode), func(sc chi.Router) {
		sc.Get("/", handlers.ShortCode)
	})
	collector.Route(fmt.Sprintf("/{%s}", UrlParamFormID), func(form chi.Router) {
		form = form.With(addRequestContext)
		form.Get("/", handlers.View)
		form.Post("/", handlers.Collect)
	})

	router.Mount(f.BuilderMountPoint, builder)
	router.Mount(f.CollectorMountPoint, collector)
}

// Middlware adds all the context necessary for frm's handlers and path helpers to function
//
// Adds the mount point where frm is mounted to the request context
func Middlware(f *frm.Frm) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			builderMountPoint := chi.RouteContext(ctx).RoutePattern()
			collectorMountPoint := chi.RouteContext(ctx).RoutePattern()
			var workspaceID string
			if rctx := chi.RouteContext(ctx); rctx != nil {
				for i, paramName := range rctx.URLParams.Keys {
					paramValue := rctx.URLParams.Values[i]
					paramPlaceholder := fmt.Sprintf("{%s}", paramName)

					// Route patterns look like: /foo/{bar}/baz
					// Where {bar} is the chi pattern placeholder. These placeholders must be replaced with the actual value
					// that it holds a place for, so we can use the full, realized routePattern as frm's mountpoint
					if strings.Contains(f.BuilderMountPoint, paramName) {
						builderMountPoint = strings.ReplaceAll(f.BuilderMountPoint, paramPlaceholder, paramValue)
					}

					if strings.Contains(f.CollectorMountPoint, paramName) {
						collectorMountPoint = strings.ReplaceAll(f.CollectorMountPoint, paramPlaceholder, paramValue)
					}

					// extract the workspace id
					if paramName == f.WorkspaceIDUrlParam {
						workspaceID = paramValue
					}
				}
			}

			// remove chi wildcard patterns from the final path
			builderMountPoint = strings.ReplaceAll(builderMountPoint, "*", "")
			collectorMountPoint = strings.ReplaceAll(collectorMountPoint, "*", "")
			ctx = context.WithValue(ctx, internal.BuilderMountPointContextKey, builderMountPoint)
			ctx = context.WithValue(ctx, internal.CollectorMountPointContextKey, collectorMountPoint)

			// Add the frm instance to the request context, using the workspace ID extracted from the chi route context
			f.WorkspaceID = workspaceID
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
				switch urlParam {
				case string(UrlParamFormID):
					formID, err := strconv.ParseInt(chi.URLParam(r, string(UrlParamFormID)), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, internal.FormIDContextKey, &formID)
				case string(UrlParamFieldID):
					fieldID, err := uuid.Parse(chi.URLParam(r, string(UrlParamFieldID)))
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, internal.FieldIDContextKey, &fieldID)
				case string(UrlParamShortCode):
					shortCode := chi.URLParam(r, string(UrlParamShortCode))
					if shortCode == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					ctx = context.WithValue(ctx, internal.ShortCodeContextKey, &shortCode)
				}
			}
		}
		h.ServeHTTP(w, r.Clone(ctx))
	}

	return http.HandlerFunc(fn)
}
