package elements_index

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/punctum/dataspace/elements"
)

func All(joinerOp joiner.Operator) ([]Item, basis.Errors) {
	elementsOp := elements.Operator(nil)

	var items []Item
	var errs basis.Errors

	for _, component := range joinerOp.ComponentsAllWithInterface(&elementsOp) {
		elementsOp, ok := component.Interface.(elements.Operator)
		if ok {
			items = append(items, Item{component.Key, elementsOp})
		} else {
			errs = append(errs, errors.Errorf("incorrect elements.Operator interface (%T) for key '%s'", component.Interface, component.Key))
		}
	}

	return items, errs
}
