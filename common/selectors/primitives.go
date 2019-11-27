package selectors

type TermOneOf struct {
	Key    string
	Values []interface{}
}

func In(key string, values ...interface{}) *TermUnary {
	return &TermUnary{TermOneOf{key, values}, nil}
}
