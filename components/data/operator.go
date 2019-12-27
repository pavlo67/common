package data

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/flow"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/tags"
)

const InterfaceKey joiner.InterfaceKey = "data"
const CollectionDefault = "data"

type TypeKey string

const TypeKeyString TypeKey = "string"
const TypeKeyHRefImage TypeKey = "href_image"
const TypeKeyHRef TypeKey = "href"

type Type struct {
	Key      TypeKey
	Exemplar interface{}
}

type Item struct {
	ID         common.Key   `bson:"_id,omitempty" json:",omitempty"`
	ExportID   string       `bson:",omitempty"    json:",omitempty"`
	URL        string       `bson:",omitempty"    json:",omitempty"`
	TypeKey    TypeKey      `bson:",omitempty"    json:",omitempty"`
	Title      string       `bson:",omitempty"    json:",omitempty"`
	Summary    string       `bson:",omitempty"    json:",omitempty"`
	Embedded   []Item       `bson:",omitempty"    json:",omitempty"`
	Tags       []tags.Item  `bson:",omitempty"    json:",omitempty"`
	Details    interface{}  `bson:"-"             json:",omitempty"`
	DetailsRaw []byte       `bson:",omitempty"    json:",omitempty"` // shouldn't be used directly
	Status     crud.History `bson:",omitempty"    json:",omitempty"`
	Origin     flow.Origin  `bson:",omitempty"    json:",omitempty"`
}

type Operator interface {
	Save([]Item, *crud.SaveOptions) ([]common.Key, error)
	Remove(common.Key, *crud.RemoveOptions) error

	Read(common.Key, *crud.GetOptions) (*Item, error)
	SetDetails(item *Item) error

	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	Count(*selectors.Term, *crud.GetOptions) (uint64, error)

	Export(afterID string, options *crud.GetOptions) ([]Item, error)
}

type Convertor interface {
	GetData() (*Item, error)
	SaveData(Item) error
}

// TODO: .History, etc...

//	Key        dataspace.Key `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Item   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Item content.Item `bson:"content,omitempty" json:"content,omitempty"`
//	Tags   links.Tags  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  common.Key `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner common.Key `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
