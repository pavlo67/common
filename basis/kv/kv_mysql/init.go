package kvmysql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/.off/kv"
	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceName = "kvmysql"

func Starter() starter.Operator {
	return &kvmysqlStarter{}
}

const tableDefault = "kv"

type kvmysqlStarter struct {
	interfaceKey joiner.InterfaceKey
	mysqlConfig  config.ServerAccess
	index        config.ServerComponentsIndex
	table        string
	addCRUD      bool
}

func (kvs *kvmysqlStarter) Name() string {
	return InterfaceName
}

func (kvs *kvmysqlStarter) Prepare(conf *config.Config, params basis.Info) error {
	var errs basis.Errors
	kvs.mysqlConfig, errs = conf.MySQL(params.StringDefault("database", "processor"), errs)

	indexPath := params.StringDefault("index_path", filelib.CurrentPath())

	kvs.index, errs = config.ComponentIndex(indexPath, errs)
	if len(errs) > 0 {
		return errs.Err()
	}

	kvs.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(kv.InterfaceKey)))

	kvs.table = params.StringDefault("table", "")
	if kvs.table == "" {
		kvs.table = tableDefault
	}

	return nil
}

func (kvs *kvmysqlStarter) Check() (info []joiner.Info, err error) {
	return mysqllib.CheckMySQLTables(
		kvs.mysqlConfig,
		kvs.index.MySQL,
		[]config.Table{{Key: "table", Name: kvs.table}},
	)

}

func (kvs *kvmysqlStarter) Setup() error {
	return mysqllib.SetupMySQLTables(
		kvs.mysqlConfig,
		kvs.index.MySQL,
		[]config.Table{{Key: "table", Name: kvs.table}},
	)
}

func (kvs *kvmysqlStarter) Init(joiner joiner.Operator) error {
	kvOp, err := New(
		kvs.mysqlConfig,
		kvs.table,
	)
	if err != nil {
		return err
	}

	err = joiner.JoinInterface(kvOp, kvs.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join kvmysql as "+string(kvs.interfaceKey))
	}

	return nil
}
