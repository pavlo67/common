package viewshtml

import (
	"html"

	"github.com/pavlo67/punctum/basis"
)

type Attributes map[string]string

type Field struct {
	Key            string        `bson:"key"                       json:"key"`
	Label          string        `bson:"label,omitempty"           json:"label,omitempty"`
	Type           string        `bson:"type,omitempty"            json:"type,omitempty"`
	Format         string        `bson:"format,omitempty"          json:"format,omitempty"`
	AttributesHTML string        `bson:"attributes_html,omitempty" json:"attributes_html,omitempty"`
	Params         basis.Options `bson:"params,omitempty"          json:"params,omitempty"`
}

const NotEmptyKey = "not_empty"
const NoEscapeKey = "no_escape"
const ModelKey = "model"
const AddBlankKey = "add_blank"
const MultiplyKey = "multiply"
const CreateNewKey = "create_new"
const CreateNewTitleKey = "create_new_title"
const GenusKey = "genus"

var NotEmpty = basis.Options{NotEmptyKey: true}
var NoEscape = basis.Options{NoEscapeKey: true}

func AttributesHTML(attributes Attributes) string {
	var attributesHTML string
	for k, v := range attributes {
		attributesHTML += " " + html.EscapeString(k) + `="` + html.EscapeString(v) + `"`
	}

	return attributesHTML
}

func AttributeHTML(key, value string) string {
	return key + `="` + html.EscapeString(value) + `"`
}

func FieldEdit(formID string, field Field, data map[string]string, options map[string]SelectString, frontOps map[string]Operator) (string, string) {

	if field.Type == "view" || field.Type == "text" {
		return FieldView(field, data, options, frontOps)

	} else if field.Type == "button" {
		attributes := AttributesHTML(
			Attributes{
				// using generalNoFormID to add listeners on html pages
				"id":           field.Key,
				"data-form_id": formID,
				"data-value":   data[field.Key],
				"value":        field.Label,
			},
		)
		return "", `<input type="button" ` + attributes + " " + field.AttributesHTML + `/>`
	}

	attributes := AttributeHTML("id", formID+field.Key) + " " + field.AttributesHTML

	var titleHTML = html.EscapeString(field.Label)
	var resHTML string

	if field.Type == "password" {
		resHTML = `<input style="width:100%" type="password" ` + attributes + ` />`
	} else if field.Type == "select" {
		resHTML = HTMLSelectEdit(attributes, options[field.Key], data[field.Key])
	} else if field.Type == "text" {
		text, _ := field.Params["text"].(string)
		resHTML = html.EscapeString(text)
	} else if field.Type == "checkbox" {
		var checked string
		if data[field.Key] != "" {
			checked = " checked"
		}
		resHTML = `<input type="checkbox" ` + attributes + checked + `/>`
	} else if frontOp, ok := frontOps[field.Type]; ok {
		params := map[string]string{
			"form_id": formID,
			"style":   "width:100%",
		}
		resHTML = frontOp.HTMLToEdit(field, data[field.Key], options[field.Key], params)
	} else {
		var value = html.EscapeString(data[field.Key])
		if field.Type == "hidden" {
			resHTML = `<input type="hidden" ` + attributes + ` value="` + value + `" /> `
			titleHTML = ""
		} else if field.Type == "textarea" {
			rows := field.Format
			if rows == "" {
				rows, _ = field.Params["rows"].(string)
			}
			resHTML = `<textarea style="width:100%" ` + attributes + ` rows=` + rows + `>` + value + `</textarea>`
		} else if field.Format == "number" {
			parameters := ` step="` + field.Params.StringDefault("step", "1") + `"`
			if min, ok := field.Params.String("min"); ok {
				parameters += ` min="` + min + `"`
			}
			if max, ok := field.Params.String("max"); ok {
				parameters += ` max="` + max + `"`
			}
			resHTML = `<input type="number"` + parameters + attributes + ` value="` + value + `" />`
		} else if (field.Format == "date") || (field.Format == "time") || (field.Format == "datetime") || (field.Format == "email") || (field.Format == "url") || (field.Format == "color") {
			resHTML = `<input type="` + field.Type + `" ` + attributes + ` value="` + value + `" />`
		} else {
			resHTML = `<input type="` + field.Type + `"style="width:100%" ` + attributes + ` value="` + value + `" />`
		}
	}
	return titleHTML, resHTML
}

// view - not editable data field
// text - text label only (no data field linked to!)

func FieldView(field Field, data map[string]string, options map[string]SelectString, frontOps map[string]Operator) (string, string) {

	var types = []string{"password", "button", "hidden"}
	for _, v := range types {
		if v == field.Type {
			return "", ""
		}
	}

	//if frontOp, ok := frontOps[field.Type]; ok {
	//	params := map[string]string{
	//		// "format": field.Options,
	//		"class": class,
	//		"style": "width:100%",
	//	}
	//	return html.EscapeString(field.Label), frontOp.HTMLToView(field, data[field.Key], nil, params)
	//}

	var resHTML string

	if field.Type == "select" {
		resHTML = HTMLSelectView(options[field.Key], data[field.Key])
	} else if field.Format == "url" {
		var url = html.EscapeString(data[field.Key])
		resHTML = `<a href="` + url + `" target=_blank>` + url + `</a>`

	} else if field.Type == "text" {
		resHTML = field.Format

	} else if field.Type == "checkbox" {
		if data[field.Key] == "on" {
			resHTML = "так"
		} else if field.Params[NotEmptyKey] != true {
			resHTML = "ні"
		}
	} else if field.Params[NotEmptyKey] == true && data[field.Key] == "0" {
		// shows nothing
	} else if field.Params[NoEscapeKey] == true {
		resHTML = data[field.Key]

	} else {
		resHTML = html.EscapeString(data[field.Key])

	}
	return html.EscapeString(field.Label), resHTML
}
