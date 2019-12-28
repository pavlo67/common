package selectors

// type ID string
type Value struct {
	V interface{} `bson:",omitempty"    json:",omitempty"`
}

// unary terms -----------------------------------------------------------------------------------------

type TermUnary struct {
	ValueUnary     interface{} `bson:",omitempty"    json:",omitempty"`
	OperationUnary `            bson:",omitempty"    json:",omitempty"`
}

type OperationUnary rune

const NopUn OperationUnary = 0
const Not OperationUnary = '!'
const Inv OperationUnary = '-'

// general terms --------------------------------------------------------------------------------------

type Term struct {
	Left      interface{} `bson:",omitempty"    json:",omitempty"`
	Right     interface{} `bson:",omitempty"    json:",omitempty"`
	Operation `            bson:",omitempty"    json:",omitempty"`
}

type Operation rune

const Add Operation = '+'
const Sub Operation = '-'
const Mult Operation = '*'
const Div Operation = '/'

const Gt Operation = '>'
const Ge Operation = 'g'
const Eq Operation = '='
const Ne Operation = 'n'
const Lt Operation = '<'
const Le Operation = 'l'

const And Operation = 'A'
const Or Operation = 'O'

const Nop Operation = 0

func Operand(value interface{}) interface{} {
	switch v := value.(type) {
	case *TermUnary:
		if v != nil && v.OperationUnary == NopUn {
			return v.ValueUnary
		}
	case TermUnary:
		if v.OperationUnary == NopUn {
			return v.ValueUnary
		}
		return &v
	}
	return value
}

func Binary(operationBinary Operation, value0, value1 interface{}) *Term {
	return &Term{Operand(value0), Operand(value1), operationBinary}
}

//func Multiple(operationBinary Operation, values ...interface{}) *Term {
//	if len(values) < 1 {
//		return nil
//	}
//
//	term := &Term{Left: Operand(values[0])}
//
//	for _, value := range values[1:] {
//		term.Right = append(term.Right, TermRight{Operand(value), operationBinary})
//	}
//
//	return term
//}
