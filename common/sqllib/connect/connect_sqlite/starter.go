package connect_sqlite

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/sqllib/sqllib_sqlite"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/connect"
)

func Starter() starter.Operator {
	return &connectSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &connectSQLiteStarter{}

type connectSQLiteStarter struct {
	cfgSQLite config.Access

	interfaceKey joiner.InterfaceKey
}

func (css *connectSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (css *connectSQLiteStarter) Prepare(cfg *config.Config, options common.Map) error {
	dbKey := options.StringDefault("db_key", "sqlite")
	if err := cfg.Value(dbKey, &css.cfgSQLite); err != nil {
		return err
	}

	css.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(connect.InterfaceSQLiteKey)))

	return nil
}

const onRun = "on connectSQLiteStarter.Run()"

func (css *connectSQLiteStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	db, err := sqllib_sqlite.Connect(css.cfgSQLite)
	if err != nil || db == nil {
		return errors.CommonError(err, fmt.Sprintf(onRun+": got %#v", db))
	}

	if err = joinerOp.Join(db, css.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *sql.DB with key '%s'", css.interfaceKey))
	}

	return nil
}
