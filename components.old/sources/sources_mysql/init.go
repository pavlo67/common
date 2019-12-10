package sources_mysql

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter/logger"
)

const DatabaseDefault = "processor"
const TableDefault = "source"

func Starter(addCRUD bool) starter.Operator {
	return &sources_mysqlStarter{
		addCRUD: addCRUD,
	}
}

var l *zap.SugaredLogger

type sources_mysqlStarter struct {
	interfaceKey     joiner.InterfaceKey
	interfaceKeyCRUD joiner.InterfaceKey

	mysqlConfig config.ServerAccess
	index       config.ServerComponentsIndex
	table       string

	addCRUD bool
}

func (ss *sources_mysqlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *sources_mysqlStarter) Prepare(cfgCommon, cfg *config.Config, params basis.Info) error {

	l = logger.zapGet()

	ss.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(sources.InterfaceKey)))

	if ss.addCRUD {
		ss.interfaceKeyCRUD = joiner.InterfaceKey(params.StringDefault("interface_key_crud", string(sources.InterfaceKeyCRUD)))
	}

	ss.table = params.StringDefault("table", TableDefault)

	var errs basis.Errors
	ss.mysqlConfig, errs = conf.MySQL(params.StringDefault("database", DatabaseDefault), errs)

	ss.index, errs = config.ComponentIndex(params.StringDefault("index_path", filelib.CurrentPath()), errs)
	return errs.Err()
}

func (ss *sources_mysqlStarter) Check() ([]joiner.Info, error) {
	return mysqllib.CheckMySQLTables(
		ss.mysqlConfig,
		ss.index.MySQL,
		[]config.Table{{Key: "table", Name: ss.table}},
	)
}

func (ss *sources_mysqlStarter) Setup() error {
	return mysqllib.SetupMySQLTables(
		ss.mysqlConfig,
		ss.index.MySQL,
		[]config.Table{{Key: "table", Name: ss.table}},
	)
}

func (ss *sources_mysqlStarter) Init(joiner joiner.Operator) error {
	grOp, _ := joiner.Interface(groups.InterfaceKey).(groups.Operator)

	srcOp, err := New(
		grOp,
		ss.mysqlConfig,
		ss.table,
		nil, // is it ok?
	)
	if err != nil {
		return err
	}

	err = joiner.JoinInterface(srcOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join sourcesmysql as "+string(ss.interfaceKey))
	}

	if ss.addCRUD {
		crudOp := sources.OperatorCRUD{
			Operator: srcOp,
		}
		err = joiner.JoinInterface(crudOp, ss.interfaceKeyCRUD)

		if err != nil {
			return errors.Wrap(err, "can't join sources.OperatorCRUD as "+string(ss.interfaceKeyCRUD))
		}
	}

	return nil
}
