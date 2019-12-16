package flow_cleaner_sqlite

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/flow/flow_cleaner"
	"github.com/pkg/errors"
)

func Starter() starter.Operator {
	return &flowCleanerSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &flowCleanerSQLiteStarter{}

type flowCleanerSQLiteStarter struct {
	config    config.Access
	table     string
	tableTags string

	interfaceKey joiner.InterfaceKey
}

func (fc *flowCleanerSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fc *flowCleanerSQLiteStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgSQLite config.Access
	err := cfg.Value("sqlite", &cfgSQLite)
	if err != nil {
		return nil, err
	}

	fc.config = cfgSQLite
	fc.table, _ = options.String("table")
	fc.tableTags, _ = options.String("table_tags")
	fc.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(flow_cleaner.InterfaceKey)))

	return nil, nil
}

func (fc *flowCleanerSQLiteStarter) Setup() error {
	return nil
}

func (fc *flowCleanerSQLiteStarter) Run(joinerOp joiner.Operator) error {

	fcOp, err := New(fc.config, fc.table, fc.tableTags, fc.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init flow_cleaner.Operator")
	}

	err = joinerOp.Join(fcOp, fc.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *flowCleanerSQLite as flow_cleaner.Operator with key '%s'", fc.interfaceKey)
	}

	return nil
}
