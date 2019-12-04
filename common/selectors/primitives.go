package selectors

type TermOneOf struct {
	Key    string
	Values []interface{}
}

func In(key string, values ...interface{}) *Term {
	return &Term{TermOneOf{key, values}, nil, Nop}
}
