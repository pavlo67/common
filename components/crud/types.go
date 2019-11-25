package crud

import (
	"time"

	"github.com/pavlo67/workshop/common"
)

type Status struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Counter map[string]uint64

type Index map[string][]common.ID
