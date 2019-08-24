package crud

import (
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/selectors"
)

const InterfaceKey joiner.InterfaceKey = "content"

type Type string

type Brief struct {
	ID      common.ID   `bson:"_id,omitempty"     json:"id,omitempty"`
	Type    Type        `bson:"type"                 json:"type"`
	Title   string      `bson:"title"             json:"title"`
	Summary string      `bson:"summary,omitempty" json:"summary,omitempty"`
	Info    common.Info `bson:"info,omitempty"    json:"info,omitempty"`
}

type Item struct {
	Brief   `            bson:",inline"           json:",inline"`
	Details interface{} `bson:"details,omitempty" json:"details,omitempty"`
}

type Description struct {
	Exemplar interface{} `json:"exemplar,omitempty"`
	Length   *int64      `json:"length,omitempty"`
}

type Part struct {
	Key   []string
	Count uint64
}

type Operator interface {
	Descript() (*Description, error)

	Save(content Item, options *SaveOptions) (id common.ID, err error)

	Read(id common.ID, options *GetOptions) (*Item, error)

	List(selector *selectors.Term, options *GetOptions) ([]Brief, error)

	Remove(id common.ID, options *RemoveOptions) error
}

type Cleaner func() error

type SaveOptions struct {
	AuthID    common.ID
	Replace   bool
	ReturnIDs bool
}

type GetOptions struct {
	AuthID  common.ID
	GroupBy []string
	OrderBy []string
}

type RemoveOptions struct {
	AuthID common.ID
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
//	RView  common.ID `bson:"r_view,omitempty"  json:"r_view,omitempty"`
//	ROwner common.ID `bson:"r_owner,omitempty" json:"r_owner,omitempty"`
