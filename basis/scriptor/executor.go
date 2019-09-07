package scriptor

import (
	"github.com/pkg/errors"
)

type Item struct {
	stack   []Element
	infixes []Infix

	Sequence []Element
}

func (item *Item) ToStack(a *Element) error {
	if len(item.stack) > len(item.infixes) {
		return errors.Errorf("too long stack (%#v) vs. infixes (%#v), adding %#v", item.stack, item.infixes, a)
	}

	if a == nil {
		item.stack = append(item.stack, Element{TypeNil, nil})
	} else {
		item.stack = append(item.stack, *a)
	}

	return nil
}

func (item *Item) ToInfixes(sign string, constants Values) error {
	if len(item.infixes) >= len(item.stack) {
		return errors.New("too many infixes for .ToInfixes()...")
	}

	infixes, ok := Infixes[sign]
	if !ok || len(infixes) < 1 {
		return errors.Errorf("no infix for sign '%s'", sign)
	}
	priority := infixes[0].Priority

	for len(item.infixes) > 0 && item.infixes[len(item.infixes)-1].Priority >= priority {
		if err := item.PrepareInfix(constants); err != nil {
			return err
		}
	}

	typeLeft := item.stack[len(item.stack)-1].Type

	var infixPtr *Infix
	for _, infix := range infixes {
		if infix.Signatura[0] == typeLeft {
			infixPtr = &infix
			break
		}
	}

	if infixPtr == nil {
		return errors.Errorf("no infix for sign '%s' and type %d", sign, typeLeft)
	}

	item.infixes = append(item.infixes, *infixPtr)

	return nil
}

func (item *Item) PrepareInfix(constants Values) error {
	leftOp := item.stack[len(item.stack)-2]
	rightOp := item.stack[len(item.stack)-1]
	infix := item.infixes[len(item.infixes)-1]

	if rightOp.Type != infix.Signatura[1] {
		return errors.Errorf("wrong right operand type (%d) for infix %#v", rightOp.Type, infix)
	}

	//log.Printf("stack: %#v", item.stack)
	//log.Printf("infixes: %#v", item.infixes)

	item.stack = append(
		item.stack[:len(item.stack)-2],
		Element{infix.Signatura[2], Executor{[]interface{}{leftOp.Value, rightOp.Value, infix.Func2}}},
		//Element{infix.Signatura[2], infix.Func2(leftOp.Value, rightOp.Value), false},
	)
	item.infixes = item.infixes[:len(item.infixes)-1]

	//log.Printf("stack: %#v", item.stack)
	//log.Printf("infixes: %#v", item.infixes)

	return nil
}

func (item *Item) PrepareInfixesAll(constants Values) error {
	if len(item.infixes)+1 != len(item.stack) {
		return errors.Errorf("len(item.infixes) + 1 != len(item.stack): %d, %d", len(item.infixes)+1, len(item.stack))
	}
	for len(item.infixes) > 0 {
		if err := item.PrepareInfix(constants); err != nil {
			return err
		}

	}

	return nil
}

func (item *Item) Run(variables Variables) error {
	return nil
}
