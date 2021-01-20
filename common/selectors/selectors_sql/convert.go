package selectors_sql

import (
	"fmt"

	"github.com/pavlo67/workshop/common/errors"

	"strings"

	"github.com/pavlo67/workshop/common/selectors"
)

func Use(term *selectors.Term) (sqlCondition string, values []interface{}, err error) {
	if term == nil {
		return "", nil, nil
	}

	if term.Right == nil {
		return use(term.Left)
	}

	sqlCondition, values, err = use(term.Left)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", term.Left)
	}

	sqlConditionRight, valuesNext, err := use(term.Right)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", term.Right)
	}

	//sqlCondition = "(" + sqlCondition + ")"
	//sqlConditionRight = "(" + sqlConditionRight + ")"

	switch term.Operation {
	case selectors.Add:
		sqlCondition = sqlCondition + " + " + sqlConditionRight
	case selectors.Sub:
		sqlCondition = sqlCondition + " - " + sqlConditionRight
	case selectors.Mult:
		sqlCondition = sqlCondition + " * " + sqlConditionRight
	case selectors.Div:
		sqlCondition = sqlCondition + " / " + sqlConditionRight
	case selectors.Gt:
		sqlCondition = sqlCondition + " > " + sqlConditionRight
	case selectors.Ge:
		sqlCondition = sqlCondition + " >= " + sqlConditionRight
	case selectors.Eq:
		sqlCondition = sqlCondition + " = " + sqlConditionRight
	case selectors.Ne:
		sqlCondition = sqlCondition + " <> " + sqlConditionRight // !=
	case selectors.Lt:
		sqlCondition = sqlCondition + " < " + sqlConditionRight
	case selectors.Le:
		sqlCondition = sqlCondition + " <= " + sqlConditionRight
	case selectors.And:
		sqlCondition = sqlCondition + " AND " + sqlConditionRight
	case selectors.Or:
		sqlCondition = sqlCondition + " OR " + sqlConditionRight
	default:
		return "", nil, fmt.Errorf("wrong .Operation on selectors_sql.use(%#v)", term.Right)
	}

	values = append(values, valuesNext...)

	return "(" + sqlCondition + ")", values, nil
}

func use(value interface{}) (sqlCondition string, values []interface{}, err error) {
	var termUnary *selectors.TermUnary

	switch v := value.(type) {
	case selectors.Term:
		return Use(&v)
	case *selectors.Term:
		return Use(v)
	case selectors.TermUnary:
		termUnary = &v
	case *selectors.TermUnary:
		termUnary = v
	case string:
		return v, nil, nil
	case selectors.Value:
		return "?", []interface{}{v.V}, nil
	case selectors.TermOneOf:
		if len(v.Values) < 1 {
			// TODO!!! is it correct?
			return "", nil, nil
		}
		return v.Key + " in (" + strings.Repeat(",?", len(v.Values))[1:] + ")", v.Values, nil
	case *selectors.TermOneOf:
		if v == nil || len(v.Values) < 1 {
			// TODO!!! is it correct?
			return "", nil, nil
		}
		return v.Key + " in (" + strings.Repeat(",?", len(v.Values))[1:] + ")", v.Values, nil
	case selectors.TermString:
		return v.String, v.Values, nil
	case *selectors.TermString:
		return v.String, v.Values, nil
	default:
		return "", nil, fmt.Errorf("wrong value for selectors_sql.use(%#v)", value)
	}

	if termUnary.OperationUnary == selectors.NopUn {
		return use(termUnary.ValueUnary)
	}

	sqlCondition, values, err = use(termUnary.ValueUnary)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", termUnary.ValueUnary)
	}

	// sqlCondition = "(" + sqlCondition + ")"

	switch termUnary.OperationUnary {
	case selectors.Not:
		return "NOT " + sqlCondition, values, nil
	case selectors.Inv:
		return "-" + sqlCondition, values, nil
	}

	return "", nil, fmt.Errorf("wrong .OperationUnary on selectors_sql.use(%#v)", termUnary)
}
