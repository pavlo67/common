package packs

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "packs"
const CollectionDefault = "packs"

const TaskAction crud.ActionKey = "task"

type Pack struct {
	Key identity.Key

	From    identity.Key `bson:",omitempty" json:",omitempty"`
	To      identity.Key `bson:",omitempty" json:",omitempty"`
	Options common.Map   `bson:",omitempty" json:",omitempty"`

	Data crud.Data `bson:",omitempty" json:",omitempty"`

	History crud.History `bson:",omitempty" json:",omitempty"`
}

type Item struct {
	ID   common.ID `bson:"_id,omitempty" json:",omitempty"`
	Pack `          bson:",inline" json:",inline"`
}

type Operator interface {
	Save(*Pack, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error
	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)

	AddHistory(common.ID, crud.History, *crud.SaveOptions) (crud.History, error)
}
