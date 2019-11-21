package elements_index

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/applications/records"
)

func All(joinerOp joiner.Operator) ([]Item, common.Errors) {
	elementsOp := records.Operator(nil)

	var items []Item
	var errs common.Errors

	for _, component := range joinerOp.InterfacesAll(&elementsOp) {
		elementsOp, ok := component.Interface.(records.Operator)
		if ok {
			items = append(items, Item{component.InterfaceKey, elementsOp})
		} else {
			errs = append(errs, errors.Errorf("incorrect elements.Operator interface (%T) for key '%s'", component.Interface, component.InterfaceKey))
		}
	}

	return items, errs
}
