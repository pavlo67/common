package selector

import "errors"

//var ErrBadObjectDetails = errors.New("bad object details")
//var ErrBadObjectManagers = errors.New("bad object managers")
//var ErrBadDataBuffer = errors.New("bad data buffer")
//var ErrEmptySelector = errors.New("empty selector")

var ErrBadSelector = errors.New("bad selector")

type Term struct {
	First TermUnary
	Next  []TermNext
}

type TermNext struct {
	TermUnary
	Operation
}

type Operation rune

const ADD Operation = '+'
const SUB Operation = '-'
const MULT Operation = '*'
const DIV Operation = '/'

const GT Operation = '>'
const GE Operation = 'g'
const EQ Operation = '='
const NE Operation = 'n'
const LT Operation = '<'
const LE Operation = 'l'

const AND Operation = 'A'
const OR Operation = 'O'

// unary terms -----------------------------------------------------------------------------------------

type TermUnary struct {
	Value interface{}
	*OperationUnary
}

type OperationUnary rune

const NOT OperationUnary = '!'
const MIN OperationUnary = '-'
const INV OperationUnary = '/'

// one-of terms ----------------------------------------------------------------------------------------

type TermOneOfStr struct {
	Key    string
	Values []string
}

// helpers ---------------------------------------------------------------------------------------------

func Unary(term *Term) *TermUnary {
	return &TermUnary{term, nil}
}

func InStrUnary(key string, values ...string) *TermUnary {
	return &TermUnary{TermOneOfStr{key, values}, nil}
}

func InStr(key string, values ...string) *Term {
	return &Term{First: TermUnary{TermOneOfStr{key, values}, nil}}
}

func Le(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func Lt(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func Eq(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func Ne(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func Ge(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func Gt(key, value interface{}) *Term {
	return &Term{
		TermUnary{key, nil},
		[]TermNext{{TermUnary{value, nil}, LE}},
	}
}

func And(termsUnary0 ...*TermUnary) *Term {
	var termsUnary []*TermUnary
	for _, termUnary0 := range termsUnary0 {
		if termUnary0 != nil {
			termsUnary = append(termsUnary, termUnary0)
		}
	}

	if len(termsUnary) < 1 {
		return nil
	}

	term := Term{
		First: *termsUnary[0],
	}
	for _, termUnary := range termsUnary[1:] {
		term.Next = append(term.Next, TermNext{*termUnary, AND})
	}

	return &term
}
