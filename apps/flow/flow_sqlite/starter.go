package flow_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/apps/flow"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/filelib"
	"github.com/pavlo67/constructor/basis/sqllib"
	"github.com/pavlo67/constructor/basis/sqllib/sqllib_sqlite"
	"github.com/pavlo67/constructor/starter"
	"github.com/pavlo67/constructor/starter/config"
	"github.com/pavlo67/constructor/starter/joiner"
	"github.com/pavlo67/constructor/starter/logger"
)

func Starter() starter.Operator {
	return &flowSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &flowSQLiteStarter{}

type flowSQLiteStarter struct {
	config       config.ServerAccess
	index        config.ComponentsIndex
	interfaceKey joiner.InterfaceKey
}

func (fs *flowSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fs *flowSQLiteStarter) Init(conf *config.Config, options basis.Info) ([]basis.Info, error) {
	var errs basis.Errors

	l = conf.Logger

	fs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(flow.InterfaceKey)))
	fs.config = conf.SQLite
	fs.index, errs = config.ComponentIndex(options.StringDefault("index_path", filelib.CurrentPath()), errs)

	sqlOp, err := sqllib_sqlite.New(fs.config)
	if err != nil {
		return nil, err
	}
	defer sqllib.Close(sqlOp)

	return sqllib.CheckTables(sqlOp, fs.index.SQLite)
}

func (fs *flowSQLiteStarter) Setup() error {

	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Name: sm.table}},
	//)
}

func (fs *flowSQLiteStarter) Run(joinerOp joiner.Operator) error {
	sqlOp, err := sqllib_sqlite.New(fs.config)
	if err != nil {
		return errors.Wrap(err, "can't init sqllib.Operator")
	}
	defer sqllib.Close(sqlOp)

	db, err := sqlOp.DB()
	if err != nil {
		return errors.Wrap(err, "can't get db from sqllib.Operator")
	}

	flowOp, err := New(db)
	if err != nil {
		return errors.Wrap(err, "can't init flow.Operator")
	}

	err = joinerOp.Join(flowOp, fs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *flowSQLite as flow.Operator with key '%s'", fs.interfaceKey)
	}

	return nil
}
