package persons

import (
	"encoding/json"
	"strings"

	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures"
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
	auth.Identity              ` json:",inline" bson:",inline"`
	structures.ItemDescription ` json:",inline" bson:",inline"`

	// hidden values
	creds auth.Creds `json:",omitempty" bson:",omitempty"`
}

type Pack struct {
	structures.PackDescription
	Items []Item
}

func (item *Item) UnfoldFromJSON(id auth.ID, rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, urnBytes, historyBytes []byte) error {
	if item == nil {
		return errors.New("nil persons.Item to be unfolded")
	}

	item.Identity.ID = id
	if len(rolesBytes) > 0 {
		if err := json.Unmarshal(rolesBytes, &item.Roles); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Roles (%s)", rolesBytes)
		}
	}

	var creds auth.Creds
	if len(credsBytes) > 0 {
		if err := json.Unmarshal(credsBytes, &creds); err != nil {
			return errors.Wrapf(err, "can't unmarshal .creds (%s)", credsBytes)
		}
	}
	if len(emailBytes) > 0 {
		if creds == nil {
			creds = auth.Creds{}
		}

		creds[auth.CredsEmail] = string(emailBytes)
	}
	item.SetCreds(creds)

	return IDUnfoldFromJSON(&item.ItemDescription, infoBytes, tagsBytes, urnBytes, historyBytes)
}

func (item *Item) FoldIntoJSON() (rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, historyBytes, urnBytes []byte, err error) {
	if item == nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.New("nil persons.Item to be folded")
	}

	rolesBytes = []byte{} // to satisfy NOT NULL constraint
	if len(item.Roles) > 0 {
		if rolesBytes, err = json.Marshal(item.Roles); err != nil {
			return nil, nil, nil, nil, nil, nil, nil, errors.Wrapf(err, "can't marshal .Roles (%#v)", item.Roles)
		}
	}

	creds := item.Creds()

	if email := strings.TrimSpace(creds[auth.CredsEmail]); email != "" {
		emailBytes = []byte(email)
	}

	delete(creds, auth.CredsEmail)

	credsBytes = []byte{} // to satisfy NOT NULL constraint
	if len(creds) > 0 {
		if credsBytes, err = json.Marshal(creds); err != nil {
			return nil, nil, nil, nil, nil, nil, nil, errors.Wrapf(err, "can't marshal creds (%#v)", creds)
		}
	}

	// TODO!!! append to item.History
	if infoBytes, tagsBytes, urnBytes, historyBytes, err = IDFoldIntoJSON(&item.ItemDescription); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, historyBytes, urnBytes, nil
}

// DEPRECATED
func IDUnfoldFromJSON(item *structures.ItemDescription, infoBytes, tagsBytes, urnBytes, historyBytes []byte) error {
	if item == nil {
		return errors.New("nil ItemDescription to be unfolded")
	}

	if len(infoBytes) > 0 {
		if err := json.Unmarshal(infoBytes, &item.Info); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Info (%s)", infoBytes)
		}
	}

	if len(tagsBytes) > 0 {
		if err := json.Unmarshal(tagsBytes, &item.Tags); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Tags (%s)", tagsBytes)
		}
	}

	item.URN = ns.URN(urnBytes)

	// TODO!!! append to item.History

	if len(historyBytes) > 0 {
		if err := json.Unmarshal(historyBytes, &item.History); err != nil {
			return errors.Wrapf(err, "can't unmarshal .History (%s)", historyBytes)
		}
	}

	return nil
}

// DEPRECATED
func IDFoldIntoJSON(item *structures.ItemDescription) (infoBytes, tagsBytes, urnBytes, historyBytes []byte, err error) {
	if item == nil {
		return nil, nil, nil, nil, errors.New("nil persons.Item to be folded")
	}

	infoBytes = []byte{} // to to satisfy NOT NULL constraint
	if len(item.Info) > 0 {
		if infoBytes, err = json.Marshal(item.Info); err != nil {
			return nil, nil, nil, nil, errors.Wrapf(err, "can't marshal .Info (%#v)", item.Info)
		}
	}

	tagsBytes = []byte{} // to to satisfy NOT NULL constraint
	if len(item.Tags) > 0 {
		if tagsBytes, err = json.Marshal(item.Tags); err != nil {
			return nil, nil, nil, nil, errors.Wrapf(err, "can't marshal .Tags (%#v)", item.Tags)
		}
	}

	if len(item.URN) > 0 {
		urnBytes = []byte(item.URN)
	}

	// TODO!!! append to item.History

	historyBytes = []byte{} // to to satisfy NOT NULL constraint
	if len(item.History) > 0 {
		historyBytes, err = json.Marshal(item.History)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrapf(err, "can't marshal .History(%#v)", item.History)
		}
	}

	return infoBytes, tagsBytes, urnBytes, historyBytes, nil
}
