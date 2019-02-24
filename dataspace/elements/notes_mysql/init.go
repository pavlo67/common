package notes_mysql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/notebook/links"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
	"go.uber.org/zap"
)

func Starter(jointLinks bool) starter.Operator {
	return &notes_mysqlStarter{
		jointLinks: jointLinks,
	}
}

var l *zap.SugaredLogger

const TableDefault = "note"

type notes_mysqlStarter struct {
	interfaceKey        joiner.InterfaceKey
	cleanerInterfaceKey joiner.InterfaceKey
	mysqlConfig         config.ServerAccess
	conf                config.PunctumConfig
	index               config.ServerComponentsIndex
	tables              []config.Table
	jointLinks          bool
}

var _ starter.Operator = &notes_mysqlStarter{}

func (nms *notes_mysqlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nms *notes_mysqlStarter) Prepare(conf *config.PunctumConfig, params basis.Options) error {

	l = logger.zapGet()

	var errs basis.Errors
	nms.mysqlConfig, errs = conf.MySQL(params.StringDefault("database", ""), errs)

	indexPath := params.StringDefault("index_path", filelib.CurrentPath())

	nms.index, errs = config.ComponentIndex(indexPath, errs)
	if len(errs) > 0 {
		return errs.Err()
	}

	nms.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(notes.InterfaceKey)))
	nms.cleanerInterfaceKey = joiner.InterfaceKey(params.StringDefault("cleaner_interface_key", string(notes.CleanerInterfaceKey)))

	table := params.StringDefault("table", TableDefault)

	nms.tables = []config.Table{
		{Key: "table", Name: table},
	}

	return nil
}

func (nms *notes_mysqlStarter) Check() (info []joiner.Info, err error) {
	return mysqllib.CheckMySQLTables(nms.mysqlConfig, nms.index.MySQL, nms.tables)
}

func (nms *notes_mysqlStarter) Setup() error {
	return mysqllib.SetupMySQLTables(nms.mysqlConfig, nms.index.MySQL, nms.tables)
}

func (nms *notes_mysqlStarter) Init(joiner joiner.Operator) error {
	grpsOp, _ := joiner.Interface(groups.InterfaceKey).(groups.Operator)
	//if !ok {
	//	return errors.New("no groups.Operator found for notes_mysql")
	//}

	linksOp, _ := joiner.Interface(links.InterfaceKey).(links.Operator)
	//if !ok {
	//	return errors.New("no links.Operator found for notes_mysql")
	//}

	//generaOp, ok := joiner.Component(genera.InterfaceKey).(genera.Operator)
	//if !ok {
	//	return errors.New("no genera.Operator found for notes_mysql")
	//}

	var err error
	notesOp, err := New(
		nms.mysqlConfig,
		nms.index.Params["table"],
		nms.jointLinks,
		grpsOp,
		linksOp,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "can't init notes_mysql")
	}

	err = joiner.JoinInterface(notesOp, nms.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join notes_mysql as notes.Operator interface")
	}

	//err = joiner.JoinInterface(dataOp.Clean, ds.cleanerInterfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
	//}

	return nil
}

//// fixturer.Operator--------------------------------------------------------------------------------------------
//
//var _ fixturer.Operator = &notes_mysqlStarter{}
//
//func (nms *notes_mysqlStarter) NameGeneric() string {
//	return string(notes.InterfaceKey)
//}
//
//func (nms *notes_mysqlStarter) Load(userIS auth.ID, selector selectors.Selector, fixture fixturer.Fixture) error {
//	return nms.objectsOp.loadFixture(userIS, selector, fixture)
//}
