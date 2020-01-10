package runner

import (
	"time"

	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/common/identity"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/transport"

	"github.com/pavlo67/workshop/components/tasks"
)

const FactoryInterfaceKey joiner.InterfaceKey = "runner_factory"

type Estimate struct {
	Duration time.Duration
}

type Actor interface {
	Name() string
	Init(params common.Map) (estimate *Estimate, err error)
	Run() (response *tasks.Task, posterior []joiner.Link, err error)
}

type Operator interface {
	Run() (estimate *Estimate, err error)
	CheckTask() (item *tasks.Item, err error)
}

type Factory interface {
	NewRunner(item tasks.Item, transportOp transport.Operator, listener identity.Key) (Operator, error)
	NewRunnerFromTask(task tasks.Task, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener identity.Key) (Operator, common.ID, error)
}
