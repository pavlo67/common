package selectors_sql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/selectors"
)

func Use(term *selectors.Term) (sqlCondition string, values []interface{}, err error) {
	if term == nil {
		return "", nil, nil
	}

	sqlCondition, values, err = use(term.Value)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", term.Value)
	}

	for _, t := range term.Next {
		sqlConditionNext, valuesNext, err := use(t.Value)
		if err != nil {
			return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", t)
		}
		sqlCondition = "(" + sqlCondition + ")"
		sqlConditionNext = "(" + sqlConditionNext + ")"

		switch t.OperationBinary {
		case selectors.Add:
			sqlCondition = sqlCondition + " + " + sqlConditionNext
		case selectors.Sub:
			sqlCondition = sqlCondition + " - " + sqlConditionNext
		case selectors.Mult:
			sqlCondition = sqlCondition + " * " + sqlConditionNext
		case selectors.Div:
			sqlCondition = sqlCondition + " / " + sqlConditionNext
		case selectors.Gt:
			sqlCondition = sqlCondition + " > " + sqlConditionNext
		case selectors.Ge:
			sqlCondition = sqlCondition + " >= " + sqlConditionNext
		case selectors.Eq:
			sqlCondition = sqlCondition + " = " + sqlConditionNext
		case selectors.Ne:
			sqlCondition = sqlCondition + " <> " + sqlConditionNext // !=
		case selectors.Lt:
			sqlCondition = sqlCondition + " < " + sqlConditionNext
		case selectors.Le:
			sqlCondition = sqlCondition + " <= " + sqlConditionNext
		case selectors.And:
			sqlCondition = sqlCondition + " And " + sqlConditionNext
		case selectors.Or:
			sqlCondition = sqlCondition + " Or " + sqlConditionNext
		default:
			return "", nil, errors.Errorf("wrong .OperationBinary on selectors_sql.use(%#v)", t)
		}

		values = append(values, valuesNext...)
	}

	return sqlCondition, values, nil
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
	case selectors.Literal:
		return string(v), nil, nil
	case selectors.Value:
		return "?", []interface{}{v.V}, nil
	default:
		return "", nil, errors.Errorf("wrong value for selectors_sql.use(%#v)", value)
	}

	sqlCondition, values, err = use(termUnary.Value)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", termUnary.Value)
	}
	sqlCondition = "(" + sqlCondition + ")"

	switch termUnary.OperationUnary {
	case selectors.Not:
		return "NOT " + sqlCondition, values, nil
	case selectors.Inv:
		return "-" + sqlCondition, values, nil
	}

	return "", nil, errors.Errorf("wrong .OperationUnary on selectors_sql.use(%#v)", termUnary)
}
