package crud

import (
	"time"

	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/identity"
)

type ActionKey string

const ProducedAction ActionKey = "produced_from"
const CreatedAction ActionKey = "created"
const UpdatedAction ActionKey = "updated"

type Action struct {
	Actor   *identity.Key `bson:",omitempty" json:",omitempty"`
	Key     ActionKey     `bson:",omitempty" json:",omitempty"`
	DoneAt  time.Time     `bson:",omitempty" json:",omitempty"`
	Related *joiner.Link  `bson:",omitempty" json:",omitempty"`
	Errors  common.Errors `bson:",omitempty" json:",omitempty"`
}

type History []Action

func (h History) FirstByKey(key ActionKey, related *joiner.Link) int {
	for i, action := range h {
		if action.Key == key {
			if action.Related == nil {
				if related == nil {
					return i
				}
			} else {
				if related != nil && *action.Related == *related {
					return i
				}
			}
		}
	}

	return -1
}

func (h History) SaveAction(action Action) History {
	i := h.FirstByKey(action.Key, action.Related)
	if i >= 0 {
		h[i] = action
	} else {
		h = append(h, action)
	}
	return h
}
