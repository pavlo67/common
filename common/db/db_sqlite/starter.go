package db_sqlite

import (
	"fmt"
	"os"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/sqllib/sqllib_sqlite"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "db_sqlite"

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
	if err := cfg.Value(options.StringDefault("db_key", "sqlite"), &css.cfgSQLite); err != nil {
		return err
	}

	css.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

const onRun = "on connectSQLiteStarter.Run()"

func (css *connectSQLiteStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if os.Getenv("SHOW_CONNECTS") != "" {
		l.Infof("CONNECTING TO SQLITE: %#v", css.cfgSQLite)
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
