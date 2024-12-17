package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/static"
	"github.com/acaloiaro/frm/types"
	"github.com/acaloiaro/frm/ui"
	"github.com/google/uuid"
)

type contextKey string

var FormIDContextKey contextKey = "frm_form_id"
var FieldIDContextKey contextKey = "frm_field_id"

var ErrNoInstanceAvailable = errors.New("no frm instance is available on the context")
var ErrFormIDNotFound = errors.New("a form ID was not found in the request context")
var ErrFieldIDNotFound = errors.New("a field ID was not found in the request context")

// StaticAssetHandler handles requests for assets embedded in frm's static file system
func StaticAssetHandler(w http.ResponseWriter, r *http.Request) {
	mp := r.Context().Value(internal.MountPointContextKey).(string)
	// Remove the mount point from the path so the static filesystem paths are resolved relative to its root
	r.URL.Path = strings.ReplaceAll(r.URL.Path, mp, "")
	http.FileServer(http.FS(static.Assets)).ServeHTTP(w, r)
}

// FormEditor displays the form editor and previewer
func FormEditor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
}

// UpdateFieldOrder handles requests updating form field order
func UpdateFieldOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
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
}

// LogicConfigurationStep3 handles changes to form field configuration logic, rendering the correct input type for the
// chosen form field
func LogicConfiguratorStep3(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fieldID, err := fieldID(ctx)
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
		ID:          *formID,
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
}

// UpdateSettings handles updates to form settings
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
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
}

// NewField creates new form fields
func NewField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
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
}

// UpdateFields handles updates to a form's fields
func UpdateFields(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
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
}

// DeleteField deletes fields
func DeleteField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fieldID, err := fieldID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	form, err := internal.Q(ctx, f.PostgresURL).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
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
}

// instance gets the frm instance from the request context
func instance(ctx context.Context) (i *frm.Frm, err error) {
	var ok bool
	i, ok = ctx.Value(internal.FrmContextKey).(*frm.Frm)
	if !ok {
		return nil, ErrNoInstanceAvailable
	}
	return
}

// formID gets the form ID from the request context
func formID(ctx context.Context) (formID *int64, err error) {
	var ok bool
	formID, ok = ctx.Value(FormIDContextKey).(*int64)
	if !ok {
		return nil, ErrFormIDNotFound
	}
	return
}

// fieldID gets the field id from the request context
func fieldID(ctx context.Context) (fieldID *uuid.UUID, err error) {
	var ok bool
	fieldID, ok = ctx.Value(FieldIDContextKey).(*uuid.UUID)
	if !ok {
		return nil, ErrFieldIDNotFound
	}
	return
}

// This regex parses form field names of the following form
//
// [FIELD_UUID]FIELD_NAME
// Example: [2ad1591d-c852-47b5-a16d-0b90892421c8]label
//
// [FIELD_UUID][SUBGROUP_NAME]FIELD_NAME
// Example: [2ad1591d-c852-47b5-a16d-0b90892421c8][logic]target_field_id
var formFieldIDExtractor = regexp.MustCompile(`^\[([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})\](\[(.+)\]){0,1}(.+){1}?$`)

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
