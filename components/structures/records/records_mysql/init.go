package records_mysql

import (
	"github.com/pavlo67/workshop/applications/records"
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/libs/sqllib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

var l logger.Operator
var _ starter.Operator = &records_mysqlStarter{}

func Starter(jointLinks bool) starter.Operator {
	return &records_mysqlStarter{
		jointLinks: jointLinks,
	}
}

const TableDefault = "records"

type records_mysqlStarter struct {
	interfaceKey        joiner.InterfaceKey
	cleanerInterfaceKey joiner.InterfaceKey
	mysqlConfig         config.ServerAccess
	conf                config.Config
	index               config.ComponentsIndex
	tables              []config.Table
	jointLinks          bool
}

func (nms *records_mysqlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nms *records_mysqlStarter) Prepare(conf *config.Config, options, runtimeOptions common.Map) error {

	l = logger.Get()

	var errs common.Errors
	nms.mysqlConfig, errs = conf.MySQL(options.StringDefault("database", ""), errs)

	indexPath := options.StringDefault("index_path", filelib.CurrentPath())

	nms.index, errs = config.ComponentIndex(indexPath, errs)
	if len(errs) > 0 {
		return errs.Err()
	}

	nms.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records.InterfaceKey)))
	nms.cleanerInterfaceKey = joiner.InterfaceKey(options.StringDefault("cleaner_interface_key", string(records.CleanerInterfaceKey)))

	table := options.StringDefault("table", TableDefault)

	nms.tables = []config.Table{
		{Key: "table", Name: table},
	}

	return nil
}

func (nms *records_mysqlStarter) Check() (info []common.Map, err error) {
	return sqllib.CheckSQLTables(nms.mysqlConfig, nms.index.MySQL, nms.tables)
}

func (nms *records_mysqlStarter) Setup() error {
	return sqllib.SetupTables(nms.mysqlConfig, nms.index.MySQL, nms.tables)
}

func (nms *records_mysqlStarter) Init(joiner joiner.Operator) error {
	//grpsOp, _ := joiner.Interface(groups.InterfaceKey).(groups.Operator)
	////if !ok {
	////	return errors.New("no groups.Operator found for notes_mysql")
	////}
	//
	//linksOp, _ := joiner.Interface(links.InterfaceKey).(links.Operator)
	//if !ok {
	//	return errors.New("no links.Operator found for notes_mysql")
	//}

	//generaOp, ok := joiner.Component(genera.InterfaceKey).(genera.Operator)
	//if !ok {
	//	return errors.New("no genera.Operator found for notes_mysql")
	//}

	//var err error
	//notesOp, err := New(
	//	nms.mysqlConfig,
	//	nms.index.Params["table"],
	//	nms.jointLinks,
	//	nil, // grpsOp,
	//	nil, // linksOp,
	//	nil,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "can't init notes_mysql")
	//}
	//
	//err = joiner.Join(notesOp, nms.interfaceKey)
	//if err != nil {
	//	return errors.Wrap(err, "can't join notes_mysql as notes.Operator interface")
	//}

	//err = joiner.Join(dataOp.Clean, ds.cleanerInterfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
	//}

	return nil
}

//// fixturer.Operator--------------------------------------------------------------------------------------------
//
//var _ fixturer.Operator = &records_mysqlStarter{}
//
//func (nms *records_mysqlStarter) NameGeneric() string {
//	return string(notes.InterfaceKey)
//}
//
//func (nms *records_mysqlStarter) Load(userIS common.ID, selector selectors.Selector, fixture fixturer.Fixture) error {
//	return nms.objectsOp.loadFixture(userIS, selector, fixture)
//}
