package scheduler

import (
	"time"

	"github.com/pavlo67/workshop/components/runner"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "scheduler"

type Operator interface {
	Init(task runner.Actor) (common.ID, error)
	Run(taskID common.ID, interval time.Duration, startImmediately bool) error
	Stop(taskID common.ID) error
}
