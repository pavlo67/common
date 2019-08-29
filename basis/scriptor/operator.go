package scriptor

import "regexp"

type Action byte

type Item struct {
	Action
	Values   map[string]interface{}
	Sequence []interface{}
}

var reSymbol = regexp.MustCompile("^[[\\](){}:]") // :=|[?]
var reWord = regexp.MustCompile("^[a-zA-Z_]\\w*")
var reInteger = regexp.MustCompile("^-?\\d+")
var reFloat = regexp.MustCompile("^-?(\\d+\\.\\d*|\\.\\d+)")
var reSpace = regexp.MustCompile("^\\s+")

func Read(s string) (*Item, error) {
	var item Item

	for {

		s = reSpace.ReplaceAllString(s, "")
		if s == "" {
			break
		}

		if s0 := reWord.FindString(s); s0 != "" {
			item.Sequence = append(item.Sequence, s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reFloat.FindString(s); s0 != "" {
			item.Sequence = append(item.Sequence, s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reInteger.FindString(s); s0 != "" {
			item.Sequence = append(item.Sequence, s0)
			s = s[len(s0):]
			continue
		}

		// reSymbol

		item.Sequence = append(item.Sequence, s[:1])
		s = s[1:]

	}

	return &item, nil
}
