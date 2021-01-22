package crud

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
)

type ActionKey string

const ProducedAction ActionKey = "produced_from"
const SavedAction ActionKey = "saved"
const CreatedAction ActionKey = "created"
const UpdatedAction ActionKey = "updated"

type Action struct {
	ActorKey common.Key    `bson:",omitempty" json:",omitempty"`
	Key      ActionKey     `bson:",omitempty" json:",omitempty"`
	DoneAt   time.Time     `bson:",omitempty" json:",omitempty"`
	Related  *joiner.Link  `bson:",omitempty" json:",omitempty"`
	Errors   errors.Errors `bson:",omitempty" json:",omitempty"`
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

func (h History) CheckOn(hOld History) error {
	if len(hOld) < 1 {
		return nil
	}

	actionLast := hOld[len(hOld)-1]
	for _, actionNew := range h {
		if reflect.DeepEqual(actionLast, actionNew) {
			return nil
		}
	}

	return fmt.Errorf("history (%#v) is inappropriate to the old one (... %#v)", h, actionLast)
}
