package elements_index

import (
	"github.com/pavlo67/workshop/basis/common"

	"github.com/pavlo67/workshop/applications/records"
	"github.com/pavlo67/workshop/basis/joiner"
)

type Item struct {
	Key      joiner.InterfaceKey
	Operator records.Operator
}

func (item *Item) UnmarshalJSON([]byte) error {
	return common.ErrNotImplemented
}

func (item *Item) MarshalJSON() ([]byte, error) {
	return nil, common.ErrNotImplemented
}
