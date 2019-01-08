package crud

import (
	"log"

	"fmt"
	"html"

	"github.com/pavlo67/punctum/basis/viewshtml"
)

func Table(description Description, data []interface{}, crudOp Operator) string {
	fields := description.ViewList
	if len(fields) < 1 {
		for _, field := range description.Fields {
			viewshtmlField := viewshtml.Field{
				Key:   field.Key,
				Label: field.Key,
			}
			if !field.Updatable {
				viewshtmlField.Type = "view"
			}
			fields = append(fields, viewshtmlField)
		}
	}

	var htmlTable string

	htmlTable += "<tr>"
	for _, viewshtmlField := range fields {
		label := viewshtmlField.Label
		if label == "" {
			label = viewshtmlField.Key
		}
		htmlTable += "<td><b>" + label + "</b></td>\n"
	}
	htmlTable += "</tr>\n"

	var ids string
	for _, row := range data {
		rowMap, err := crudOp.NativeToStringMap(row)
		if err != nil {
			log.Print(err)
		}

		htmlTable += "<tr>"
		for _, viewshtmlField := range fields {
			var update string
			if viewshtmlField.Type != "view" {
				if len(description.FieldsKey) < 1 {
					ids = " data-id_id=\"" + html.EscapeString(rowMap["id"]) + "\""
				} else {
					ids = ""
					for _, key := range description.FieldsKey {
						ids = " data-id_" + html.EscapeString(key) + " =\"" + html.EscapeString(rowMap[key]) + "\""
					}
				}

				update = fmt.Sprintf("&nbsp;<a href=# id=editInTable data-key=\"%s\" %s>//</a>", html.EscapeString(viewshtmlField.Key), ids)
			}
			htmlTable += "<td>" + rowMap[viewshtmlField.Key] + update + "</td>\n"
		}
		htmlTable += "</tr>\n"
	}

	return "<table class=\"cell_border\" width=100% cellpadding=5 cellspacing=0>" + htmlTable + "</table>\n"
}
