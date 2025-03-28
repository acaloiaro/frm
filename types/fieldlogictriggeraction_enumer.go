// Code generated by "enumer -type FieldLogicTriggerAction -trimprefix FieldLogicTriggerAction -transform=snake -json -text"; DO NOT EDIT.

package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _FieldLogicTriggerActionName = "field_logic_trigger_showfield_logic_trigger_require"

var _FieldLogicTriggerActionIndex = [...]uint8{0, 24, 51}

const _FieldLogicTriggerActionLowerName = "field_logic_trigger_showfield_logic_trigger_require"

func (i FieldLogicTriggerAction) String() string {
	if i < 0 || i >= FieldLogicTriggerAction(len(_FieldLogicTriggerActionIndex)-1) {
		return fmt.Sprintf("FieldLogicTriggerAction(%d)", i)
	}
	return _FieldLogicTriggerActionName[_FieldLogicTriggerActionIndex[i]:_FieldLogicTriggerActionIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _FieldLogicTriggerActionNoOp() {
	var x [1]struct{}
	_ = x[FieldLogicTriggerShow-(0)]
	_ = x[FieldLogicTriggerRequire-(1)]
}

var _FieldLogicTriggerActionValues = []FieldLogicTriggerAction{FieldLogicTriggerShow, FieldLogicTriggerRequire}

var _FieldLogicTriggerActionNameToValueMap = map[string]FieldLogicTriggerAction{
	_FieldLogicTriggerActionName[0:24]:       FieldLogicTriggerShow,
	_FieldLogicTriggerActionLowerName[0:24]:  FieldLogicTriggerShow,
	_FieldLogicTriggerActionName[24:51]:      FieldLogicTriggerRequire,
	_FieldLogicTriggerActionLowerName[24:51]: FieldLogicTriggerRequire,
}

var _FieldLogicTriggerActionNames = []string{
	_FieldLogicTriggerActionName[0:24],
	_FieldLogicTriggerActionName[24:51],
}

// FieldLogicTriggerActionString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FieldLogicTriggerActionString(s string) (FieldLogicTriggerAction, error) {
	if val, ok := _FieldLogicTriggerActionNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _FieldLogicTriggerActionNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to FieldLogicTriggerAction values", s)
}

// FieldLogicTriggerActionValues returns all values of the enum
func FieldLogicTriggerActionValues() []FieldLogicTriggerAction {
	return _FieldLogicTriggerActionValues
}

// FieldLogicTriggerActionStrings returns a slice of all String values of the enum
func FieldLogicTriggerActionStrings() []string {
	strs := make([]string, len(_FieldLogicTriggerActionNames))
	copy(strs, _FieldLogicTriggerActionNames)
	return strs
}

// IsAFieldLogicTriggerAction returns "true" if the value is listed in the enum definition. "false" otherwise
func (i FieldLogicTriggerAction) IsAFieldLogicTriggerAction() bool {
	for _, v := range _FieldLogicTriggerActionValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for FieldLogicTriggerAction
func (i FieldLogicTriggerAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for FieldLogicTriggerAction
func (i *FieldLogicTriggerAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FieldLogicTriggerAction should be a string, got %s", data)
	}

	var err error
	*i, err = FieldLogicTriggerActionString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for FieldLogicTriggerAction
func (i FieldLogicTriggerAction) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for FieldLogicTriggerAction
func (i *FieldLogicTriggerAction) UnmarshalText(text []byte) error {
	var err error
	*i, err = FieldLogicTriggerActionString(string(text))
	return err
}
