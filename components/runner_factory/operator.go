package runner_factory

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/tasks"
	"github.com/pavlo67/workshop/components/transport"
)

type Factory interface {
	TaskRunner(task tasks.Item, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener *transport.Listener) (runner.Operator, common.ID, error)
}
