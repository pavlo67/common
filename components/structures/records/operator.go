package records

import (
	"time"

	"github.com/pavlo67/workshop/basis/joiner"

	"github.com/pavlo67/workshop/applications/links"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/vcs"
	"github.com/pavlo67/workshop/dataspace"
	"github.com/pavlo67/workshop/dataspace/content"
)

const InterfaceKey joiner.InterfaceKey = "records"
const CleanerInterfaceKey joiner.InterfaceKey = "records.cleaner"

// const GenusDefault = "note"
// const GenusFieldName = "genus"

//type Asked struct {
//	ID    string `bson:"_id,omitempty" json:"id"`
//	Genus string `bson:"genus"         json:"genus"`
//	Title  string `bson:"name"          json:"name"`
//}

type Item struct {
	ID        common.ID   `bson:"id"                   json:"id"`
	Version   vcs.Version `bson:"version,omitempty"    json:"version,omitempty"`
	CreatedAt time.Time   `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time  `bson:"updated_at,omitempty" json:"updated_at,omitempty"`

	Title   string       `bson:"title"             json:"title"`
	Brief   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
	Content content.Item `bson:"content,omitempty" json:"content,omitempty"`
	Links   links.Links  `bson:"links,omitempty"   json:"links,omitempty"`

	RView  common.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
	ROwner common.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
}

type Operator interface {
	Create(common.ID, Item) (dataspace.ID, error)

	Read(common.ID, dataspace.ID) (*Item, error)

	ReadList(common.ID, content.ListOptions) ([]Item, *uint64, error)

	Update(common.ID, Item) error

	UpdateLinks(common.ID, dataspace.ID, links.Links) error

	Delete(common.ID, dataspace.ID) error
}
