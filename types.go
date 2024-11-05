package frm

import (
	"encoding/json"

	"github.com/acaloiaro/frm/internal"
)

// Form is a form
type Form internal.Form

// JSON returns the form's JSON-seralized string representation
func (f Form) JSON() string {
	b, err := json.Marshal(f)
	if err != nil {
		return ""
	}

	return string(b)
}
