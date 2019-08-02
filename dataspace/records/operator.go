package records

import (
	"time"

	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/starter/joiner"

	"github.com/pavlo67/constructor/dataspace"
	"github.com/pavlo67/constructor/dataspace/content"
	"github.com/pavlo67/constructor/dataspace/links"
	"github.com/pavlo67/constructor/dataspace/vcs"
)

const InterfaceKey joiner.ComponentKey = "records"
const CleanerInterfaceKey joiner.ComponentKey = "records.cleaner"

// const GenusDefault = "note"
// const GenusFieldName = "genus"

//type Asked struct {
//	ID    string `bson:"_id,omitempty" json:"id"`
//	Genus string `bson:"genus"         json:"genus"`
//	Name  string `bson:"name"          json:"name"`
//}

type Item struct {
	ID        dataspace.ID `bson:"id"                   json:"id"`
	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
	CreatedAt time.Time    `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`

	Title   string       `bson:"title"             json:"title"`
	Brief   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
	Content content.Item `bson:"content,omitempty" json:"content,omitempty"`
	Links   links.Links  `bson:"links,omitempty"   json:"links,omitempty"`

	RView  auth.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
	ROwner auth.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
}

type Operator interface {
	Create(auth.ID, Item) (dataspace.ID, error)

	Read(auth.ID, dataspace.ID) (*Item, error)

	ReadList(auth.ID, content.ListOptions) ([]Item, *uint64, error)

	Update(auth.ID, Item) error

	UpdateLinks(auth.ID, dataspace.ID, links.Links) error

	Delete(auth.ID, dataspace.ID) error
}
