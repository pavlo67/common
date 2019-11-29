package logic

import "github.com/pavlo67/workshop/common/selectors"

func AND(value0, value1 interface{}) *selectors.Term {
	return &selectors.Term{selectors.TermUnary{Value: value0}, []selectors.TermNext{{selectors.TermUnary{Value: value1}, selectors.And}}}
}

func OR(value0, value1 interface{}) *selectors.Term {
	return &selectors.Term{selectors.TermUnary{Value: value0}, []selectors.TermNext{{selectors.TermUnary{Value: value1}, selectors.Or}}}
}

func NOT(value interface{}) *selectors.Term {
	return &selectors.Term{selectors.TermUnary{Value: value, OperationUnary: selectors.Not}, nil}
}
