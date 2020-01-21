package flowimporter_task

import (
	"github.com/pavlo67/workshop/components/sources"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

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
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(flow.ImporterTaskInterfaceKey)))

	return nil, nil
}

func (ts *importerTasksStarter) Setup() error {
	return nil
}

func (ts *importerTasksStarter) Run(joinerOp joiner.Operator) error {
	datataggedOp, ok := joinerOp.Interface(ts.datataggedKey).(datatagged.Operator)
	if !ok {
		return errors.Errorf("no datatagged.ActorKey with key %s", ts.datataggedKey)
	}

	sourcesOp, ok := joinerOp.Interface(sources.InterfaceKey).(sources.Operator)
	if !ok {
		return errors.Errorf("no sources.ActorKey with key %s", sources.InterfaceKey)
	}

	impOp, err := New(datataggedOp, sourcesOp)
	if err != nil {
		return errors.Wrap(err, "can't init *loadTask")
	}

	err = joinerOp.Join(impOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *loadTask as actor.ActorKey with key '%s'", ts.interfaceKey)
	}

	return nil

}
