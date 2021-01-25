package persons

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
)

type Item struct {
	auth.Identity `           json:",inline"`
	Data          common.Map `json:",omitempty"`
	CreatedAt     time.Time  `json:",omitempty"`
	UpdatedAt     *time.Time `json:",omitempty"`
}

type Operator interface {
	Add(identity auth.Identity, data common.Map, options *crud.Options) (auth.ID, error)
	Change(Item, *crud.Options) error
	Read(auth.ID, *crud.Options) (*Item, error)
	Remove(auth.ID, *crud.Options) error
	List(options *crud.Options) ([]Item, error)
}
