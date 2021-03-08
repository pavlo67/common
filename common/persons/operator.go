package persons

import (
	"time"

	"github.com/pavlo67/data_exchange/components/vcs"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
)

type Operator interface {
	Add(identity auth.Identity, creds auth.Creds, data common.Map, options *crud.Options) (auth.ID, error)
	Change(Item, *crud.Options) (*Item, error)
	Read(auth.ID, *crud.Options) (*Item, error)
	Remove(auth.ID, *crud.Options) error
	List(options *crud.Options) ([]Item, error)

	HasEmail(email string) (selectors.Term, error)
	HasNickname(nickname string) (selectors.Term, error)
}

type Item struct {
	auth.Identity `            json:",inline"    bson:",inline"`
	creds         auth.Creds  `json:",omitempty" bson:",omitempty"`
	Data          common.Map  `json:",omitempty" bson:",omitempty"`
	History       vcs.History `json:",omitempty" bson:",omitempty"`
	CreatedAt     time.Time   `json:",omitempty" bson:",omitempty"`
	UpdatedAt     *time.Time  `json:",omitempty" bson:",omitempty"`
}
