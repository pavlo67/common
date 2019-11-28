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

type Item struct {
	ID       common.ID `bson:"_id,omitempty"`
	URL      string
	Title    string
	Summary  string
	Embedded []Item
	Tags     []tagger.Tag

	// Details should be used with Operator.Save only (and use Operator.Details to get .Details value)
	Details interface{} `bson:"-" json:"-"`

	// DetailsRaw shouldn't be used directly
	DetailsRaw []byte

	crud.Status
	flow.Origin
}

type Operator interface {
	Save([]Item, *crud.SaveOptions) ([]common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error

	Read(common.ID, *crud.GetOptions) (*Item, error)
	Details(item *Item, exemplar interface{}) error
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	Count(*selectors.Term, *crud.GetOptions) ([]crud.Counter, error)
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
