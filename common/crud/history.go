package crud

import (
	"time"

	"github.com/pavlo67/workshop/common"
)

const CreatedAction = "created_at"
const UpdatedAction = "updated_at"

type Action struct {
	Key    string    `bson:",omitempty" json:",omitempty"`
	Actor  common.ID `bson:",omitempty" json:",omitempty"`
	DoneAt time.Time `bson:",omitempty" json:",omitempty"`
}

// DEPRECATED
type History struct {
	Actions []Action `bson:",omitempty" json:",omitempty"`
}
