package handlers

import (
	"errors"
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
	"github.com/jackc/pgx/v5"
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
	if err != nil || formID == nil {
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
	}
	sc := submission.Get("short_code")
	submission.Del("short_code")
	arg := internal.GetShortCodeParams{
		WorkspaceID: i.WorkspaceID,
		ShortCode:   sc,
	}
	// Submissions without short codes are anonymous, and valid
	shortCode, err := internal.Q(ctx, i.DBArgs).GetShortCode(ctx, arg)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Info("[collector] unable to find provided short code for workspace", "errors", err, "params", arg)
		w.WriteHeader(http.StatusInternalServerError)
	} else if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("[collector] short code not found", "params", arg)
	}

	// TODO: Keep track of the submission id
	// submissionID := r.Form.Get("id")

	// Validation renders whether there are errors or not errors, so that non-erroneous fields can be cleared of error messages
	// as the user corrects validation errors
	allFields := slices.Collect(maps.Keys(f.Fields))
	err = ui.Validation(allFields, errs).Render(ctx, w)
	if err != nil {
		slog.Error("[collector] error while reporting validation error", "error", err)
		return
	}
	// submitted forms only have a submission id when they've been previously submitted and the subject has re-submitted
	if submission.Has("submission_id") {
		// TODO do something with the submission id
		submission.Del("submission_id")
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
	var s internal.FormSubmission
	s, err = internal.Q(ctx, i.DBArgs).SaveSubmission(ctx, internal.SaveSubmissionParams{
		// ID:          submissionID, TODO save submission id
		FormID:      *formID,
		WorkspaceID: i.WorkspaceID,
		SubjectID:   &shortCode.SubjectID,
		Status:      internal.SubmissionStatusPartial,
		Fields:      formFieldValues,
	})
	if err != nil {
		slog.Error("[collector] unable to save submission", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if i.Receiver != nil {
		err = i.Receiver(ctx, (frm.FormSubmission)(s))
		if err != nil {
			slog.Error("[collector] unable to execute submission receiver", "error", err)
		}
	}
	err = ui.ThankYou().Render(ctx, w)
	if err != nil {
		slog.Error("[collector] unable to render thank you page", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
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
