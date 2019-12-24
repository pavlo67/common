package tasks_postgres

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
	return &tasksSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &tasksSQLiteStarter{}

type tasksSQLiteStarter struct {
	config       config.Access
	table        string
	interfaceKey joiner.InterfaceKey
}

func (ts *tasksSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *tasksSQLiteStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgSQLite config.Access
	err := cfg.Value("postgres", &cfgSQLite)
	if err != nil {
		return nil, err
	}

	ts.config = cfgSQLite
	ts.table, _ = options.String("table")
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(tasks.InterfaceKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (ts *tasksSQLiteStarter) Setup() error {
	return nil
}

func (ts *tasksSQLiteStarter) Run(joinerOp joiner.Operator) error {
	tasksOp, _, err := New(ts.config, ts.table, ts.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init tasks.Operator")
	}

	err = joinerOp.Join(tasksOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join &tasksSQLite as tasks.Operator with key '%s'", ts.interfaceKey)
	}

	return nil
}
