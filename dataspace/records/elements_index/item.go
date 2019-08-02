package elements_index

import (
	"github.com/pavlo67/constructor/basis"

	"github.com/pavlo67/constructor/dataspace/records"
	"github.com/pavlo67/constructor/starter/joiner"
)

type Item struct {
	Key      joiner.ComponentKey
	Operator records.Operator
}

func (item *Item) UnmarshalJSON([]byte) error {
	return basis.ErrNotImplemented
}

func (item *Item) MarshalJSON() ([]byte, error) {
	return nil, basis.ErrNotImplemented
}
