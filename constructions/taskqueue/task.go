package taskqueue

import (
	"time"

	"github.com/pavlo67/workshop/common"
)

type Timing struct {
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type Result struct {
	Timing    `            bson:",inline"    json:",inline"`
	Success   bool        `bson:",omitempty" json:",omitempty"`
	Info      common.Map  `bson:",omitempty" json:",omitempty"`
	Posterior []common.ID `bson:",omitempty" json:",omitempty"`
}

type Status struct {
}
