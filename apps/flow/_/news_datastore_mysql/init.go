package news_datastore_mysql

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/punctum/flow/datastore"
	"github.com/pavlo67/punctum/flow/datastore/datastore_mysql"
	"github.com/pavlo67/punctum/processor.old/news"
	"github.com/pavlo67/punctum/starter/logger"
)

const DatabaseDefault = "processor"
const TableDefault = "flow"

func Starter(addCRUD bool) starter.Operator {
	return &news_datastore_mysqlStarter{
		addCRUD: addCRUD,
	}
}

var l *zap.SugaredLogger

type news_datastore_mysqlStarter struct {
	interfaceKey        joiner.InterfaceKey
	interfaceKeyCRUD    joiner.InterfaceKey
	dsmysqlStarter      starter.Operator
	dsmysqlInterfaceKey joiner.InterfaceKey
	database            string
	table               string

	addCRUD bool
}

func (fs *news_datastore_mysqlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

const onPrepare = "on news_datastore_mysqlStarter.Init()"

func (fs *news_datastore_mysqlStarter) Prepare(conf *config.Config, params basis.Info) error {
	l = logger.zapGet()

	fs.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(news.InterfaceKey)))

	fs.dsmysqlStarter = datastore_mysql.Starter(false, &news.Content{})
	if fs.dsmysqlStarter == nil {
		return errors.Wrap(basis.ErrNull, onPrepare+": no datamysqlStarter.Operator")
	}

	fs.dsmysqlInterfaceKey = fs.interfaceKey
	fs.table = params.StringDefault("table", TableDefault)
	fs.database = params.StringDefault("database", DatabaseDefault)

	if fs.addCRUD {
		fs.interfaceKeyCRUD = joiner.InterfaceKey(params.StringDefault("interface_key_crud", string(news.InterfaceKeyCRUD)))
	}

	daParams := basis.Info{
		"interface_key": string(fs.dsmysqlInterfaceKey),
		"database":      fs.database,
		"table":         fs.table,
	}
	return fs.dsmysqlStarter.Init(conf, daParams)
}

func (fs *news_datastore_mysqlStarter) Check() ([]joiner.Info, error) {
	return fs.dsmysqlStarter.Check()
}

func (fs *news_datastore_mysqlStarter) Setup() error {
	return fs.dsmysqlStarter.Setup()
}

func (fs *news_datastore_mysqlStarter) Init(joiner joiner.Operator) error {
	var joinerForDS joiner.Joiner

	joinerForDS = joiner.NewJoiner()
	err := fs.dsmysqlStarter.Run(joinerForDS)
	if err != nil {
		return err
	}

	dsOpRaw := joinerForDS.Interface(fs.dsmysqlInterfaceKey)
	dsOp, ok := dsOpRaw.(datastore.Operator)
	if !ok {
		return errors.Errorf("no datastore.Operator with key %s for news_datastore_mysqlStarter, got %#v", fs.dsmysqlInterfaceKey, dsOpRaw)
	}

	flowOp, err := New(dsOp)
	if err != nil {
		return err
	}

	err = joiner.JoinInterface(flowOp, fs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join flowmysql.Operator as %s", fs.interfaceKey)
	}

	if fs.addCRUD {
		crudOp := news.OperatorCRUD{
			Operator: flowOp,
		}
		err = joiner.JoinInterface(crudOp, fs.interfaceKeyCRUD)

		if err != nil {
			return errors.Wrap(err, "can't join sources.OperatorCRUD as "+string(fs.interfaceKeyCRUD))
		}
	}

	return nil
}
