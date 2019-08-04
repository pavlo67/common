package structura

import (
	"github.com/pavlo67/constructor/apps/content"
	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/selectors"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "content"

type Item struct {
	content.Brief `            bson:",inline"           json:",inline"`
	Details       interface{} `bson:"details,omitempty" json:"details,omitempty"`
}

type Description struct {
	Exemplar interface{} `json:"exemplar,omitempty"`
	Length   *int64      `json:"length,omitempty"`
}

type Operator interface {
	Descript() (*Description, error)

	Save(content Item, options *SaveOptions) (id basis.ID, err error)

	Read(id basis.ID, options *GetOptions) (*Item, error)

	List(selector *selectors.Term, options *GetOptions) ([]content.Brief, error)

	Remove(id basis.ID, options *RemoveOptions) error
}

type Cleaner func() error

type SaveOptions struct {
	AuthID    auth.ID
	Replace   bool
	ReturnIDs bool
}

type GetOptions struct {
	AuthID auth.ID
	SortBy []string
}

type RemoveOptions struct {
	AuthID auth.ID
	Delete bool
}

//	ID        dataspace.ID `bson:"id"                   json:"id"`
//	Version   vcs.Version  `bson:"version,omitempty"    json:"version,omitempty"`
//	CreatedAt time.Time    `bson:"created_at,omitempty" json:"created_at"`
//	UpdatedAt *time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
//
//	Title   string       `bson:"title"             json:"title"`
//	Brief   string       `bson:"brief,omitempty"   json:"brief,omitempty"`
//	Author  string       `bson:"author,omitempty"  json:"author,omitempty"`
//	Item content.Item `bson:"content,omitempty" json:"content,omitempty"`
//	Links   links.Links  `bson:"links,omitempty"   json:"links,omitempty"`
//
//	RView  auth.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner auth.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
