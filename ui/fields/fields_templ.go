// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package fields

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/types"
	"github.com/acaloiaro/frm/ui/selector"

	"github.com/google/uuid"
)

// RequiredFieldIndicator visually indicates that a field is required
func RequiredFieldIndicator() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<span class=\"text-red-500 required-dot\">*</span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// FieldLabel is a HTML label for a field
func FieldLabel(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<label for=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 22, Col: 31}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\" class=\"text-slate-700 text-xl pb-5\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(field.Label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 23, Col: 15}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if field.Required {
			templ_7745c5c3_Err = RequiredFieldIndicator().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</label>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// FieldLabel is a HTML label for a field
//
// The "required" indicator CSS watches conditionally displays when any descendant input element has the "required" attribute.
// This is done by being part of a tailwind 'group' that is applied to the container element in [fields.View()]
func LabeledField(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<label for=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 35, Col: 31}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "\" class=\"text-slate-700 text-xl\"><div class=\"pb-5\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(field.Label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 37, Col: 22}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</span> <span class=\"invisible group-has-[:required]:visible text-red-500 required-dot\">*</span></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var5.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</label>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func FieldTypeIcon(fieldType types.FormFieldType) templ.Component {
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
		templ_7745c5c3_Var8 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var8 == nil {
			templ_7745c5c3_Var8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "<div class=\"p-1 rounded-md flex items-center justify-center bg-blue-100 text-blue-900 ml-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch int(fieldType) {
		case int(types.FormFieldTypeTextSingle), int(types.FormFieldTypeTextMultiple):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "<!-- hericons: bards-3-bottom-left --> <svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12\"></path></svg>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeSingleSelect), int(types.FormFieldTypeMultiSelect):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "<!-- heroicons: chevron-up-down --> <svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M8.25 15 12 18.75 15.75 15m-7.5-6L12 5.25 15.75 9\"></path></svg>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeSingleChoice), int(types.FormFieldTypeSingleChoiceSpaced):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "<!-- heroicons: outline star --> <svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z\"></path></svg>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// View displays fields with the appropriate UI component
func View(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var9 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var9 == nil {
			templ_7745c5c3_Var9 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "<div id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("field-container-%s", field.ID.String()))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 71, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if field.Hidden {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, " class=\"group hidden\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, " class=\"group\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "><div class=\"flex flex-col bg-sky-200 rounded-xl p-6 border-red-500 group-has-[:user-invalid]:border-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch field.Type {
		case types.FormFieldTypeTextSingle:
			templ_7745c5c3_Err = singleLineTextView(field).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case types.FormFieldTypeTextMultiple:
			templ_7745c5c3_Err = multiLineTextView(field).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case types.FormFieldTypeSingleSelect, types.FormFieldTypeMultiSelect:
			templ_7745c5c3_Err = selectView(field).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case types.FormFieldTypeSingleChoice:
			templ_7745c5c3_Err = SingleChoice(field).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case types.FormFieldTypeSingleChoiceSpaced:
			templ_7745c5c3_Err = SingleChoiceSpaced(field).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func singleLineTextView(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var11 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var11 == nil {
			templ_7745c5c3_Var11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var12 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "<input id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 100, Col: 25}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "\" name=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var14 string
			templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 101, Col: 27}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "\" placeholder=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var15 string
			templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(field.Placeholder)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 102, Col: 34}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "\" type=\"text\" autocomplete=\"off\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if field.Required {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, " required")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, " _=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var16 string
			templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("on keyup debounced at 250ms trigger field_change(field_id: '%s', value: my.value)", field.ID.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 108, Col: 122}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "\" class=\"flex-1 appearance-none border border-gray-300 dark:border-gray-600 w-full text-gray-700 dark:bg-notion-dark-light dark:text-gray-300 dark:placeholder-gray-500 placeholder-gray-400 shadow-sm focus:outline-none focus:ring-2 focus:border-2 focus:ring-opacity-100 px-4 py-2 text-base resize-y block rounded-xl bg-sky-50\"><div hx-swap-oob=\"true\" id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var17 string
			templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("errors-%s", field.ID.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 111, Col: 74}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "\"></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = LabeledField(field).Render(templ.WithChildren(ctx, templ_7745c5c3_Var12), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func multiLineTextView(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var18 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var18 == nil {
			templ_7745c5c3_Var18 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var19 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "<textarea id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var20 string
			templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 118, Col: 25}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "\" name=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var21 string
			templ_7745c5c3_Var21, templ_7745c5c3_Err = templ.JoinStringErrs(field.ID.String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 119, Col: 27}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var21))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "\" class=\"flex-1 appearance-none border border-gray-300 dark:border-gray-600 w-full text-gray-700 dark:bg-notion-dark-light dark:text-gray-300 dark:placeholder-gray-500 placeholder-gray-400 shadow-sm focus:outline-none focus:ring-2 focus:border-transparent focus:ring-opacity-100 px-4 py-2 text-base resize-y block rounded-xl bg-sky-50\" placeholder=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var22 string
			templ_7745c5c3_Var22, templ_7745c5c3_Err = templ.JoinStringErrs(field.Placeholder)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 121, Col: 34}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var22))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "\" autocomplete=\"off\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if field.Required {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, " required")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, " _=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var23 string
			templ_7745c5c3_Var23, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("on keyup debounced at 250ms trigger field_change(field_id: '%s', value: my.value)", field.ID.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 126, Col: 122}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var23))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "\" rows=\"3\"></textarea><div id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var24 string
			templ_7745c5c3_Var24, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("errors-%s", field.ID.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 130, Col: 51}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var24))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 36, "\" class=\"text-red-400\"></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = LabeledField(field).Render(templ.WithChildren(ctx, templ_7745c5c3_Var19), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func selectView(field types.FormField) templ.Component {
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
		templ_7745c5c3_Var25 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var25 == nil {
			templ_7745c5c3_Var25 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var26 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = selector.Selector(selector.SelectArgs{
				ID:             field.ID.String(),
				Name:           field.ID.String(),
				Label:          "", // make its defulat label styling not render by being blank
				LabelClass:     "text-slate-700 text-xl py-1",
				Required:       field.Required,
				Placeholder:    field.Placeholder,
				Multiple:       field.Type == types.FormFieldTypeMultiSelect,
				Options:        ToSelectorOpts(field.Options, false),
				SearchDisabled: true,
				EditItems:      false,
				Hyperscript:    fmt.Sprintf("on change trigger field_change(field_id: '%s', value: my.value)", field.ID.String()),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 37, " <div id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var27 string
			templ_7745c5c3_Var27, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("errors-%s", field.ID.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/fields/fields.templ`, Line: 152, Col: 51}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var27))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 38, "\" class=\"text-red-400\"></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = LabeledField(field).Render(templ.WithChildren(ctx, templ_7745c5c3_Var26), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// FormFieldTypeLabel is the UI label for FormFieldTypes, e.g. 'text_single' -> "Single line text"
func FormFieldTypeLabel(fieldType types.FormFieldType) templ.Component {
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
		templ_7745c5c3_Var28 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var28 == nil {
			templ_7745c5c3_Var28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 39, "<div class=\"flex gap-3\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = FieldTypeIcon(types.FormFieldType(fieldType)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 40, "<label class=\"w-full cursor-pointer truncate\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch int(fieldType) {
		case int(types.FormFieldTypeTextSingle):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 41, "Single-line text")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeTextMultiple):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 42, "Multi-line text")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeSingleSelect):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 43, "Single select")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeMultiSelect):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 44, "Multi select")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeSingleChoice):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 45, "Single Choice")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case int(types.FormFieldTypeSingleChoiceSpaced):
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 46, "Single Choice (spaced)")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 47, "</label></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// ToSelectorOpts converts a slices of types.Option to a slice of selector.Option
func ToSelectorOpts(opts []types.Option, selectAll bool) (sopts []selector.Option) {
	for _, opt := range opts {
		if selectAll {
			opt.Selected = true
		}
		sopts = append(sopts, (selector.Option)(opt))
	}

	return
}

// ToSelectorOptsStr converts a slices of types.Option to a slice of selector.Option
func ToSelectorOptsStr(opts []string, selectAll bool) (sopts []selector.Option) {
	// TODO: Fix -- adding an empty option because the first select <option> is selected by default, for some reason
	sopts = append(sopts, (selector.Option)(selector.Option{
		Value: "",
	}))
	for _, opt := range opts {
		sopts = append(sopts, selector.Option{
			Label:    opt,
			Value:    opt,
			Selected: selectAll,
		})
	}

	return
}

// FieldsAsSelector returns all of a form's fields as selector.Options to be used in a selector.Selector dropdown
// fieldID is the ID of the field for which the options are being rendered
func FieldsAsSelectorOptions(form frm.Form, fieldID uuid.UUID) (options []selector.Option) {
	for _, field := range form.Fields {
		// fields should not show themselves as options
		if field.ID == fieldID {
			continue
		}
		selected := false
		for _, f := range form.Fields {
			if f.ID == fieldID && f.Logic != nil && f.Logic.TargetFieldID == field.ID {
				selected = true
			}
		}
		options = append(options, selector.Option{
			ID:       field.ID,
			Label:    field.Label,
			Value:    field.ID.String(),
			Selected: selected,
		})
	}
	return
}

// fieldTypeToDefaultValuesJSON converts field types to their devault values.
// This function takes a fieldType and converts it to the default values for that field type,
// and encodes it as a JSON string.
func fieldTypeToDefaultValuesJSON(fieldType types.FormFieldType) string {
	label := ""
	placeholder := ""
	switch int(fieldType) {
	case int(types.FormFieldTypeTextSingle), int(types.FormFieldTypeTextMultiple):
		label = "New Text Field"
		placeholder = "Enter some text"
	case int(types.FormFieldTypeSingleSelect):
		label = "New Single Select Field"
		placeholder = "Choose An Item"
	case int(types.FormFieldTypeMultiSelect):
		label = "New Multi Select Field"
		placeholder = "Choose Items"
	}

	field := types.FormField{
		Label:       label,
		Placeholder: placeholder,
		Required:    false,
		Hidden:      false,
		Type:        fieldType,
	}
	json, err := json.Marshal(field)
	if err != nil {
		return ""
	}
	return string(json)
}

// toFrmFields converts frm.FormFields (map[string]frm.FormField) --> []frm.FormField
func toFrmFields(fields types.FormFields) (ffields []types.FormField) {
	for _, field := range fields {
		ffields = append(ffields, (types.FormField)(field))
	}
	return
}

// Sort form fields by Order
func SortFields(fields types.FormFields) (sorted []types.FormField) {
	sorted = toFrmFields(fields)
	sort.Sort(types.FormFieldSortByOrder(sorted))
	return
}

// FieldName generates the HTML form field name for form fileds
func FieldName(field types.FormField, group, name string) string {
	if group == "" {
		return fmt.Sprintf("[%s]%s", field.ID.String(), name)
	}

	return fmt.Sprintf("[%s][%s]%s", field.ID.String(), group, name)
}

var _ = templruntime.GeneratedTemplate
