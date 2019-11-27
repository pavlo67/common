package kv_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/kv"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceName = "kv_sqlite"

func Starter() starter.Operator {
	return &kv_sqliteStarter{}
}

const tableDefault = "kv"

type kv_sqliteStarter struct {
	interfaceKey joiner.InterfaceKey
	config       config.ServerAccess
	table        string
}

func (kvs *kv_sqliteStarter) Name() string {
	return InterfaceName
}

func (kvs *kv_sqliteStarter) Init(conf *config.Config, params common.Map) (info []common.Map, err error) {

	kvs.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(kv.InterfaceKey)))

	kvs.table = tableDefault
	//kvs.table = params.StringDefault("table", "")
	//if kvs.table == "" {
	//	kvs.table = tableDefault
	//}

	// TODO: validate
	kvs.config = conf.SQLite

	//return mysqllib.CheckMySQLTables(
	//	kvs.mysqlConfig,
	//	kvs.index.MySQL,
	//	[]config.Table{{Key: "table", Name: kvs.table}},
	//)

	return nil, nil
}

func (kvs *kv_sqliteStarter) Setup() error {
	return nil
	//return mysqllib.SetupMySQLTables(
	//	kvs.mysqlConfig,
	//	kvs.index.MySQL,
	//	[]config.Table{{Key: "table", Name: kvs.table}},
	//)
}

func (kvs *kv_sqliteStarter) Run(joiner joiner.Operator) error {
	kvOp, err := New(kvs.config, kvs.table)
	if err != nil {
		return err
	}

	err = joiner.Join(kvOp, kvs.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join kvSQLite as "+string(kvs.interfaceKey))
	}

	return nil
}
