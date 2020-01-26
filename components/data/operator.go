package data

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/tagger"
)

const InterfaceKey joiner.InterfaceKey = "data"
const CleanerInterfaceKey joiner.InterfaceKey = "datacleaner"

const ItemsTypeKey crud.TypeKey = "data_items"

const CollectionDefault = "data"

type Item struct {
	ID  common.ID    `bson:"_id,omitempty" json:",omitempty"`
	Key identity.Key `bson:",omitempty"    json:",omitempty"`

	URL      string       `bson:",omitempty" json:",omitempty"`
	Title    string       `bson:",omitempty" json:",omitempty"`
	Summary  string       `bson:",omitempty" json:",omitempty"`
	Embedded []Item       `bson:",omitempty" json:",omitempty"`
	Tags     []tagger.Tag `bson:",omitempty" json:",omitempty"`
	Data     crud.Data    `bson:",omitempty" json:",omitempty"`

	OwnerKey  identity.Key `bson:",omitempty" json:",omitempty"`
	ViewerKey identity.Key `bson:",omitempty" json:",omitempty"`

	History crud.History `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	Save(Item, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error

	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	Count(*selectors.Term, *crud.GetOptions) (uint64, error)
}

type Convertor interface {
	GetData() (*Item, error)
	SaveData(Item) error
}
