package sources

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/instruments"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "founts"

type Item struct {
	ID    common.ID             `bson:"_id,omitempty"   json:"id,omitempty"`
	Title string                `bson:"title,omitempty" json:"title,omitempty"`
	URL   string                `bson:"url,omitempty"   json:"url,omitempty"`
	Tags  []string              `bson:"tags,omitempty"  json:"tags,omitempty"`
	Log   []instruments.LogItem `bson:"log,omitempty" json:"log,omitempty"`

	//Type      joiner.InterfaceKey
	//Params    basis.Options // for Create/Update methods for ex. tags list to set them on each imported item
	//ParamsRaw string     // for Read/ReadList methods

	SavedAt time.Time `bson:"saved_at"      json:"saved_at"`
}

// TODO!!!

type Operator interface {
	Save(url string, logItems ...instruments.LogItem) error
	Read(url string) (*Item, error)
	List(crud.GetOptions) ([]Item, *uint64, error)
	Delete(url string) error
}
