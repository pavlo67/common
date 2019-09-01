package scriptor

import "github.com/pkg/errors"

type Item struct {
	Sequence

	stack    []Action
	prefixes []Action
}

func (item *Item) ToStack(a *Action) error {
	if len(item.stack) > len(item.prefixes) {
		return errors.New("too long sequence without prefixes...")
	}

	if a == nil {
		item.stack = append(item.stack, Action{TypeNil, nil})
	} else {
		item.stack = append(item.stack, *a)
	}

	return nil
}

func (item *Item) ToActions(a Action, constants Values) error {
	// TODO!!! check if prefix only

	if len(item.prefixes) >= len(item.stack) {
		return errors.New("too many prefixes for the .Sequence...")
	}

	if false {
		// TODO: check prefixes priority
		if err := item.Prepare(constants); err != nil {
			return err
		}

	}

	item.prefixes = append(item.prefixes, a)

	return nil
}

func (item *Item) Prepare(constants Values) error {
	//if len(item.prefixes) > 0 {
	//	return &item, s0, errors.Errorf("open prefixes (%#v) remain: %s", item.prefixes, sOriginal[:offset+len(s0)])


	return nil
}

func (item *Item) Run(variables Variables) error {
	return nil
}
