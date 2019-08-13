package elements_index

import (
	"github.com/pavlo67/constructor/components/basis"

	"github.com/pavlo67/constructor/applications/records"
	"github.com/pavlo67/constructor/components/basis/joiner"
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
