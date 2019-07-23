package structure

import (
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/dataspace"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "structure"
const CleanerInterfaceKey joiner.InterfaceKey = "structure.cleaner"

type Item struct {
}

//	ID        dataspace.ID `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//	CreatedAt time.Time    `bson:"created_at,omitempty" json:"created_at"`
//	UpdatedAt *time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Brief   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Content content.Item `bson:"content,omitempty" json:"content,omitempty"`
//	Links   links.Links  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  auth.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner auth.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`

type Operator interface {
	Add(Item) (dataspace.ID, error)

	Remove(dataspace.ID) error

	List(crud.ReadOptions) ([]Item, error)
}
