package tasks_pg

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/tasks"
)

func Starter() starter.Operator {
	return &tasksPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &tasksPgStarter{}

type tasksPgStarter struct {
	config       config.Access
	table        string
	interfaceKey joiner.InterfaceKey
}

func (ts *tasksPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *tasksPgStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgPg config.Access
	err := cfg.Value("postgres", &cfgPg)
	if err != nil {
		return nil, err
	}

	ts.config = cfgPg
	ts.table, _ = options.String("table")
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(tasks.InterfaceKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (ts *tasksPgStarter) Setup() error {
	return nil
}

func (ts *tasksPgStarter) Run(joinerOp joiner.Operator) error {
	tasksOp, _, err := New(ts.config, ts.table, ts.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init tasks.Operator")
	}

	err = joinerOp.Join(tasksOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join &tasksPg as tasks.Operator with key '%s'", ts.interfaceKey)
	}

	return nil
}
