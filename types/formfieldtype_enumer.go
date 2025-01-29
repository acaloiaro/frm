// Code generated by "enumer -type FormFieldType -trimprefix FormFieldType -transform=snake -json"; DO NOT EDIT.

package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _FormFieldTypeName = "text_singletext_multiplesingle_selectmulti_selectsingle_choice"

var _FormFieldTypeIndex = [...]uint8{0, 11, 24, 37, 49, 62}

const _FormFieldTypeLowerName = "text_singletext_multiplesingle_selectmulti_selectsingle_choice"

func (i FormFieldType) String() string {
	if i < 0 || i >= FormFieldType(len(_FormFieldTypeIndex)-1) {
		return fmt.Sprintf("FormFieldType(%d)", i)
	}
	return _FormFieldTypeName[_FormFieldTypeIndex[i]:_FormFieldTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _FormFieldTypeNoOp() {
	var x [1]struct{}
	_ = x[FormFieldTypeTextSingle-(0)]
	_ = x[FormFieldTypeTextMultiple-(1)]
	_ = x[FormFieldTypeSingleSelect-(2)]
	_ = x[FormFieldTypeMultiSelect-(3)]
	_ = x[FormFieldTypeSingleChoice-(4)]
}

var _FormFieldTypeValues = []FormFieldType{FormFieldTypeTextSingle, FormFieldTypeTextMultiple, FormFieldTypeSingleSelect, FormFieldTypeMultiSelect, FormFieldTypeSingleChoice}

var _FormFieldTypeNameToValueMap = map[string]FormFieldType{
	_FormFieldTypeName[0:11]:       FormFieldTypeTextSingle,
	_FormFieldTypeLowerName[0:11]:  FormFieldTypeTextSingle,
	_FormFieldTypeName[11:24]:      FormFieldTypeTextMultiple,
	_FormFieldTypeLowerName[11:24]: FormFieldTypeTextMultiple,
	_FormFieldTypeName[24:37]:      FormFieldTypeSingleSelect,
	_FormFieldTypeLowerName[24:37]: FormFieldTypeSingleSelect,
	_FormFieldTypeName[37:49]:      FormFieldTypeMultiSelect,
	_FormFieldTypeLowerName[37:49]: FormFieldTypeMultiSelect,
	_FormFieldTypeName[49:62]:      FormFieldTypeSingleChoice,
	_FormFieldTypeLowerName[49:62]: FormFieldTypeSingleChoice,
}

var _FormFieldTypeNames = []string{
	_FormFieldTypeName[0:11],
	_FormFieldTypeName[11:24],
	_FormFieldTypeName[24:37],
	_FormFieldTypeName[37:49],
	_FormFieldTypeName[49:62],
}

// FormFieldTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FormFieldTypeString(s string) (FormFieldType, error) {
	if val, ok := _FormFieldTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _FormFieldTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to FormFieldType values", s)
}

// FormFieldTypeValues returns all values of the enum
func FormFieldTypeValues() []FormFieldType {
	return _FormFieldTypeValues
}

// FormFieldTypeStrings returns a slice of all String values of the enum
func FormFieldTypeStrings() []string {
	strs := make([]string, len(_FormFieldTypeNames))
	copy(strs, _FormFieldTypeNames)
	return strs
}

// IsAFormFieldType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i FormFieldType) IsAFormFieldType() bool {
	for _, v := range _FormFieldTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for FormFieldType
func (i FormFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for FormFieldType
func (i *FormFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FormFieldType should be a string, got %s", data)
	}

	var err error
	*i, err = FormFieldTypeString(s)
	return err
}
