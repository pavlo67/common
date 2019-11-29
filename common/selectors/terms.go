package selectors

type Literal string
type Value struct{ V interface{} }

// unary terms -----------------------------------------------------------------------------------------

type TermUnary struct {
	Value interface{}
	OperationUnary
}

type OperationUnary rune

const Not OperationUnary = '!'
const Inv OperationUnary = '-'

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

const Add OperationBinary = '+'
const Sub OperationBinary = '-'
const Mult OperationBinary = '*'
const Div OperationBinary = '/'

const Gt OperationBinary = '>'
const Ge OperationBinary = 'g'
const Eq OperationBinary = '='
const Ne OperationBinary = 'n'
const Lt OperationBinary = '<'
const Le OperationBinary = 'l'

const And OperationBinary = 'A'
const Or OperationBinary = 'O'

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
