package crud

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/identity"
)

type ActionKey string

const ProducedAction ActionKey = "produced_from"
const CreatedAction ActionKey = "created"
const UpdatedAction ActionKey = "updated"

type Action struct {
	Identity   *identity.Item `bson:",omitempty" json:",omitempty"`
	Key        ActionKey      `bson:",omitempty" json:",omitempty"`
	DoneAt     time.Time      `bson:",omitempty" json:",omitempty"`
	RelatedIDs []common.ID    `bson:",omitempty" json:",omitempty"`
	Errors     common.Errors  `bson:",omitempty" json:",omitempty"`
}

type History []Action

func (h History) FirstByKey(key ActionKey) *Action {
	for _, action := range h {
		if action.Key == key {
			return &action
		}
	}

	return nil
}
