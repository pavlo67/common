package selectors

import (
	"github.com/pavlo67/punctum/auth"
	"github.com/pkg/errors"
)

// !!! this package requires on auth.IDentity & basis.UserIS types due to its application usage

var ErrBadSelector = errors.New("bad selectors")

var NoValues = []interface{}{}

type PrepareIDs func(userIS auth.ID) ([]string, error)

// FieldEqual --------------------------------------------------------------------

type in struct {
	field  string
	values []interface{}
}

func FieldEqual(field string, values ...interface{}) Selector {
	return &in{field, values}
}

func FieldStr(field string, valuesStr ...string) Selector {
	var values []interface{}
	for _, v := range valuesStr {
		values = append(values, v)
	}

	return &in{field, values}
}

// FieldOp --------------------------------------------------------------------

type op struct {
	field, op string
	value     interface{}
}

func FieldOp(field, operation string, value interface{}) Selector {
	return &op{field, operation, value}
}

// Match --------------------------------------------------------------------

type match struct {
	field string
	value string
	mode  string
}

func Match(field, value, mode string) Selector {
	return &match{field, value, mode}
}

// And, Or ------------------------------------------------------------------

type multi struct {
	oper   string
	values []interface{}
}

func Or(values ...interface{}) Selector {
	return &multi{"OR", values}
}

func And(values ...interface{}) Selector {
	return &multi{"AND", values}
}

// Not ----------------------------------------------------------------------

type not struct {
	term interface{}
}

func Not(value interface{}) Selector {
	a := not{value}
	return &a
}

// check map ----------------------------------------------------------------

//type CheckMapHelper func(map[string]string) (bool, error)
//
//func CheckMap(selector Selector, data map[string]string) (bool, error) {
//	if selector == nil {
//		return true, nil
//	}
//
//	// log.Println("%s: %#v", reflect.TypeOf(selector), selector)
//
//	switch s := selector.(type) {
//	case CheckMapHelper:
//		return s(data)
//	case *in:
//		for _, v := range s.values {
//			//log.Println("check:", s.field, v, data[s.field], reflect.TypeOf(v))
//			if value, ok := v.(string); ok {
//				if value == data[s.field] {
//					return true, nil
//				}
//			} else if value, ok := v.(int); ok {
//				if strconv.Itoa(value) == data[s.field] {
//					return true, nil
//				}
//			} else if value, ok := v.(auth.ID); ok {
//				if string(value) == data[s.field] {
//					return true, nil
//				}
//			} else {
//				return false, errors.Wrapf(errors.New("can't compare"), "type: %v", reflect.TypeOf(v))
//			}
//		}
//		return false, nil
//	case *multi:
//		if len(s.values) > 1 {
//			return boltMultiArgs(s.values, s.oper, data)
//		} else if len(s.values) == 1 {
//			return CheckMap(s.values[0], data)
//		} else if strings.ToLower(s.oper) == "or" {
//			return false, nil
//		} else {
//			return true, nil
//		}
//	case *not:
//		res, err := CheckMap(s, data)
//		if err != nil {
//			return false, err
//		}
//		return !res, nil
//	default:
//		return false, errors.Wrapf(ErrBadSelector, "%#v: %s", selector, reflect.TypeOf(selector))
//	}
//	return false, nil
//}
