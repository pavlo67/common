package data

import (
	"github.com/pavlo67/workshop/common"

	"github.com/pavlo67/workshop/components/editor"
)

var _ editor.Operator = &Item{}

func (item Item) PrepareToEdit() ([]editor.Field, error) {
	//var fields []editor.Field
	//
	//// TODO!!!
	//
	//if len(item.Data.Content) > 0 {
	//	editDetailsOp, ok := item.Data.Content.(editor.Operator)
	//	if !ok {
	//		return nil, errors.Errorf("item.Details (%#v) isn't editor.ActorKey", item.Details)
	//	}
	//
	//	detailsFields, err := editDetailsOp.PrepareToEdit()
	//	if err != nil {
	//		return nil, errors.Wrapf(err, "can't .PrepareToEdit(%#v)", item.Details)
	//	}
	//
	//	fields = append(fields, detailsFields...)
	//}
	//
	//return fields, nil

	return nil, common.ErrNotImplemented
}

func (item *Item) SaveEdited([]editor.Field) error {
	return common.ErrNotImplemented
}
