package data

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/constructions/editor"
)

var _ editor.Operator = &Item{}

func (item Item) PrepareToEdit() ([]editor.Field, error) {
	var fields []editor.Field

	// TODO!!!

	if item.Details != nil {
		editDetailsOp, ok := item.Details.(editor.Operator)
		if !ok {
			return nil, errors.Errorf("item.Details (%#v) isn't editor.Operator", item.Details)
		}

		detailsFields, err := editDetailsOp.PrepareToEdit()
		if err != nil {
			return nil, errors.Wrapf(err, "can't .PrepareToEdit(%#v)", item.Details)
		}

		fields = append(fields, detailsFields...)
	}

	return fields, nil
}

func (item *Item) SaveEdited([]editor.Field) error {
	return common.ErrNotImplemented
}
