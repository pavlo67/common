package content

import (
	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/selectors"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.ComponentKey = "content"

type Brief struct {
	ID    basis.ID   `json:"id,omitempty"`
	Title string     `json:"title,omitempty"`
	Brief string     `json:"brief,omitempty"`
	Info  basis.Info `json:"info,omitempty"`
}

type Item struct {
	Brief    `json:",inline"`
	Operator `json:",inline"`

	Details interface{} `json:"details,omitempty"`
}

type Description struct {
	Exemplar interface{} `json:"exemplar,omitempty"`
	Length   *int64      `json:"length,omitempty"`
}

type Operator interface {
	Descript() (*Description, error)

	Save(content Item, options *SaveOptions) (id basis.ID, err error)

	List(selector *selectors.Term, options *ListOptions) ([]Brief, error)

	Read(id basis.ID, options *ReadOptions) (*Item, error)

	Remove(id basis.ID, options *RemoveOptions) error
}

type Cleaner func() error

type SaveOptions struct {
	AuthID      auth.ID
	DontReplace bool
}

type ListOptions struct {
	AuthID auth.ID
	SortBy []string
}

type ReadOptions struct {
	AuthID auth.ID
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
