package manager_storage_sqlite

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/libraries/filelib"
	"github.com/pavlo67/workshop/libraries/sqllib"
	"github.com/pavlo67/workshop/libraries/sqllib/sqllib_sqlite"

	"github.com/pavlo67/workshop/components/data"
)

func Starter() starter.Operator {
	return &managerStorageSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &managerStorageSQLiteStarter{}

type managerStorageSQLiteStarter struct {
	config       config.ServerAccess
	index        config.ComponentsIndex
	interfaceKey joiner.InterfaceKey
}

func (fs *managerStorageSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fs *managerStorageSQLiteStarter) Init(conf *config.Config, options common.Options) ([]common.Options, error) {
	var errs common.Errors

	l = conf.Logger

	fs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))
	fs.config = conf.SQLite
	fs.index, errs = config.ComponentIndex(options.StringDefault("index_path", filelib.CurrentPath()), errs)

	sqlOp, err := sqllib_sqlite.New(fs.config)
	if err != nil {
		return nil, err
	}

	return sqllib.CheckTables(sqlOp, fs.index.SQLite)
}

func (fs *managerStorageSQLiteStarter) Setup() error {

	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Title: sm.table}},
	//)
}

func (fs *managerStorageSQLiteStarter) Run(joinerOp joiner.Operator) error {
	//sqlOp, err := sqllib_sqlite.New(fs.config)
	//if err != nil {
	//	return errors.Wrap(err, "can't init sqllib.Operator")
	//}
	//
	//db := sqlOp.DB()
	//
	//flowOp, err := New(db, 0)
	//if err != nil {
	//	return errors.Wrap(err, "can't init flow.Operator")
	//}
	//
	//err = joinerOp.Join(flowOp, fs.interfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join *flowSQLite as flow.Operator with key '%s'", fs.interfaceKey)
	//}

	return nil
}
