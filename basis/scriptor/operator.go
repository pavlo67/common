package scriptor

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type Type byte
type Value interface{}
type Element struct {
	Type
	Value
}

const (
	TypeNil Type = iota
	TypeByte
	TypeInt
	TypeFloat
	TypeString
	TypePrefix
	TypeInfix
	TypePostfix
	TypePostfix2
	TypeSequence
)

type Values map[string]Element
type Sequence []Element

type Variables struct {
	Values
	Top *Variables
}

var reOpenItem = regexp.MustCompile("^[([{]")
var reCloseItem = regexp.MustCompile("^[)\\]}]")
var reWord = regexp.MustCompile("^[a-zA-Z_]\\w*")
var reInteger = regexp.MustCompile("^-?\\d+")
var reFloat = regexp.MustCompile("^-?(\\d+\\.\\d*|\\.\\d+)")
var reSpace = regexp.MustCompile("^\\s+")

const openBr = "("
const closeBr = ")"

var itemPairs = map[string]string{
	openBr: closeBr,
	"[":    "]",
	"{":    "}",
}

func Read(sOriginal string, openedWith string, constants Values) (action *Element, rest string, err error) {
	item := Item{}
	if constants == nil {
		constants = Values{}
	}

	s := sOriginal
	offset := 0

	for {

		s = reSpace.ReplaceAllString(s, "")
		if s == "" {
			break
		}

		if s0 := reWord.FindString(s); s0 != "" {
			if err := item.ToStack(&Element{TypeString, s0}); err != nil {
				return nil, s, err // TODO!!! show details
			}

			offset += len(s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reFloat.FindString(s); s0 != "" {
			f0, _ := strconv.ParseFloat(s0, 64)
			if err := item.ToStack(&Element{TypeFloat, f0}); err != nil {
				return nil, s, err // TODO!!! show details
			}

			offset += len(s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reInteger.FindString(s); s0 != "" {
			i0, _ := strconv.ParseInt(s0, 10, 64)
			if err := item.ToStack(&Element{TypeInt, i0}); err != nil {
				return nil, s, err // TODO!!! show details
			}

			offset += len(s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reInfix.FindString(s); s0 != "" {

			if err := item.ToInfixes(s0, constants); err != nil {
				return nil, s, err // TODO!!! show details errors.Errorf("open infixes (%#v) remain: %s", item.infixes, sOriginal[:offset+len(s0)])
			}

			offset += len(s0)
			s = s[len(s0):]
			continue
		}

		if s0 := reOpenItem.FindString(s); s0 != "" {
			value, s1, err := Read(s[len(s0):], s0, constants)
			if err != nil {
				return nil, s, nil
			}
			if err := item.ToStack(value); err != nil {
				return nil, s, err // TODO!!! show details
			}

			offset += len(s) - len(s1)
			s = s1
			continue
		}

		if s0 := reCloseItem.FindString(s); s0 != "" {

			// original string closed with some bracket ----------------------------------------------------

			if itemPairs[openedWith] != s0 {
				return nil, s, errors.Errorf("wrong close bracket: %s", openedWith+sOriginal[:offset+len(s0)])
			}
			if err := item.PrepareInfixesAll(constants); err != nil {
				return nil, s, err
			}
			if len(item.stack) > 1 {
				return nil, s, errors.Errorf("open stack (%#v) remain: %s", item.stack, sOriginal[:offset+len(s0)])
			} else if len(item.stack) == 1 {
				if openedWith != openBr {
					return nil, s, errors.Errorf("open stack (%#v) remain: %s", item.stack, sOriginal[:offset+len(s0)])
				} else if len(item.Sequence) > 0 {
					return nil, s, errors.Errorf("open stack (%#v) remain: %s", item.stack, sOriginal[:offset+len(s0)])
				}
				return &item.stack[0], s[len(s0):], nil
			}
			if openedWith != openBr {
				return nil, s[len(s0):], nil
			}
			return &Element{TypeSequence, item.Sequence}, s[len(s0):], nil

			// /original string closed with some bracket ---------------------------------------------------

		}

		return nil, s, errors.Errorf("wrong symbol: %s", s)
	}

	// original string finished -----------------------------------------------------------------

	if openedWith != "" {
		return nil, s, errors.Errorf("no close bracket: %s", openedWith+sOriginal)
	}
	if err := item.PrepareInfixesAll(constants); err != nil {
		return nil, "", err
	}
	if len(item.stack) > 0 {
		return nil, s, errors.Errorf("open stack (%#v) remain: %s", item.stack, sOriginal)
	}
	return &Element{TypeSequence, item.Sequence}, "", nil

	// /original string finished ----------------------------------------------------------------

}
