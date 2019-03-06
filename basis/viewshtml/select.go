package viewshtml

import (
	"html"
	"strings"

	"github.com/pavlo67/partes/validator"
	"github.com/pavlo67/punctum/basis"
)

type SelectString [][2]string

const InlineFields = "inline"

func HTMLSelectEdit(attributes string, options SelectString, selected string) string {
	body := ""
	var option string
	for i := 0; i < len(options); i++ {
		body += "<option"
		if options[i][1] != "" {
			option = options[i][1]
			body += ` value="` + html.EscapeString(options[i][1]) + `"`
		} else {
			option = options[i][0]
		}
		if option == selected {
			body += " selected"
		}
		body += ">" + html.EscapeString(options[i][0]) + "</option>\n"
	}
	return `<select ` + attributes + `>` + body + "</select>\n"
}

func HTMLSelectView(values SelectString, selected string) string {
	for i := 0; i < len(values); i++ {
		option := values[i][0]
		if values[i][1] != "" {
			option = values[i][1]
		}
		if option == selected {
			return values[i][0]
		}
	}
	return ""
}

// select string validation ---------------------------------------------------------------------------------

var _ validator.Operator = &SelectStringValidator{}

type SelectStringValidator struct {
	data  string
	label string
	value string
	errs  basis.Errors
}

func NewSelectString(data, label string, values SelectString, trim bool) SelectStringValidator {
	if trim {
		data = strings.TrimSpace(data)
	}
	value := ""
	errs := basis.Errors{validator.BadValue}
	for _, v := range values {
		if v[1] == "" {
			if v[0] == data {
				value = data
				errs = nil
			}
		} else {
			if v[1] == data {
				value = data
				errs = nil
			}
		}
	}

	return SelectStringValidator{label: label, data: data, value: value, errs: errs}
}

func (v SelectStringValidator) Label() string {
	return v.label
}

func (v SelectStringValidator) Errs() basis.Errors {
	return v.errs
}

func (v SelectStringValidator) Value() string {
	return v.value
}
