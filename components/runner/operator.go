package runner

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/tasks"
)

const FactoryInterfaceKey joiner.InterfaceKey = "runner_factory"

type Estimate struct {
}

type Actor interface {
	Name() string
	Init(params common.Map) (estimate *Estimate, err error)
	Run() (posterior []joiner.Link, info common.Map, err error)
}

type Operator interface {
	Run() (estimate *Estimate, err error)
	CheckResults() (status *tasks.Status, posterior []joiner.Link, info common.Map, err error)
}

type Factory interface {
	NewRunner(item tasks.Item) (Operator, error)
}
