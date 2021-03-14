package persons

import (
	"encoding/json"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/vcs"
)

const HasEmail selectors.Key = "has_email"
const HasNickname selectors.Key = "has_nickname"

type Operator interface {
	Save(Item, *auth.Identity) (auth.ID, error)
	Read(auth.ID, *auth.Identity) (*Item, error)
	Remove(auth.ID, *auth.Identity) error
	List(*selectors.Term, *auth.Identity) ([]Item, error)
	Stat(*selectors.Term, *auth.Identity) (db.StatMap, error)
}

type Item struct {
	auth.Identity `            json:",inline"    bson:",inline"`
	Data          common.Map  `json:",omitempty" bson:",omitempty"`
	History       vcs.History `json:",omitempty" bson:",omitempty"`
	CreatedAt     time.Time   `json:",omitempty" bson:",omitempty"`
	UpdatedAt     *time.Time  `json:",omitempty" bson:",omitempty"`

	// hidden values
	creds auth.Creds `json:",omitempty" bson:",omitempty"`
}

type Pack struct {
	structures.PackDescription
	Data []Item
}

func (item *Item) CompletePersonFromJSON(id auth.ID, rolesBytes, credsBytes, dataBytes, historyBytes []byte, email string) error {
	if item == nil {
		return errors.New("nil persons.Item to be completed")
	}

	item.Identity.ID = id

	if len(rolesBytes) > 0 {
		if err := json.Unmarshal(rolesBytes, &item.Identity.Roles); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Roles (%s)", rolesBytes)
		}
	}

	if len(credsBytes) > 0 {
		var creds auth.Creds
		if err := json.Unmarshal(credsBytes, &creds); err != nil {
			return errors.Wrapf(err, "can't unmarshal .creds (%s)", credsBytes)
		}
		creds[auth.CredsEmail] = email
		item.SetCreds(creds)
	}

	if len(dataBytes) > 0 {
		if err := json.Unmarshal(dataBytes, &item.Data); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Data (%s)", dataBytes)
		}
	}

	if len(historyBytes) > 0 {
		if err := json.Unmarshal(historyBytes, &item.History); err != nil {
			return errors.Wrapf(err, "can't unmarshal .History (%s)", historyBytes)
		}
	}

	return nil
}
