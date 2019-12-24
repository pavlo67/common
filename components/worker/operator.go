package worker

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/tasks"
)

type Operator interface {
	Name() string
	Run(task *tasks.Item, label string) (info common.Map, posterior []joiner.Link, err error)
}
