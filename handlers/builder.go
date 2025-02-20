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
	"github.com/acaloiaro/frm/ui/builder"
	"github.com/acaloiaro/frm/ui/collector"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrFormIDNotFound = errors.New("a form ID was not found in the request context")
var ErrFieldIDNotFound = errors.New("a field ID was not found in the request context")

// StaticAssetHandler handles requests for assets embedded in frm's static file system
func StaticAssetHandler(w http.ResponseWriter, r *http.Request) {
	// static assets could feasibly be loaded from either the collector or builder mount point. this arbitrarily chooses
	// the collector's mount point for all static assets
	mp := r.Context().Value(internal.CollectorMountPointContextKey).(string)

	// mp ends in a slash (e.g. foo/bar/), and we want to remove /foo/bar/static from the path prefix before searching for
	// paths in the static file system. Join mp with "static" to form foo/bar/static as path prefix to be removed
	// before searching inside the file system for files.
	pathPrefix := fmt.Sprintf("%s/%s", mp, "static")

	// Remove the mount point and static prefix from the path so the static filesystem paths are resolved relative to
	// the file system's root, e.g. if frm is mounted at /foo/bar, /foo/bar/static is removed from the URL path before
	// searching inside the file system.
	r.URL.Path = strings.ReplaceAll(r.URL.Path, pathPrefix, "")

	http.FileServer(http.FS(static.Assets)).ServeHTTP(w, r)
}

// DraftEditor displays the form editor and previewer
func DraftEditor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(r.Context(), internal.GetDraftParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *id,
	})
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("unable to get forms", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = builder.Builder((frm.Form)(draft)).Render(r.Context(), w)
	if err != nil {
		slog.Error("unable to render builder", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// UpdateFieldOrder handles requests updating form field order
func UpdateFieldOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
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
				oldField := draft.Fields[fieldID]
				oldField.Order = order
				updatedFields[fieldID] = oldField
			}
		default:
		}
	}

	draft.Fields = updatedFields
	draft, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		ID:     draft.ID,
		Name:   draft.Name,
		Fields: updatedFields,
	})
	if err != nil {
		slog.Error("unable to save draft", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form fields form UI
	err = builder.FormFieldsForm((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = collector.FormView(collector.ViewerArgs{Form: (frm.Form)(draft), Preview: true}).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// LogicConfiguratorChoices handles changes to form field configuration logic, rendering the correct input type for the
// chosen form field
func LogicConfiguratorChoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
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

	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	targetField, ok := draft.Fields[targetFieldID.String()]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = builder.LogicConfiguratorStepThree((frm.Form)(draft), draft.Fields[fieldID.String()], targetField).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(200)
}

// UpdateSettings handles updates to form settings
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	draftID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	form, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *draftID,
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
	form, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		ID:     draftID,
		Name:   formName,
		Fields: form.Fields,
	})
	if err != nil {
		slog.Error("unable to save new form field")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Re-render the form fields form UI
	err = builder.FormFieldsForm((frm.Form)(form)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form preview
	err = collector.FormView(collector.ViewerArgs{Form: (frm.Form)(form), Preview: true}).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the configurator form
	err = builder.FormFieldConfigurator((frm.Form)(form)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the configurator form
	err = builder.FormSettings((frm.Form)(form)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render nav, so the title of the form updates
	err = builder.FormBuilderNavTitle((frm.Form)(form)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(200)
}

// NewField creates new form fields
func NewField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
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

	fields := draft.Fields
	order := len(fields) + 1 // place the new field at the end of the field list
	fieldID := uuid.New()
	newField := &types.FormField{ID: fieldID, Type: fieldType, Order: order}

	// FIELD_TYPES: field types may be added/modified/removed below
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
		newField.Placeholder = "Choose items"
	case types.FormFieldTypeSingleChoice:
		newField.Label = "New single choice field"
	case types.FormFieldTypeSingleChoiceSpaced:
		newField.Label = "New single choice field (spaced)"
	}

	fields[fieldID.String()] = *newField
	_, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		ID:     formID,
		Name:   draft.Name,
		Fields: fields,
	})
	if err != nil {
		slog.Error("unable to save new form field")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Re-render the form fields form UI
	err = builder.FormFieldsForm((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form preview
	err = collector.FormView(collector.ViewerArgs{Form: (frm.Form)(draft), Preview: true}).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the configurator form
	err = builder.FormFieldConfigurator((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(200)
}

// UpdateFields handles updates to form draft fields
//
// This endpoint updates every draft field. If a field is not present in the request, it will not be present on the
// draft after this endpoint succeeds.
func UpdateFields(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
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

	// newFields becomes the draft's full set of fields after this endpoint succeeds
	newFields := map[string]*types.FormField{}

	// iterate over all submitted fields, adding each one to 'newFields'
	for formFieldName, formFieldValues := range r.Form {
		matches := formFieldIDExtractor.FindStringSubmatch(formFieldName)
		if len(matches) < 4 {
			slog.Warn("skipping field: unable to match field naming convention", "field_name", formFieldName)
			continue
		}
		fieldID := matches[1]
		fieldGroup := matches[3]
		fieldName := matches[4]
		fieldValues := formFieldValues

		id, err := uuid.Parse(fieldID)
		if err != nil {
			slog.Error("skipping field: unable to parse field id", "field_id", fieldID)
			continue
		}

		var field *types.FormField
		var present bool
		if field, present = newFields[fieldID]; !present {
			// The PUT /fields endpoint is only called after a field exists. Thus, it is appropriate to use the 'order'
			// and 'type' of the draft's field when updating, since PUT /fields does not affect order or type
			field = &types.FormField{
				ID:    id,
				Order: draft.Fields[fieldID].Order,
				Type:  draft.Fields[fieldID].Type,
				Logic: &types.FieldLogic{},
			}
			newFields[fieldID] = field
		}

		// parse specific field update requests and update the corresponding field accordingly
		switch {
		case fieldName == "required":
			required := (len(fieldValues) > 1 && fieldValues[1] == "on") || (len(fieldValues) > 0 && fieldValues[0] == "on")
			field.Required = required
		case fieldName == "hidden":
			hidden := (len(fieldValues) > 1 && fieldValues[1] == "on") || (len(fieldValues) > 0 && fieldValues[0] == "on")
			field.Hidden = hidden
		case fieldName == "label":
			field.Label = fieldValues[0]
		case fieldName == "placeholder":
			field.Placeholder = fieldValues[0]
		case fieldName == "options":
			field.Options = toFormFieldOption(draft.Fields[fieldID], fieldValues)
		case fieldName == "option_labels":
			field.OptionLabels = fieldValues
		case fieldName == "option_ordering":
			field.OptionOrder, err = types.FormFieldOptionOrderString(fieldValues[0])
			if err != nil {
				field.OptionOrder = types.OptionOrderNatural
			}
		// field logic, target field chosen
		case fieldGroup == builder.FieldGroupLogic && fieldName == builder.FieldLogicTargetFieldID:
			targetFieldID, err := uuid.Parse(fieldValues[0])
			if err != nil {
				continue
			}
			field.Logic.TargetFieldID = targetFieldID
		// field logic, subject field chosen
		case fieldGroup == builder.FieldGroupLogic && fieldName == builder.FieldLogicTargetFieldValue:
			field.Logic.TriggerValues = fieldValues
		// field logic, comparator chosen
		case fieldGroup == builder.FieldGroupLogic && fieldName == builder.FieldLogicComparator:
			field.Logic.TriggerComparator, _ = types.FieldLogicComparatorString(fieldValues[0])
		// field logic, actions to take
		case fieldGroup == builder.FieldGroupLogic && fieldName == builder.FieldLogicActions:
			if len(fieldValues) > 0 {
				actions := []types.FieldLogicTriggerAction{}
				for _, fv := range fieldValues {
					flta, err := types.FieldLogicTriggerActionString(fv)
					if err != nil {
						continue
					}
					actions = append(actions, flta)
				}
				field.Logic.TriggerActions = actions
			}
		}
	}

	ff := types.FormFields{}
	for fieldID, fptr := range newFields {
		ff[fieldID] = *fptr
	}
	draft, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		ID:     draft.ID,
		Name:   draft.Name,
		Fields: ff,
	})
	if err != nil {
		slog.Error("unable to save form", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form fields form UI
	err = builder.FormFieldsForm((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form preview
	err = collector.FormView(collector.ViewerArgs{Form: (frm.Form)(draft), Preview: true}).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteField deletes fields
func DeleteField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fieldID, err := fieldID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	draft, err := internal.Q(ctx, f.DBArgs).GetDraft(ctx, internal.GetDraftParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	updatedFields := draft.Fields
	delete(updatedFields, fmt.Sprint(fieldID))
	draft.Fields = updatedFields
	draft, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		ID:     draft.ID,
		Name:   draft.Name,
		Fields: updatedFields,
	})
	if err != nil {
		slog.Error("unable to delete form field", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form fields form UI
	err = builder.FormFieldsForm((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the form preview
	err = collector.FormView(collector.ViewerArgs{Form: (frm.Form)(draft), Preview: true}).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Re-render the configurator form
	err = builder.FormFieldConfigurator((frm.Form)(draft)).Render(ctx, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// ChangeStatus changes the status of forms
func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		slog.Error("unable to change form status", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	form, err := internal.Q(ctx, f.DBArgs).GetForm(ctx, internal.GetFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		slog.Error("unable to get form", slog.Any("error", err), slog.Any("workspace_id", f.WorkspaceID), slog.Any("form_id", *formID))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status := frm.FormStatus(r.FormValue("status"))
	form, err = internal.Q(ctx, f.DBArgs).SaveForm(ctx, internal.SaveFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          form.ID,
		Name:        form.Name,
		Fields:      form.Fields,
		Status:      status,
	})
	if err != nil {
		slog.Error("unable to change form status", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// NewDraft creates new drafts
//
// Drafts may be created from scratch by not providing a formID in the request
// Drafts may be created from existing forms by providing a formID in the request
// Drafts may be "clones" of existing forms by providing the "clone" URL parameter, the only differences between clones
// and drafts created from existing forms is that cloned form names are suffixed, and don't recall their parent form.
func NewDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		slog.Error("unable to create draft", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isClone := r.URL.Query().Get("clone") != ""
	suffix := r.URL.Query().Get("name_suffix")

	var draft internal.Form
	// dont check for errors here, because this endpoint handles both new drafts and drafts from existing forms
	formID, _ := formID(ctx, f)

	q := internal.Q(ctx, f.DBArgs)
	draftParams := &internal.SaveFormParams{
		WorkspaceID: f.WorkspaceID,
		Fields:      types.FormFields{},
	}
	if formID != nil {
		copy, err := f.CopyForm(ctx, frm.CopyFormArgs{
			ID:               *formID,
			ForgetParentForm: isClone,
			NameSuffix:       suffix,
		})
		if err != nil {
			slog.Error("unable to create draft", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		draft = (internal.Form)(copy)
	} else {
		draftParams.Name = "New form"
		draft, err = q.SaveForm(ctx, *draftParams)
		if err != nil {
			slog.Error("unable to create draft", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var event string
	if isClone {
		event = frm.EventCloneCreated
	} else {
		event = frm.EventDraftCreated
	}
	w.Header().Add("HX-Trigger", fmt.Sprintf("{\"%s\": {\"draft_id\": \"%d\"}}", event, draft.ID))
	w.WriteHeader(http.StatusCreated)
}

// PublishDraft copies drafts to published forms
func PublishDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		slog.Error("unable to publish draft", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	draftID, err := formID(ctx, f)
	if err != nil {
		slog.Error("unable to publish draft", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tx, err := internal.Tx(ctx, f.DBArgs)
	if err != nil {
		slog.Error("unable to publish draft", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := internal.Q(ctx, f.DBArgs).WithTx(tx)
	_, err = q.PublishDraft(ctx, *draftID)
	if err != nil {
		slog.Error("unable to publish draft", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = q.DeleteForm(ctx, internal.DeleteFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *draftID,
	})
	if err != nil {
		slog.Error("unable to publish draft", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)

	w.WriteHeader(http.StatusNoContent)
}

// DeleteForm deletes forms
func DeleteForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := frm.Instance(ctx)
	if err != nil {
		slog.Error("unable to delete form", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	formID, err := formID(ctx, f)
	if err != nil {
		slog.Error("unable to delete form", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = internal.Q(ctx, f.DBArgs).DeleteForm(ctx, internal.DeleteFormParams{
		WorkspaceID: f.WorkspaceID,
		ID:          *formID,
	})
	if err != nil {
		slog.Error("unable to delete form", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// formID gets the form ID from the request context
func formID(ctx context.Context, f *frm.Frm) (formID *int64, err error) {
	var ok bool
	formID, ok = ctx.Value(internal.FormIDContextKey).(*int64)
	if !ok {
		if shortCode, ok := ctx.Value(internal.ShortCodeContextKey).(*string); ok {
			s, err := internal.Q(ctx, f.DBArgs).GetShortCode(ctx, internal.GetShortCodeParams{
				WorkspaceID: f.WorkspaceID,
				ShortCode:   *shortCode,
			})
			if err != nil {
				return nil, err
			}
			return s.FormID, nil
		}
		return nil, ErrFormIDNotFound
	}
	return
}

// fieldID gets the field id from the request context
func fieldID(ctx context.Context) (fieldID *uuid.UUID, err error) {
	var ok bool
	fieldID, ok = ctx.Value(internal.FieldIDContextKey).(*uuid.UUID)
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

// toFormFieldOption takes a list of options as strings and determines whether the string options represent new options
// being created, in which case an ID/value must be generated for the option, or if the option is amongst the existing
// options for the field being updated.
func toFormFieldOption(field types.FormField, options []string) types.FieldOptions {
	fieldOptions := types.FieldOptions{}
	for i, option := range options {
		var id uuid.UUID
		optionID, err := uuid.Parse(option)
		if err != nil {
			id = uuid.New()
			fieldOptions = append(fieldOptions, types.Option{
				ID:    id,
				Value: id.String(),
				Label: option,
				Order: i,
			})
		} else {
			for _, opt := range field.Options {
				opt.Order = i
				if opt.ID == optionID {
					fieldOptions = append(fieldOptions, opt)
				}
			}
		}
	}

	return fieldOptions
}
