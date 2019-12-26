package taskscheduler

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/actor"
)

const InterfaceKey joiner.InterfaceKey = "scheduler"

type Operator interface {
	Init(task actor.Operator) (common.Key, error)
	Run(taskID common.Key, interval time.Duration, startImmediately bool) error
	Stop(taskID common.Key) error
}
