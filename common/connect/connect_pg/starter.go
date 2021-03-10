package connect_pg

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/connect"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/sqllib/sqllib_pg"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &connectPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &connectPgStarter{}

type connectPgStarter struct {
	cfgPg config.Access

	interfaceKey joiner.InterfaceKey
}

func (css *connectPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (css *connectPgStarter) Prepare(cfg *config.Config, options common.Map) error {
	if err := cfg.Value(options.StringDefault("db_key", "pg"), &css.cfgPg); err != nil {
		return err
	}

	css.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(connect.InterfaceSQLiteKey)))

	return nil
}

const onRun = "on connectPgStarter.Run()"

func (css *connectPgStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	db, err := sqllib_pg.Connect(css.cfgPg)
	if err != nil || db == nil {
		return errors.CommonError(err, fmt.Sprintf(onRun+": got %#v", db))
	}

	if err = joinerOp.Join(db, css.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *sql.DB with key '%s'", css.interfaceKey))
	}

	return nil
}
