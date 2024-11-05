package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

// FormFieldType enum enumerates all possible form field types
//
//go:generate enumer -type FormFieldType -trimprefix FormFieldType -transform=snake -json
type FormFieldType int

const (
	FormFieldTypeTextSingle   FormFieldType = iota // single line of text
	FormFieldTypeTextMultiple                      // multiple lines of text
	FormFieldTypeSingleSelect                      // single-select dropdown
	FormFieldTypeMultiSelect                       // multi-select dropdown
)

// FieldLogicComparator enum enumerates all possible form field logic comparators
//
//go:generate enumer -type FieldLogicComparator -trimprefix FieldLogicComparator -transform=snake -json -text
type FieldLogicComparator int

const (
	FieldLogicComparatorEqual    FieldLogicComparator = iota // target field value is equal to the subject value
	FieldLogicComparatorContains                             // target field value contains the subject value
)

// FieldLogicTriggerAction enum enumerates all possible field logic trigger actions
//
//go:generate enumer -type FieldLogicTriggerAction -trimprefix FieldLogicTriggerAction -transform=snake -json -text
type FieldLogicTriggerAction int

const (
	FieldLogicTriggerShow FieldLogicTriggerAction = iota // show the field in the to the user
)

// FormFields is a collection of form fields associated with a Form
//
// The underlying type is a map, where keys are form field IDs and values are the corresponding form field
type FormFields map[string]FormField

// FieldOptions are options for single or multi-selector fields
type FieldOptions []Option

// FormField is a field associated with a form
type FormField struct {
	ID          uuid.UUID     `json:"id"`          // field's unique id
	Order       int           `json:"order"`       // order in which the field appears on forms
	Label       string        `json:"label"`       // field's label (name)
	Logic       *FieldLogic   `json:"logic"`       // UI logic for this field
	Options     FieldOptions  `json:"options"`     // single/multi-select options
	Placeholder string        `json:"placeholder"` // placeholder value
	Required    bool          `json:"required"`    // whether the field is required
	Hidden      bool          `json:"hidden"`      // whether the field is hidden
	Type        FormFieldType `json:"type"`        // field type
}

// FieldLogic defines logic associated with a field
type FieldLogic struct {
	TargetFieldID     uuid.UUID                `json:"target_field_id"`  // ID of the field to monitor for logic evaluation
	TriggerComparator FieldLogicComparator     `json:"field_comparator"` // comparator to use evaluating target field's value with trigger values
	TriggerValues     []string                 `json:"trigger_values"`   // values that target field's value is compared with
	TriggerActions    FieldLogicTriggerActions `json:"actions"`          // actions to take when the field comparator evaluates true
}

type FieldLogicTriggerActions []FieldLogicTriggerAction

// Contains determines whether FieldLogicTriggerActions contains some other trigger action
func (f FieldLogicTriggerActions) Contains(a FieldLogicTriggerAction) bool {
	for _, ta := range f {
		if ta == a {
			return true
		}
	}
	return false
}

func (f *FieldLogic) Configured() bool {
	return f.TargetFieldID != uuid.Nil && len(f.TriggerValues) > 0 && len(f.TriggerActions) > 0
}

// Option is a select option (single and multi)
type Option struct {
	ID       uuid.UUID `json:"id"`
	Value    string    `json:"value"`
	Label    string    `json:"label"`
	Selected bool      `json:"-"`
	Disabled bool      `json:"-"`
}

// FormFieldSortByOrder implements sort.Interface for []FormField based on
// the Order field.
type FormFieldSortByOrder []FormField

func (f FormFieldSortByOrder) Len() int           { return len(f) }
func (f FormFieldSortByOrder) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FormFieldSortByOrder) Less(i, j int) bool { return f[i].Order < f[j].Order }

// MarshalJSON implements the json.Marshaler interface for FormFieldType
func (f FormField) MarshalJSON() ([]byte, error) {
	id := uuid.Nil
	if f.ID != id {
		id = f.ID
	}

	d := struct {
		ID          uuid.UUID     `json:"id"`    // field's unique id
		Order       int           `json:"order"` // order in which the field appears on forms
		Label       string        `json:"label"` // field's label (name)
		Logic       *FieldLogic   `json:"logic"`
		Options     FieldOptions  `json:"options"`     // single/multi-select options
		Placeholder string        `json:"placeholder"` // placeholder value
		Required    bool          `json:"required"`    // whether the field is required
		Hidden      bool          `json:"hidden"`      // whether the field is hidden
		Type        FormFieldType `json:"type"`        // field type
	}{

		ID:          id,
		Order:       f.Order,
		Label:       f.Label,
		Options:     f.Options,
		Placeholder: f.Placeholder,
		Required:    f.Required,
		Hidden:      f.Hidden,
		Type:        f.Type,
	}

	if f.Logic != nil && f.Logic.Configured() {
		d.Logic = f.Logic
	} else {
		d.Logic = nil
	}

	return json.Marshal(d)
}
