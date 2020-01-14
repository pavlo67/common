package runner_factory_goroutine

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/tasks"
)

func Starter() starter.Operator {
	return &runnerFactoryStarter{}
}

var l logger.Operator
var _ starter.Operator = &runnerFactoryStarter{}

type runnerFactoryStarter struct {
	interfaceKey joiner.InterfaceKey
	joinerOp     joiner.Operator
}

func (rfs *runnerFactoryStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rfs *runnerFactoryStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	rfs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(runner.FactoryInterfaceKey)))

	return nil, nil
}

func (rfs *runnerFactoryStarter) Setup() error {
	return nil
}

func (rfs *runnerFactoryStarter) Run(joinerOp joiner.Operator) error {
	tasksOp, ok := joinerOp.Interface(tasks.InterfaceKey).(tasks.Operator)
	if !ok {
		return errors.Errorf("no tasks.Operator with key %s", tasks.InterfaceKey)
	}

	runnerFactory, err := New(tasksOp, joinerOp)
	if err != nil {
		return errors.Wrap(err, "can't init runner.factory")
	}

	err = joinerOp.Join(runnerFactory, rfs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *runnerFactory as runner.Factory with key '%s'", rfs.interfaceKey)
	}

	return nil
}
