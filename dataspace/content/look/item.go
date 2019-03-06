package look

import (
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/dataspace/content"
	"github.com/pavlo67/punctum/dataspace/records"
	"github.com/pavlo67/punctum/dataspace/records/elements_index"
)

const Type content.Type = "look"

var _ content.Item = &Item{}

type Item struct {
	Operator elements_index.Item `bson:"operator"          json:"operator"`
	Options  crud.ReadOptions    `bson:"options,omitempty" json:"options,omitempty"`
	Items    []records.Item      `bson:"items,omitempty"   json:"items,omitempty"`
}

func (look Item) Type() content.Type {
	return Type
}

func (look Item) Key() string {
	return ""
}

func (look Item) Set(interface{}) error {
	return nil
}

func (look Item) Refresh() error {
	return nil
}

func (look Item) String() string {
	return ""
}
