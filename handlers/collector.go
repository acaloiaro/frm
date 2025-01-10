package handlers

import (
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"slices"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/types"
	"github.com/acaloiaro/frm/ui"
)

// View renders the form viewer for the collector
func View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	i, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := internal.Q(ctx, i.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: i.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Render the form collector
	err = ui.Viewer((frm.Form)(f)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Collect handles collector form submissions
func Collect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	i, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := internal.Q(ctx, i.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: i.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = r.ParseForm()
	if err != nil {
		slog.Error("[collector] unable to parse form", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := validate(f, r.Form)
	if errs.Any() {
		slog.Info("[collector] failed validation", "errors", errs)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Add("hx-redirect", frm.CollectorPathForm(ctx, *formID))
		w.WriteHeader(http.StatusOK)
	}

	allFields := slices.Collect(maps.Keys(f.Fields))
	err = ui.Validation(allFields, errs).Render(ctx, w)
	if err != nil {
		slog.Error("[collector] error while reporting validation error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// validate validates forms
func validate(f internal.Form, submission url.Values) (errs types.ValidationErrors) {
	errs = types.ValidationErrors{}
	for submittedFieldID := range maps.Keys(submission) {
		field := f.Fields[submittedFieldID]
		formFieldValue := submission[submittedFieldID]
		if err := field.Validate(formFieldValue); err != nil {
			errs[submittedFieldID] = err
		}
	}
	return errs
}
