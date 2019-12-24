package crud

import (
	"time"

	"github.com/pavlo67/workshop/common"
)

type History struct {
	CreatedAt time.Time  `bson:",omitempty" json:",omitempty"`
	UpdatedAt *time.Time `bson:",omitempty" json:",omitempty"`

	// History string    `bson:",omitempty" json:",omitempty"`
}

type Counter map[string]uint64

type Index map[string][]common.ID
