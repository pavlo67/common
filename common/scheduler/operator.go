package scheduler

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "scheduler"

type Task interface {
	Name() string
	Run(timeSheduled time.Time) error
}

type Operator interface {
	Init(task Task) (common.ID, error)
	Run(taskID common.ID, interval time.Duration, startImmediately bool) error
	Stop(taskID common.ID) error
}
