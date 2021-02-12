package persons

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/selectors"
)

type Item struct {
	auth.Identity `           json:",inline"`
	Data          common.Map `json:",omitempty"`
	CreatedAt     time.Time  `json:",omitempty"`
	UpdatedAt     *time.Time `json:",omitempty"`
}

type Operator interface {
	Add(identity auth.Identity, data common.Map, options *crud.Options) (auth.ID, error)
	Change(Item, *crud.Options) (*Item, error)
	Read(auth.ID, *crud.Options) (*Item, error)
	Remove(auth.ID, *crud.Options) error
	List(options *crud.Options) ([]Item, error)

	HasEmail(email string) (selectors.Term, error)
	HasNickname(nickname string) (selectors.Term, error)
}
