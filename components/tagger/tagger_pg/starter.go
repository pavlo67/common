package tagger_pg

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/tagger"
)

func Starter() starter.Operator {
	return &taggerPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &taggerPgStarter{}

type taggerPgStarter struct {
	config              config.Access
	interfaceKey        joiner.InterfaceKey
	cleanerInterfaceKey joiner.InterfaceKey
}

func (ts *taggerPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *taggerPgStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	cfgPg := config.Access{}
	err := cfg.Value("pg", &cfgPg)
	if err != nil {
		return nil, err
	}

	ts.config = cfgPg
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(tagger.InterfaceKey)))
	ts.cleanerInterfaceKey = joiner.InterfaceKey(options.StringDefault("cleaner_interface_key", string(tagger.CleanerInterfaceKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (ts *taggerPgStarter) Setup() error {
	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Title: sm.table}},
	//)
}

func (ts *taggerPgStarter) Run(joinerOp joiner.Operator) error {
	taggerOp, taggerCleanerOp, err := New(ts.config, ts.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init tagger.Operator")
	}

	err = joinerOp.Join(taggerOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *taggerPg as tagger.Operator with key '%s'", ts.interfaceKey)
	}

	err = joinerOp.Join(taggerCleanerOp, ts.cleanerInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *taggerPg as tagger.Cleaner with key '%s'", ts.cleanerInterfaceKey)
	}

	return nil
}
