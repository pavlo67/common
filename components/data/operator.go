package data

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/tagger"
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
	ID       common.ID    `bson:"_id,omitempty" json:",omitempty"`
	ExportID string       `bson:",omitempty"    json:",omitempty"`
	URL      string       `bson:",omitempty"    json:",omitempty"`
	TypeKey  TypeKey      `bson:",omitempty"    json:",omitempty"`
	Title    string       `bson:",omitempty"    json:",omitempty"`
	Summary  string       `bson:",omitempty"    json:",omitempty"`
	Embedded []Item       `bson:",omitempty"    json:",omitempty"`
	Tags     []tagger.Tag `bson:",omitempty"    json:",omitempty"`

	// Details should be used with Operator.Save only (and use Operator.Details to get .Details value)
	Details interface{} `bson:"-" json:",omitempty"`

	// DetailsRaw shouldn't be used directly
	DetailsRaw []byte `bson:",omitempty"    json:",omitempty"`

	Status crud.Status `bson:",omitempty" json:",omitempty"`
	Origin flow.Origin `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	Save([]Item, *crud.SaveOptions) ([]common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error

	Read(common.ID, *crud.GetOptions) (*Item, error)
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

//	ID        dataspace.ID `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Item   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Item content.Item `bson:"content,omitempty" json:"content,omitempty"`
//	Tags   links.Tags  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  common.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner common.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
