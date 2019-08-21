package elements_index

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/joiner"

	"github.com/pavlo67/constructor/applications/records"
)

func All(joinerOp joiner.Operator) ([]Item, common.Errors) {
	elementsOp := records.Operator(nil)

	var items []Item
	var errs common.Errors

	for _, component := range joinerOp.ComponentsAllWithInterface(&elementsOp) {
		elementsOp, ok := component.Interface.(records.Operator)
		if ok {
			items = append(items, Item{component.InterfaceKey, elementsOp})
		} else {
			errs = append(errs, errors.Errorf("incorrect elements.Operator interface (%T) for key '%s'", component.Interface, component.InterfaceKey))
		}
	}

	return items, errs
}
