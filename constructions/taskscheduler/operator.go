package taskscheduler

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/actor"
)

const InterfaceKey joiner.InterfaceKey = "scheduler"

type Operator interface {
	Init(task actor.Operator) (common.ID, error)
	Run(taskID common.ID, interval time.Duration, startImmediately bool) error
	Stop(taskID common.ID) error
}
