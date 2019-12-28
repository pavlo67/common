package crud

import (
	"time"

	"github.com/pavlo67/workshop/common"
)

type Action struct {
	Key    string    `bson:",omitempty" json:",omitempty"`
	Actor  common.ID `bson:",omitempty" json:",omitempty"`
	DoneAt time.Time `bson:",omitempty" json:",omitempty"`
}

type History struct {
	CreatedAt time.Time `bson:",omitempty" json:",omitempty"`
	Actions   []Action  `bson:",omitempty" json:",omitempty"`
}
