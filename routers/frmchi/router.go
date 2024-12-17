package frmchi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/handlers"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/static"
	"github.com/acaloiaro/frm/types"
	"github.com/acaloiaro/frm/ui"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// This regex parses form field names of the following form
//
// [FIELD_UUID]FIELD_NAME
// Example: [2ad1591d-c852-47b5-a16d-0b90892421c8]label
//
// [FIELD_UUID][SUBGROUP_NAME]FIELD_NAME
// Example: [2ad1591d-c852-47b5-a16d-0b90892421c8][logic]target_field_id
var formFieldIDExtractor = regexp.MustCompile(`^\[([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})\](\[(.+)\]){0,1}(.+){1}?$`)

// Mount mounts frm to the router at the given path
func Mount(f *frm.Frm, router chi.Router, mountPoint string) {
	r := chi.NewRouter()
	r.Use(addMountPointContext)
	router.Mount(mountPoint, r)

	// any requests for which there are no defined chi routes are sent to the "file system"
	// server, serving static content from the embedded filesystem
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		mp := r.Context().Value(handlers.MountPointContextKey).(string)
		// Remove the mount point from the path so the static filesystem paths are resolved relative to its root
		r.URL.Path = strings.ReplaceAll(r.URL.Path, mp, "")
		http.FileServer(http.FS(static.Assets)).ServeHTTP(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		forms, err := internal.Q(ctx, f.PostgresURL).ListForms(r.Context(), f.WorkspaceID)
		if err != nil {
			slog.Error("unable to get forms", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(forms) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		form := forms[0]
		err = ui.Builder((frm.Form)(form)).Render(r.Context(), w)
		if err != nil {
			slog.Error("unable to render builder", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	r.Put("/{form_id}/fields/order", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updatedFields := types.FormFields{}
		for key, field := range r.Form {
			switch key {
			case "order":
				for order, fieldID := range field {
					oldField := form.Fields[fieldID]
					oldField.Order = order
					updatedFields[fieldID] = oldField
				}
			default:
			}
		}

		form.Fields = updatedFields
		form, err = internal.Q(ctx, f.PostgresURL).SaveForm(ctx, internal.SaveFormParams{
			ID:     form.ID,
			Name:   form.Name,
			Fields: updatedFields,
		})
		if err != nil {
			slog.Error("unable to save form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form fields form UI
		err = ui.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = ui.FormView((frm.Form)(form), true).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Get("/{form_id}/logic_configurator/{field_id}/step3", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fieldID, err := uuid.Parse(chi.URLParam(r, "field_id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// the URL param 'id' represents the _target_ field id
		targetFieldID, err := uuid.Parse(r.Form.Get("id"))
		if err != nil {
			slog.Error("unable to parse chosen field id", slog.Any("error", err))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		targetField, ok := form.Fields[targetFieldID.String()]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = ui.LogicConfiguratorStepThree((frm.Form)(form), form.Fields[fieldID.String()], targetField).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(200)
	})

	r.Put("/{form_id}/settings", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		formName := r.Form.Get("name")
		if len(formName) == 0 {
			formName = "New form"
		}
		form, err = internal.Q(ctx, f.PostgresURL).SaveForm(ctx, internal.SaveFormParams{
			ID:     formID,
			Name:   formName,
			Fields: form.Fields,
		})
		if err != nil {
			slog.Error("unable to save new form field")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Re-render the form fields form UI
		err = ui.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form preview
		err = ui.FormView((frm.Form)(form), true).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the configurator form
		err = ui.FormFieldConfigurator((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the configurator form
		err = ui.FormSettings((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render nav, so the title of the form updates
		err = ui.FormBuilderNavTitle((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(200)
	})

	r.Post("/{form_id}/fields", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ft := r.Form.Get("field_type")
		fieldType, err := types.FormFieldTypeString(ft)
		if err != nil {
			slog.Error("unknown form field type", "type", ft)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fields := form.Fields
		order := len(fields) + 1 // place the new field at the end of the field list
		fieldID := uuid.New()
		newField := &types.FormField{ID: fieldID, Type: fieldType, Order: order}
		switch fieldType {
		case types.FormFieldTypeTextSingle:
			newField.Label = "New text field"
			newField.Placeholder = "Write some text"
		case types.FormFieldTypeTextMultiple:
			newField.Label = "New multi-line text field"
			newField.Placeholder = "Write some text"
		case types.FormFieldTypeSingleSelect:
			newField.Label = "New select field"
			newField.Placeholder = "Choose an item"
		case types.FormFieldTypeMultiSelect:
			newField.Label = "New multi select field"
			newField.Placeholder = "Choose items item"
		}

		fields[fieldID.String()] = *newField
		_, err = internal.Q(ctx, f.PostgresURL).SaveForm(ctx, internal.SaveFormParams{
			ID:     formID,
			Name:   form.Name,
			Fields: fields,
		})
		if err != nil {
			slog.Error("unable to save new form field")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Re-render the form fields form UI
		err = ui.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form preview
		err = ui.FormView((frm.Form)(form), true).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the configurator form
		err = ui.FormFieldConfigurator((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(200)
	})

	r.Put("/{form_id}/fields", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updatedFields := form.Fields

		// clear out the logic fields, ensuring logic only gets configured when the UI reflects fully configured logic
		for fid, f := range updatedFields {
			f.Logic = nil
			updatedFields[fid] = f
		}

		for formFieldName, formFieldValues := range r.Form {
			matches := formFieldIDExtractor.FindStringSubmatch(formFieldName)
			if len(matches) < 4 {
				continue
			}
			fieldID := matches[1]
			fieldGroup := matches[3]
			fieldName := matches[4]
			fieldValues := formFieldValues
			isset := len(fieldValues) > 0
			if !isset {
				continue
			}
			switch {
			case fieldName == "required":
				val := (len(fieldValues) > 1 && fieldValues[1] == "on") || (len(fieldValues) > 0 && fieldValues[0] == "on")
				oldField := form.Fields[fieldID]
				oldField.Required = val
				updatedFields[fieldID] = oldField
			case fieldName == "hidden":
				val := (len(fieldValues) > 1 && fieldValues[1] == "on") || (len(fieldValues) > 0 && fieldValues[0] == "on")
				oldField := form.Fields[fieldID]
				oldField.Hidden = val
				updatedFields[fieldID] = oldField
			case fieldName == "label":
				oldField := form.Fields[fieldID]
				oldField.Label = fieldValues[0]
				updatedFields[fieldID] = oldField
			case fieldName == "placeholder":
				oldField := form.Fields[fieldID]
				oldField.Placeholder = fieldValues[0]
				updatedFields[fieldID] = oldField
			case fieldName == "options":
				oldField := form.Fields[fieldID]
				oldField.Options = toFormFieldOption(oldField, fieldValues)
				updatedFields[fieldID] = oldField
			case fieldGroup == ui.FieldGroupLogic && fieldName == ui.FieldLogicTargetFieldID:
				targetFieldID, err := uuid.Parse(fieldValues[0])
				if err != nil {
					continue
				}
				oldField := form.Fields[fieldID]
				if oldField.Logic == nil {
					oldField.Logic = &types.FieldLogic{}
				}
				oldField.Logic.TargetFieldID = targetFieldID
				updatedFields[fieldID] = oldField
			case fieldGroup == ui.FieldGroupLogic && fieldName == ui.FieldLogicTargetFieldValue:
				oldField := form.Fields[fieldID]
				if oldField.Logic == nil {
					oldField.Logic = &types.FieldLogic{}
				}
				oldField.Logic.TriggerValues = fieldValues
				updatedFields[fieldID] = oldField
			case fieldGroup == ui.FieldGroupLogic && fieldName == ui.FieldLogicComparator:
				oldField := form.Fields[fieldID]
				if oldField.Logic == nil {
					oldField.Logic = &types.FieldLogic{}
				}
				oldField.Logic.TriggerComparator, _ = types.FieldLogicComparatorString(fieldValues[0])
				updatedFields[fieldID] = oldField
			case fieldGroup == ui.FieldGroupLogic && fieldName == types.FieldLogicTriggerShow.String():
				oldField := form.Fields[fieldID]
				if oldField.Logic == nil {
					oldField.Logic = &types.FieldLogic{}
				}
				if len(fieldValues) > 0 && fieldValues[0] == "on" {
					oldField.Logic.TriggerActions = []types.FieldLogicTriggerAction{types.FieldLogicTriggerShow}
				}
				updatedFields[fieldID] = oldField
			}
		}

		form.Fields = updatedFields
		form, err = internal.Q(ctx, f.PostgresURL).SaveForm(ctx, internal.SaveFormParams{
			ID:     form.ID,
			Name:   form.Name,
			Fields: updatedFields,
		})
		if err != nil {
			slog.Error("unable to save form", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form fields form UI
		err = ui.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form preview
		err = ui.FormView((frm.Form)(form), true).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Delete("/{form_id}/fields/{field_id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		formID, err := strconv.ParseInt(chi.URLParam(r, "form_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fieldID, err := uuid.Parse(chi.URLParam(r, "field_id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
			WorkspaceID: f.WorkspaceID,
			ID:          formID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		updatedFields := form.Fields
		delete(updatedFields, fmt.Sprint(fieldID))
		form.Fields = updatedFields
		form, err = internal.Q(ctx, f.PostgresURL).SaveForm(ctx, internal.SaveFormParams{
			ID:     form.ID,
			Name:   form.Name,
			Fields: updatedFields,
		})
		if err != nil {
			slog.Error("unable to delete form field", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form fields form UI
		err = ui.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the form preview
		err = ui.FormView((frm.Form)(form), true).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Re-render the configurator form
		err = ui.FormFieldConfigurator((frm.Form)(form)).Render(ctx, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}

// toFormFieldOption takes a list of options sa strings and determines whether the string options represent new options
// being created, in which case an ID/value must be generated for the option, or if the option is amongst the existing
// options for the field being updated.
func toFormFieldOption(field types.FormField, options []string) types.FieldOptions {
	fieldOptions := types.FieldOptions{}
	for _, option := range options {
		var id uuid.UUID
		optionID, err := uuid.Parse(option)
		if err != nil {
			id = uuid.New()
			fieldOptions = append(fieldOptions, types.Option{
				ID:    id,
				Value: id.String(),
				Label: option,
			})
		} else {
			for _, opt := range field.Options {
				if opt.ID == optionID {
					fieldOptions = append(fieldOptions, opt)
				}
			}
		}
	}

	return fieldOptions
}

// addMountPointContext adds the mount point where frm is mounted to the request context
func addMountPointContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mountPoint := chi.RouteContext(r.Context()).RoutePattern()
		if rctx := chi.RouteContext(ctx); rctx != nil {
			for i, urlParam := range rctx.URLParams.Keys {
				if strings.Contains(mountPoint, urlParam) {
					// routePatterns look like: /foo/{bar}/baz
					// Where {bar} is the chi pattern placeholder. These placeholders must be replaced with the actual value
					// that it holders the place for, so we can use the full, realized routePattern as frm's mountpoint
					mountPoint = strings.ReplaceAll(mountPoint, fmt.Sprintf("{%s}", urlParam), rctx.URLParams.Values[i])
				}
			}
		}

		// remove extraneous chi wildcard patterns from the final path
		mountPoint = strings.ReplaceAll(mountPoint, "*", "")

		ctx = context.WithValue(ctx, handlers.MountPointContextKey, mountPoint)
		h.ServeHTTP(w, r.Clone(ctx))
	}

	return http.HandlerFunc(fn)
}
