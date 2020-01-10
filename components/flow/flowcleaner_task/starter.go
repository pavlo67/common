package flowcleaner_task

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/dataimporter"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/flow"
)

func Starter() starter.Operator {
	return &importerTasksStarter{}
}

var l logger.Operator
var _ starter.Operator = &importerTasksStarter{}

type importerTasksStarter struct {
	datataggedKey joiner.InterfaceKey
	interfaceKey  joiner.InterfaceKey
}

// ------------------------------------------------

func (ts *importerTasksStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *importerTasksStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	ts.datataggedKey = joiner.InterfaceKey(options.StringDefault("datatagged_key", string(flow.InterfaceKey)))
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(dataimporter.TaskInterfaceKey)))

	return nil, nil
}

func (ts *importerTasksStarter) Setup() error {
	return nil
}

func (ts *importerTasksStarter) Run(joinerOp joiner.Operator) error {
	datataggedOp, ok := joinerOp.Interface(ts.datataggedKey).(datatagged.Operator)
	if !ok {
		return errors.Errorf("no datatagged.Actor with key %s", ts.datataggedKey)
	}

	impOp, err := NewLoader(datataggedOp)
	if err != nil {
		return errors.Wrap(err, "can't init flowimporter.Actor")
	}

	err = joinerOp.Join(impOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *cleanTask as actor.Actor with key '%s'", ts.interfaceKey)
	}

	return nil

}
