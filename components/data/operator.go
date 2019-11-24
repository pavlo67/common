package data

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/crud"
	"github.com/pavlo67/workshop/components/selector"
	"github.com/pavlo67/workshop/components/tagger"
)

const InterfaceKey joiner.InterfaceKey = "data"

type Item struct {
	ID       common.ID `bson:"_id,omitempty"`
	Title    string
	Summary  string
	URL      string
	Embedded []Item
	Tags     []tagger.Tag

	Status

	// Details should be used with Operator.Save only (and use Operator.Details to get .Details value)
	Details interface{} `bson:"-" json:"-"`

	// DetailsRaw shouldn't be used directly
	DetailsRaw []byte `bson:"Details" json:"Details"`
}

type Status struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Part struct {
	Key   []string
	Count uint64
}

type Operator interface {
	Save(Item, *crud.SaveOptions) (*common.ID, error)
	Read(common.ID, *crud.GetOptions) (*Item, error)
	Details(item *Item, exemplar interface{}) error

	Exists(*selector.Term, *crud.GetOptions) ([]Part, error)
	List(*selector.Term, *crud.GetOptions) ([]Item, error)
	Remove(*selector.Term, *crud.RemoveOptions) error
}

// TODO: .History, etc...

//	ID        dataspace.ID `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//	CreatedAt time.Time    `bson:"created_at,omitempty" json:"created_at"`
//	UpdatedAt *time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Item   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Item content.Item `bson:"content,omitempty" json:"content,omitempty"`
//	Tags   links.Tags  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  common.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner common.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`

//  Origin     `           bson:"origin,omitempty"      json:"origin,omitempty"`
//  OriginTime *time.Time `bson:"origin_time,omitempty" json:"origin_time,omitempty"`
//  OriginData string     `bson:"origin_data,omitempty" json:"origin_data,omitempty"`

//  type Origin struct {
//	  ID  common.ID `bson:"id,omitempty"  json:"id,omitempty"`
//	  Key string    `bson:"key,omitempty" json:"key,omitempty"`
//  }
