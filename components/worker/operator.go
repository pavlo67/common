package worker

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/tasks"
)

type Operator interface {
	Name() string
	Run(task *tasks.Task, label string) (posterior []joiner.Link, info common.Map, err error)
}
