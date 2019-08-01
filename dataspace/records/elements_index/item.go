package elements_index

import (
	"github.com/pavlo67/associatio/basis"

	"github.com/pavlo67/associatio/dataspace/records"
	"github.com/pavlo67/associatio/starter/joiner"
)

type Item struct {
	Key      joiner.InterfaceKey
	Operator records.Operator
}

func (item *Item) UnmarshalJSON([]byte) error {
	return basis.ErrNotImplemented
}

func (item *Item) MarshalJSON() ([]byte, error) {
	return nil, basis.ErrNotImplemented
}
