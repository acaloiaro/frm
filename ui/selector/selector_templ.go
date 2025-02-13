// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
// package selector provides a single and multi-slect HTML element powered by Choices.js

package selector

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type FieldOptions []Option

// Option is a select option
type Option struct {
	ID       uuid.UUID `json:"id"`
	Value    string    `json:"value"`
	Label    string    `json:"label"`
	Selected bool      `json:"selected"`
	Disabled bool      `json:"disabled"`
}

// ContentID is the HTML element ID where content for this option is rendered when OptionsContent is rendered
func (i Option) ContentID() string {
	return fmt.Sprintf("option-%s-content", i.ID)
}

// SelectArgs are the arguments used to initialize selector.Selector
//
// ID: element ID for the <select> element
//
// Label: The <label> for the <select> element
//
// LabelClass: The CSS class(es) for the <label> element
//
// Multiple: Allow multiple items to be selected
//
// Name: The <select> form element name
//
// Options: The <option>s available to be <select>ed. Use either Options or OptionsContent to provide options to selector.Seletor
//
// OptionsContent: The <option>s available to be <select>ed, and the content to be conditionally rendered when option is selected. Use either Options or OptionsContent, but not both.
//
// Placeholder: Placeholder text shown in the search box
//
// SearchDisabled: Disable the ability to search for items
//
// EditItems: Allow new items to be created through the ui
//
// SelectionChangeEvent: The name of the DOM event to trigger when selections change. This event bubbles. It is triggered on the selector element, which does not handle it. An element up the DOM tree should handle it.
//
// Hyperscript: Attach a hyperscript _ attribute to the <select>.
type SelectArgs struct {
	ID                   string
	Label                string
	LabelClass           string
	Multiple             bool
	Name                 string
	Options              FieldOptions
	OptionsContent       map[Option]templ.Component
	Placeholder          string
	Required             bool
	SearchDisabled       bool
	EditItems            bool
	SelectionChangeEvent string
	Hyperscript          string
}

// MarhsahJSON allows SelectArgs to be serialized to JSON, to be added as a `data-*` field for latest access
func (a SelectArgs) MarshalJSON() (b []byte, err error) {
	d := struct {
		ID                   string
		Name                 string
		Label                string
		LabelClass           string
		Options              []Option
		OptionsContent       map[string]string // mapping of opton ID to the ID of the HTML element where content should render when selected
		Placeholder          string
		Multiple             bool
		SearchDisabled       bool
		EditItems            bool
		SelectionChangeEvent string
		Hyperscript          string
	}{
		ID:                   a.ID,
		Name:                 a.Name,
		Label:                a.Label,
		LabelClass:           a.LabelClass,
		Options:              a.Options,
		OptionsContent:       map[string]string{},
		Placeholder:          a.Placeholder,
		Multiple:             a.Multiple,
		SearchDisabled:       a.SearchDisabled,
		EditItems:            a.EditItems,
		SelectionChangeEvent: a.SelectionChangeEvent,
		Hyperscript:          a.Hyperscript,
	}

	if len(a.OptionsContent) > 0 {
		for option := range a.OptionsContent {
			d.Options = append(d.Options, option)
			d.OptionsContent[option.ID.String()] = option.ContentID()
		}
	}

	b, err = json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return
}

// dataElementID is the element ID of the HTML element that contains the JSON-encoded data for 'args'
func dataElementID(args SelectArgs) string {
	return fmt.Sprintf("selector-data-%s", args.ID)
}

// Selector converts a standard 'select' element into a Choices.js-enriched select element
//
// To use seletor as a multi-selet, the select element must have a the 'multiple' attribute, e.g.
// <select id="my-id" multiple></select>
func Selector(args SelectArgs) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if args.Label != "" {
			var templ_7745c5c3_Var2 = []any{args.LabelClass}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var2...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<label for=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(args.ID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 123, Col: 22}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var2).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 1, Col: 0}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(args.Label)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 124, Col: 15}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if args.Required {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<span class=\"text-red-500 required-dot\">*</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</label> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<select id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(args.ID)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 131, Col: 14}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "\" name=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(args.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 132, Col: 18}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Multiple {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, " multiple")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if args.Required {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, " required=\"true\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, " data-placeholder=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(args.Placeholder)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 139, Col: 37}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Hyperscript != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, " _=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(args.Hyperscript)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 141, Col: 23}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "><option disabled>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(args.Placeholder)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 144, Col: 37}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "</option></select><div id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 string
		templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(dataElementID(args))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 146, Col: 30}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "\" data-args=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var12 string
		templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(templ.JSONString(args))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 146, Col: 67}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "\" _=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var13 string
		templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("init initializeSelect('%s')", dataElementID(args)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 146, Col: 137}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "\"></div><!-- when OptionsContent\n are present, their components render inside of this element when chosen -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(args.OptionsContent) > 0 {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "<!-- content displayed when specific options are selected --> <div id=\"selected-content\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for item, component := range args.OptionsContent {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "<div id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var14 string
				templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(item.ContentID())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/selector/selector.templ`, Line: 153, Col: 30}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "\" class=\"hidden\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = component.Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "<script>\n\t\t/**\n\t\t * Function to wait for predicates.\n\t\t * @param {function() : Promise.<Boolean> | function() : Boolean} predicate\n\t\t * - A function that returns or resolves a bool\n\t\t * @param {Number} [timeout] - Optional maximum waiting time in ms after rejected\n\t\t */\n\t\tfunction waitFor(predicate, timeout) {\n\t\t    return new Promise((resolve, reject) => {\n\t\t        let running = true;\n\n\t\t        const check = async () => {\n\t\t            const res = await predicate();\n\t\t            if (res) return resolve(res);\n\t\t            if (running) setTimeout(check, 100);\n\t\t        };\n\n\t\t        check();\n\n\t\t        if (!timeout) return;\n\t\t        setTimeout(() => {\n\t\t            running = false;\n\t\t            reject();\n\t\t        }, timeout);\n\t\t    });\n\t\t}\n\n\t\tasync function initializeSelect(dataElementID) {\n\t\t\t// wait 10 seconds for choices.js to load \n\t\t\ttry {\n\t\t\t\tawait waitFor(() => typeof Choices !== 'undefined', 10000);\n\t\t\t} catch {\n\t\t\t\tconsole.log(\"unable to load choices!\")\n\t\t\t}\n\n\t\t\tconst argsElement = document.getElementById(dataElementID);\n\t\t\tconst args = JSON.parse(argsElement.getAttribute('data-args'));\n\t\t\tif (args.Options != null) {\n\t\t\t\targs.Options.map((option) => {\n\t\t\t\t\tif (option.value === \"\") {\n\t\t\t\t\t\tconsole.error(`option id: (${option.id}) label: (${option.label}) should have a value! ensure every option provided to SelectorArgs has a 'Value'`)\n\t\t\t\t\t}\n\t\t\t\t});\n\t\t\t}\n\t\t\tvar searchEnabled = true\n\t\t\tif (args.SearchDisabled) {\n\t\t\t\tsearchEnabled = false\n\t\t\t}\n\n\t\t\tvar element = document.getElementById(args.ID)\n\t\t\tvar choices = new Choices(element, {\n\t\t\t\tchoices: args.Options != null ? args.Options : [],\n\t\t\t\tduplicateItemsAllowed: false,\n\t\t\t\tplaceholder: true,\n\t\t\t\tplaceholderValue: args.Placeholder,\n\t\t\t\tsearchEnabled: searchEnabled,\n\t\t\t\teditItems: args.EditItems,\n\t\t\t\tremoveItems: true,\n\t\t\t\tremoveItemButton: true,\n\t\t\t\taddChoices: true,\n\t\t\t\taddItems: true,\n\t\t\t\titemSelectText: \"\",\n\t\t\t\tshouldSort: false,\n\t\t\t});\n\t\t\telement._choices = choices;\n\t\t\n\t\t\t/**\n\t\t\t* When the user's selection(s) have changed, notify a DOM element with an event\n\t\t\t*/\n\t\t\tfunction notifySelected(event) {\n\t\t\t\tif (args.SelectionChangeEvent == \"\") {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t\tvar element = document.getElementById(args.ID)\n\t\t\t\tif (element == null) {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t\telement.dispatchEvent(new CustomEvent(args.SelectionChangeEvent, {bubbles: true, detail: { id: args.ID, value: event.detail.value}}))\n\t\t\t}\n\t\t\t\n\t\t\telement.addEventListener('change', notifySelected, false);\n\t\t}\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
