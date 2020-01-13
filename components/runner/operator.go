package runner

import (
	"time"

	"github.com/pavlo67/workshop/common/crud"

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
	Run() (info common.Map, posterior []joiner.Link, err error)
}

type Operator interface {
	Run() (estimate *Estimate, err error)
	CheckTask() (item *tasks.Item, err error)
}

type Factory interface {
	ItemRunner(item tasks.Item, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener *transport.Listener) (Operator, error)
	TaskRunner(task crud.Data, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener *transport.Listener) (Operator, common.ID, error)
}

func DataInterfaceKey(key crud.TypeKey) joiner.InterfaceKey {
	return joiner.InterfaceKey("actor for " + string(key))
}
