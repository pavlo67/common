package runner

import (
	"time"

	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/tasks"
)

const FactoryInterfaceKey joiner.InterfaceKey = "runner_factory_goroutine"

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

func DataInterfaceKey(key crud.TypeKey) joiner.InterfaceKey {
	return joiner.InterfaceKey("actor for " + string(key))
}
