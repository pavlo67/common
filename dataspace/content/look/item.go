package look

import (
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/dataspace/content"
	"github.com/pavlo67/punctum/dataspace/elements"
)

const Type content.Type = "look"

var _ content.Item = &Item{}

type Item struct {
	Options  crud.ReadOptions `bson:"options,omitempty"  json:"options,omitempty"`
	Elements []elements.Item  `bson:"elements,omitempty" json:"elements,omitempty"`
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
