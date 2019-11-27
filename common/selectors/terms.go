package selectors

// unary terms -----------------------------------------------------------------------------------------

type TermUnary struct {
	Value interface{}
	OperationUnary
}

type OperationUnary rune

const NOT OperationUnary = '!'
const INV OperationUnary = '-'

// general terms --------------------------------------------------------------------------------------

type Term struct {
	Value interface{}
	Next  []TermNext
}

type TermNext struct {
	Value interface{}
	OperationBinary
}

type OperationBinary rune

const ADD OperationBinary = '+'
const SUB OperationBinary = '-'
const MULT OperationBinary = '*'
const DIV OperationBinary = '/'

const GT OperationBinary = '>'
const GE OperationBinary = 'g'
const EQ OperationBinary = '='
const NE OperationBinary = 'n'
const LT OperationBinary = '<'
const LE OperationBinary = 'l'

const AND OperationBinary = 'A'
const OR OperationBinary = 'O'

func TermBinary(operationBinary OperationBinary, value0, value1 interface{}) *Term {
	var term0 *TermUnary
	switch v := value0.(type) {
	case *TermUnary:
		term0 = v
	case TermUnary:
		term0 = &v
	default:
		term0 = &TermUnary{Value: value0}
	}

	var term1 *TermUnary
	switch v := value1.(type) {
	case *TermUnary:
		term1 = v
	case TermUnary:
		term1 = &v
	default:
		term1 = &TermUnary{Value: value1}
	}

	return &Term{*term0, []TermNext{{*term1, operationBinary}}}
}

func TermMultiple(operationBinary OperationBinary, values ...interface{}) *Term {
	if len(values) < 1 {
		return nil
	}

	var term *Term
	switch v := values[0].(type) {
	case *TermUnary:
		term = &Term{Value: *v}
	case TermUnary:
		term = &Term{Value: v}
	default:
		term = &Term{Value: TermUnary{Value: values[0]}}
	}

	for _, value := range values[1:] {
		switch v := value.(type) {
		case *TermUnary:
			term.Next = append(term.Next, TermNext{*v, operationBinary})
		case TermUnary:
			term.Next = append(term.Next, TermNext{v, operationBinary})
		default:
			term.Next = append(term.Next, TermNext{TermUnary{Value: value}, operationBinary})
		}
	}

	return term
}
