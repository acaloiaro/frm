package types

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"slices"
	"sort"

	"github.com/google/uuid"
)

// ErrRequiredNoValueProvided is a form validation error for required fields missing values
var ErrRequiredNoValueProvided = errors.New("This field is required")
var ErrUnknownOptionProvided = errors.New("This field is required, please choose a valid option")

// ValidationErrors is a mapping of form field IDs to the errors validating values submitted to those fields
type ValidationErrors map[string]error

func (v ValidationErrors) Any() bool {
	return len(v) > 0
}

// FormFieldType enum enumerates all possible form field types
//
//go:generate enumer -type FormFieldType -trimprefix FormFieldType -transform=snake -json
type FormFieldType int

const (
	FormFieldTypeTextSingle         FormFieldType = iota // single line of text
	FormFieldTypeTextMultiple                            // multiple lines of text
	FormFieldTypeSingleSelect                            // single-select dropdown
	FormFieldTypeMultiSelect                             // multi-select dropdown
	FormFieldTypeSingleChoice                            // nicely styled radio buttons
	FormFieldTypeSingleChoiceSpaced                      // nicely styled radio buttons, spaced out
)

// FormFieldDataType enum enumerates all possible data types for form fields
//
// This type informs how form field submissions may be used by 'frm' users.
//
//go:generate enumer -type FormFieldDataType -trimprefix FormFieldDataType -transform=snake -json
type FormFieldDataType int

const (
	FormFieldDataTypeText    FormFieldDataType = iota // textual data
	FormFieldDataTypeNumeric                          // numeric data
	FormFieldDataTypeRating                           // chosen values represent a 'rating'
)

func FormFieldDataTypes() []FormFieldDataType {
	return []FormFieldDataType{FormFieldDataTypeText, FormFieldDataTypeNumeric, FormFieldDataTypeRating}
}

// FieldLogicComparator enum enumerates all possible form field logic comparators
//
//go:generate enumer -type FieldLogicComparator -trimprefix FieldLogicComparator -transform=snake -json -text
type FieldLogicComparator int

const (
	FieldLogicComparatorEqual    FieldLogicComparator = iota // target field value is equal to the subject value
	FieldLogicComparatorContains                             // target field value contains the subject value
	FieldLogicComparatorNot                                  // target field value is "not" the subject value
)

// FormFieldOptionOrder enum enumerates all possible ways to order FieldOptions
//
//go:generate enumer -type FormFieldOptionOrder -trimprefix FormFieldOptionOrder -transform=snake -json -text
type FormFieldOptionOrder int

const (
	OptionOrderNatural FormFieldOptionOrder = iota // FieldOptions are ordered naturally according to their order field
	OptionOrderRandom                              // FieldOptions are ordered randomly
)

// FieldLogicTriggerAction enum enumerates all possible field logic trigger actions
//
//go:generate enumer -type FieldLogicTriggerAction -trimprefix FieldLogicTriggerAction -transform=snake -json -text
type FieldLogicTriggerAction int

const (
	FieldLogicTriggerShow    FieldLogicTriggerAction = iota // make the field visible to the user
	FieldLogicTriggerRequire                                // require the user to enter a value
)

// FormFields is a collection of form fields associated with a Form
//
// The underlying type is a map, where keys are form field IDs and values are the corresponding form field
type FormFields map[string]FormField

// FormFieldValues is a collection of form fields submitted to a form
//
// The underlying type is a map, where keys are form field IDs and values are what was submited to the form representing that field
type FormFieldValues map[string]FormFieldSubmission

// FieldOptions are options for single or multi-selector fields
type FieldOptions []Option

// FormField is a field associated with a form
type FormField struct {
	ID           uuid.UUID            `json:"id"`            // field's unique id
	Order        int                  `json:"order"`         // order in which the field appears on forms
	Label        string               `json:"label"`         // field's label (name)
	Logic        *FieldLogic          `json:"logic"`         // UI logic for this field
	Options      FieldOptions         `json:"options"`       // single/multi-select options
	OptionLabels []string             `json:"option_labels"` // option labels are shown below [types.FormFieldTypeSingleChoice] options
	OptionOrder  FormFieldOptionOrder `json:"option_order"`  // the order in which options appear to viewers
	Placeholder  string               `json:"placeholder"`   // placeholder value
	Required     bool                 `json:"required"`      // whether the field is required
	Hidden       bool                 `json:"hidden"`        // whether the field is hidden
	Type         FormFieldType        `json:"type"`          // field type
	DataType     FormFieldDataType    `json:"data_type"`     // the data type form submissions to this field
}

// FormFieldSubmission is a form submission for a particular form field. Form submissions consists of one or more form field submission
type FormFieldSubmission struct {
	ID          uuid.UUID         `json:"id"` // field submission's unique id
	FormFieldID uuid.UUID         `json:"form_field_id"`
	Order       int               `json:"order"`     // order in which the field appeared on the submitted form
	Required    bool              `json:"required"`  // whether the field was requird
	Hidden      bool              `json:"hidden"`    // whether the field was hidden
	Type        FormFieldType     `json:"type"`      // field type
	DataType    FormFieldDataType `json:"data_type"` // the data type of the Value
	Value       any               `json:"value"`     // the value that was submitted
}

// FieldLogic defines logic associated with a field
type FieldLogic struct {
	TargetFieldID     uuid.UUID                `json:"target_field_id"`  // ID of the field to monitor for logic evaluation
	TriggerComparator FieldLogicComparator     `json:"field_comparator"` // comparator to use evaluating target field's value with trigger values
	TriggerValues     []string                 `json:"trigger_values"`   // values that target field's value is compared with
	TriggerActions    FieldLogicTriggerActions `json:"actions"`          // actions to take when the field comparator evaluates true
}

// FieldLogicTriggerActions is a collection of field logic trigger actions
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

// Option is a select option (single and multi)
type Option struct {
	ID       uuid.UUID `json:"id"`
	Value    string    `json:"value"`
	Label    string    `json:"label"`
	Order    int       `json:"order"`
	Selected bool      `json:"-"`
	Disabled bool      `json:"-"`
}

// FormFieldSortByOrder implements sort.Interface for []FormField based on
// the Order field.
type FormFieldSortByOrder []FormField

func (f FormFieldSortByOrder) Len() int           { return len(f) }
func (f FormFieldSortByOrder) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FormFieldSortByOrder) Less(i, j int) bool { return f[i].Order < f[j].Order }

// FormFieldOptionSortNatural implements sort.Interface for [[]Option], sorting options naturally by Order
type FormFieldOptionSortNatural []Option

func (f FormFieldOptionSortNatural) Len() int           { return len(f) }
func (f FormFieldOptionSortNatural) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FormFieldOptionSortNatural) Less(i, j int) bool { return f[i].Order < f[j].Order }

// FormFieldOptionSortRand implements sort.Interface for [[]Option], sorting options randomly
type FormFieldOptionSortRand []Option

func (f FormFieldOptionSortRand) Len() int           { return len(f) }
func (f FormFieldOptionSortRand) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FormFieldOptionSortRand) Less(i, j int) bool { return rand.Int64()%2 == 0 }

// Validate validates values submitted to a form field
func (f FormField) Validate(value []string) (err error) {
	if f.Required {
		if len(value) == 0 {
			return ErrRequiredNoValueProvided
		}
		for _, ffv := range value {
			if ffv == "" {
				return ErrRequiredNoValueProvided
			}
		}
	}

	switch f.Type {
	// ensure that the provided value is one of this field's available options
	case FormFieldTypeSingleSelect, FormFieldTypeMultiSelect, FormFieldTypeSingleChoice, FormFieldTypeSingleChoiceSpaced:
		// Choices-js causes fields to be submitted with am empty value, rather than excluding it. This is a bit hacky.
		if !f.Required && len(value) == 1 && value[0] == "" {
			return nil
		}

		if !allValid(f, value) {
			return ErrUnknownOptionProvided
		}
		return nil
	default:
		return nil
	}
}

// MarshalJSON implements the json.Marshaler interface for FormFieldType
func (f FormField) MarshalJSON() ([]byte, error) {
	id := uuid.Nil
	if f.ID != id {
		id = f.ID
	}

	// only confiugre logic when logic is _completely_ configured
	var logic *FieldLogic
	if f.Logic != nil && f.Logic.TargetFieldID != uuid.Nil && len(f.Logic.TriggerValues) > 0 && f.Logic.TriggerValues[0] != "" {
		logic = f.Logic
	}

	d := struct {
		ID           uuid.UUID            `json:"id"`            // field's unique id
		Order        int                  `json:"order"`         // order in which the field appears on forms
		Label        string               `json:"label"`         // field's label (name)
		Logic        *FieldLogic          `json:"logic"`         // field's logic configuration
		Options      FieldOptions         `json:"options"`       // single/multi-select options
		OptionLabels []string             `json:"option_labels"` // labels for [FormFieldTypeSingleChoice] options
		OptionOrder  FormFieldOptionOrder `json:"option_order"`  // the order in which options appear
		Placeholder  string               `json:"placeholder"`   // placeholder value
		Required     bool                 `json:"required"`      // whether the field is required
		Hidden       bool                 `json:"hidden"`        // whether the field is hidden
		Type         FormFieldType        `json:"type"`          // field type
		DataType     FormFieldDataType    `json:"data_type"`     // field's data type
	}{

		ID:           id,
		Order:        f.Order,
		Label:        f.Label,
		Options:      f.Options,
		OptionLabels: f.OptionLabels,
		OptionOrder:  f.OptionOrder,
		Placeholder:  f.Placeholder,
		Required:     f.Required,
		Hidden:       f.Hidden,
		Type:         f.Type,
		Logic:        logic,
		DataType:     f.DataType,
	}

	return json.Marshal(d)
}

// SortedOptions returns a field's options sorted according to its [OptionOrder]
func (f *FormField) SortedOptions() (sorted []Option) {
	sorted = slices.Clone(f.Options)
	switch f.OptionOrder {
	case OptionOrderNatural:
		sort.Sort(FormFieldOptionSortNatural(sorted))
		return
	case OptionOrderRandom:
		sort.Sort(FormFieldOptionSortRand(sorted))
		return
	default:
		sort.Sort(FormFieldOptionSortNatural(sorted))
		return
	}
}

// allValid checks if all field submission values are valid options
func allValid(field FormField, subset []string) bool {
	set := make(map[string]bool)
	for _, v := range field.Options {
		set[v.Value] = true
	}

	// Check if all form responses are valid options;
	for _, v := range subset {
		if !set[v] {
			return false
		}
	}

	return true
}
