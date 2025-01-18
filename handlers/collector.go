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
	"github.com/google/uuid"
)

// View renders the form viewer for the collector
func View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	i, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, i)
	if err != nil {
		slog.Error("unable to view form!", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := internal.Q(ctx, i.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: i.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		slog.Error("unable to view form!!", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Render the form collector
	err = ui.Viewer((frm.Form)(f)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ShortCode handles requsts for form short codes and renders the corresponding form
//
// When Forms are submitted via short URL, submissions are attributed to the subject with which the short code was
// generated.
func ShortCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	i, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, i)
	if err != nil {
		slog.Error("unable to view form", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := internal.Q(ctx, i.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: i.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		slog.Error("unable to view form", "error", err)
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
	formID, err := formID(ctx, i)
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
	submission := r.Form
	errs := validate(f, submission)
	if errs.Any() {
		slog.Info("[collector] failed validation", "errors", errs)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// TODO Redirect to a thank-you page
		w.Header().Add("hx-redirect", frm.CollectorPathForm(ctx, *formID))

		w.WriteHeader(http.StatusOK)
	}

	// Validation renders whether there are errors or not errors, so that non-erroneous fields can be cleared of error messages
	// as the user corrects validation errors
	allFields := slices.Collect(maps.Keys(f.Fields))
	err = ui.Validation(allFields, errs).Render(ctx, w)
	if err != nil {
		slog.Error("[collector] error while reporting validation error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	formFieldValues := types.FormFieldValues{}
	for fieldID, fieldValue := range submission {
		formFieldValues[fieldID] = types.FormFieldSubmission{
			ID:          uuid.New(),
			FormFieldID: uuid.MustParse(fieldID), // TODO do not use MustParse
			Order:       f.Fields[fieldID].Order,
			Required:    f.Fields[fieldID].Required,
			Hidden:      f.Fields[fieldID].Hidden,
			Type:        f.Fields[fieldID].Type,
			DataType:    f.Fields[fieldID].DataType,
			Value:       fieldValue,
		}
	}
	_, err = internal.Q(ctx, i.DBArgs).SaveSubmission(ctx, internal.SaveSubmissionParams{
		FormID:      formID,
		WorkspaceID: i.WorkspaceID,
		Status:      internal.SubmissionStatusPartial,
		Fields:      formFieldValues,
	})
	if err != nil {
		slog.Error("[collector] unable to save submission")
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
