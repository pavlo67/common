package data

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/components/tagger"
)

const InterfaceKey joiner.InterfaceKey = "data"
const CleanerInterfaceKey joiner.InterfaceKey = "datacleaner"

const CollectionDefault = "data"

type Item struct {
	ID  common.ID    `bson:"_id,omitempty" json:",omitempty"`
	Key identity.Key `bson:",omitempty"    json:",omitempty"`

	URL      string       `bson:",omitempty"    json:",omitempty"`
	Title    string       `bson:",omitempty"    json:",omitempty"`
	Summary  string       `bson:",omitempty"    json:",omitempty"`
	Embedded []Item       `bson:",omitempty"    json:",omitempty"`
	Tags     []tagger.Tag `bson:",omitempty"    json:",omitempty"`
	Data     crud.Data    `bson:",omitempty"    json:",omitempty"`

	History crud.History `bson:",omitempty"    json:",omitempty"`
}

type Operator interface {
	Save([]Item, *crud.SaveOptions) ([]common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error

	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	Count(*selectors.Term, *crud.GetOptions) (uint64, error)

	Export(afterID string, options *crud.GetOptions) ([]Item, error)
}

type Convertor interface {
	GetData() (*Item, error)
	SaveData(Item) error
}

// TODO: .History, etc...

//	ID        dataspace.ID `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Tag   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Tag content.Tag `bson:"content,omitempty" json:"content,omitempty"`
//	Tags   links.Tags  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  common.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner common.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
