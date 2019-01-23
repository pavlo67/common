package selectors

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
)

type Selector interface{}

func FromParams(params basis.Options) (Selector, error) {
	if params == nil {
		return nil, nil
	}

	var ok bool
	selectorType := params["type"]
	delete(params, "type")

	switch selectorType {
	case "field":
		var name string
		var values []interface{}
		for k, v := range params {
			if k == "name" {
				name, ok = v.(string)
				if !ok {
					return nil, errors.Errorf("wrong name value type (%#v) for field selector: %#v", v, params)
				}
			} else {
				values = append(values, v)
			}
		}
		return FieldEqual(name, values...), nil

	case "and":
		var values []interface{}
		for _, v := range params {
			values = append(values, v)
		}
		return And(values...), nil

	case "or":
		var values []interface{}
		for _, v := range params {
			values = append(values, v)
		}
		return Or(values...), nil
	}
	// case "match":
	// case "not":
	// case "check_map":

	return nil, errors.Errorf("wrong params value for Selector: %#v", params)
}
