package selectors

import (
	"strings"

	"github.com/pavlo67/punctum/auth"
	"github.com/pkg/errors"
)

// !!! this package requires on auth.IDentity & basis.UserIS types due to its application usage

// TODO: simplify result condition looking through down for nil values

// TODO:
//for k, v := range values {
//   if value, ok := v.(basis.UserIS); ok {
//      values[k] = string(value)
//   }
//}

func Mysql(userIS auth.ID, selector Selector) (string, []interface{}, error) {
	if selector == nil {
		return "", NoValues, nil
	}

	switch s := selector.(type) {
	case *multi:
		args := []Selector{}
		for _, arg := range s.values {
			if arg == nil {
				continue
			}
			if term, ok := arg.(Selector); ok {
				args = append(args, term)
				continue
			}
			return "", NoValues, errors.Wrapf(ErrBadSelector, "%+v in %+v", arg, selector)
		}
		if len(args) > 1 {
			return mysqlMultiArgs(userIS, args, s.oper)
		} else if len(args) == 1 {
			return Mysql(userIS, args[0])
		} else if strings.ToLower(s.oper) == "or" {
			return "0", NoValues, nil
		} else {
			return "1", NoValues, nil
		}

	case *not:
		condition, values, err := Mysql(userIS, s.term)
		if err != nil {
			return "", nil, err
		}
		return "NOT " + condition, values, nil

	case *in:
		if len(s.values) < 1 {
			return "0", NoValues, nil
		}

		for k, v := range s.values {
			if value, ok := v.(auth.ID); ok {
				s.values[k] = string(value)
			}
		}

		if len(s.values) > 1 {
			return "`" + s.field + "` in (" + wildcards(len(s.values)) + ")", correct(s.values), nil
		} else if s.values[0] != nil {
			return "`" + s.field + "` = ?", correct(s.values), nil
		} else {
			return "`" + s.field + "` UserIS NULL", nil, nil
		}

	case *op:
		return "`" + s.field + "` " + s.op + " ?", correct([]interface{}{s.value}), nil

	case *match:
		if len(s.value) < 1 {
			return "1", NoValues, nil
		} else if len(s.field) < 1 {
			return "0", NoValues, nil
		} else if strings.ToLower(s.mode) == "like" {
			return s.field + " like " + s.mode + " ?", []interface{}{s.value}, nil
		} else {
			return "MATCH (" + s.field + ") AGAINST (? " + s.mode + ")", []interface{}{s.value}, nil
		}

	case PrepareIDs:
		ids, err := s(userIS)
		if err != nil {
			return "", NoValues, errors.Wrapf(ErrBadSelector, "%#s", selector)
		}
		return Mysql(userIS, FieldStr("id", ids...))

	default:
		return "", NoValues, errors.Wrapf(ErrBadSelector, "%#s", selector)
	}
}

func correct(values []interface{}) []interface{} {
	var corrected []interface{}
	for _, v := range values {
		if is, ok := v.(auth.ID); ok {
			corrected = append(corrected, string(is))
		} else {
			corrected = append(corrected, v)
		}
	}

	return corrected
}

func wildcards(length int) string {
	result := make([]string, length)
	for i := range result {
		result[i] = "?"
	}

	return strings.Join(result, ",")
}

func mysqlMultiArgs(userIS auth.ID, args []Selector, oper string) (condition string, values []interface{}, err error) {
	conditions := []string{}
	for _, arg := range args {
		cond, val, err := Mysql(userIS, arg)
		if err != nil {
			return "", NoValues, err
		}
		conditions = append(conditions, "("+cond+")")
		values = append(values, val...)
	}

	return strings.Join(conditions, " "+oper+" "), values, nil
}

// example for tests
//
// {
//	 "$and"
//   [
//      {
//         "$or",
//         [
//            {
//               "b",
//               [4]
//            },
//            {
//               "c",
//               [1, 2]
//            },
//         ]
//      },
//      {
//         "a",
//         [1, 2, 3]
//      },
//   ]
// }
//
// -->
//
// ((`b` = ?) or (`c` in (?,?))) and (`a` in ?,?,?),
// [4, 1, 2, 1, 2, 3],
// nil
