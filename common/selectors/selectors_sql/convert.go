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
		case selectors.ADD:
			sqlCondition = sqlCondition + " + " + sqlConditionNext
		case selectors.SUB:
			sqlCondition = sqlCondition + " - " + sqlConditionNext
		case selectors.MULT:
			sqlCondition = sqlCondition + " * " + sqlConditionNext
		case selectors.DIV:
			sqlCondition = sqlCondition + " / " + sqlConditionNext
		case selectors.GT:
			sqlCondition = sqlCondition + " > " + sqlConditionNext
		case selectors.GE:
			sqlCondition = sqlCondition + " >= " + sqlConditionNext
		case selectors.EQ:
			sqlCondition = sqlCondition + " = " + sqlConditionNext
		case selectors.NE:
			sqlCondition = sqlCondition + " <> " + sqlConditionNext // !=
		case selectors.LT:
			sqlCondition = sqlCondition + " < " + sqlConditionNext
		case selectors.LE:
			sqlCondition = sqlCondition + " <= " + sqlConditionNext
		case selectors.AND:
			sqlCondition = sqlCondition + " AND " + sqlConditionNext
		case selectors.OR:
			sqlCondition = sqlCondition + " OR " + sqlConditionNext
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
	default:
		// TODO: check if value is suitable for SQL
		return "?", []interface{}{value}, nil
	}

	sqlCondition, values, err = use(termUnary.Value)
	if err != nil {
		return "", nil, errors.Wrapf(err, "on selectors_sql.use(%#v)", termUnary.Value)
	}
	sqlCondition = "(" + sqlCondition + ")"

	switch termUnary.OperationUnary {
	case selectors.NOT:
		return "NOT " + sqlCondition, values, nil
	case selectors.INV:
		return "-" + sqlCondition, values, nil
	}

	return "", nil, errors.Errorf("wrong .OperationUnary on selectors_sql.use(%#v)", termUnary)
}
