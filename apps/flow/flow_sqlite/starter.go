package flow_sqlite

//import (
//	"github.com/pavlo67/constructor/basis"
//	"github.com/pavlo67/constructor/basis/filelib"
//	"github.com/pavlo67/constructor/basis/mysqllib"
//	"github.com/pavlo67/constructor/starter"
//	"github.com/pavlo67/constructor/starter/config"
//	"github.com/pavlo67/constructor/starter/logger"
//	"github.com/pkg/errors"
//	"go.uber.org/zap"
//)
//
//const DatabaseDefault = "processor"
//const TableDefault = "data"
//const TableTemporaryDefault = "data_temporary"
//
//func Starter(addCRUD, staged bool, contentTemplate interface{}) starter.Operator {
//	return &datamysqlStarter{
//		addCRUD:         addCRUD,
//		staged:          staged,
//		contentTemplate: contentTemplate,
//	}
//}
//
//var l *zap.SugaredLogger
//
//type datamysqlStarter struct {
//	interfaceKey        program.InterfaceKey
//	cleanerInterfaceKey program.InterfaceKey
//	mysqlConfig         config.ServerAccess
//	index               config.ServerComponentsIndex
//	tables              []config.Table
//	tableTemporary      string
//	addCRUD             bool
//	staged              bool
//	contentTemplate     interface{}
//}
//
//// TODO: implement addCRUD feature!!!
//
//func (ds *datamysqlStarter) Name() string {
//	return logger.GetCallInfo().PackageName
//}
//
//func (ds *datamysqlStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
//	l = logger.Get()
//
//	var errs basis.Errors
//	ds.mysqlConfig, errs = conf.MySQL(params.StringKeyDefault("database", DatabaseDefault), errs)
//
//	indexPath := params.StringKeyDefault("index_path", filelib.CurrentPath())
//
//	ds.index, errs = config.ComponentIndex(indexPath, errs)
//	if len(errs) > 0 {
//		return errs.Err()
//	}
//
//	ds.interfaceKey = program.InterfaceKey(params.StringKeyDefault("interface_key", string(datastore.InterfaceKey)))
//	ds.cleanerInterfaceKey = program.InterfaceKey(params.StringKeyDefault("cleaner_interface_key", string(datastore.CleanerInterfaceKey)))
//
//	table := params.StringKeyDefault("table", TableDefault)
//
//	ds.tables = []config.Table{
//		{Key: "table", Name: table},
//	}
//
//	if ds.staged {
//		ds.tableTemporary = params.StringKeyDefault("table_temporary", "")
//		if ds.tableTemporary == "" {
//			ds.tableTemporary = TableTemporaryDefault
//		}
//		ds.tables = append(ds.tables, config.Table{Key: "table_temporary", Name: ds.tableTemporary})
//	}
//
//	return nil
//}
//
//func (ds *datamysqlStarter) Check() (info []program.Info, err error) {
//	return mysqllib.CheckMySQLTables(ds.mysqlConfig, ds.index.MySQL, ds.tables)
//}
//
//func (ds *datamysqlStarter) Setup() error {
//	return mysqllib.SetupMySQLTables(ds.mysqlConfig, ds.index.MySQL, ds.tables)
//}
//
//func (ds *datamysqlStarter) Init(joiner program.Joiner) error {
//	dataOp, err := New(
//		ds.mysqlConfig,
//		ds.tables[0].Name,
//		ds.tableTemporary,
//		ds.contentTemplate,
//	)
//	if err != nil {
//		return err
//	}
//
//	err = joiner.JoinInterface(dataOp, ds.interfaceKey)
//	if err != nil {
//		return errors.Wrapf(err, "can't join datastoremysql.Operator as %s", ds.interfaceKey)
//	}
//
//	err = joiner.JoinInterface(dataOp.Clean, ds.cleanerInterfaceKey)
//	if err != nil {
//		return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
//	}
//
//	return nil
//}
