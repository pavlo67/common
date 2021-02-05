package selectors

type TermOneOf struct {
	Key    string
	Values []interface{}
}

func In(key string, values ...interface{}) *Term {
	return &Term{TermOneOf{key, values}, nil, Nop2}
}

type TermString struct {
	String string
	Values []interface{}
}

func String(str string, values ...interface{}) *Term {
	return &Term{TermString{str, values}, nil, Nop2}
}
