package links_mysql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/confidenter/groups"
	"github.com/pavlo67/workshop/notebook/links"
	"go.uber.org/zap"
)

var l *zap.SugaredLogger

const TableDefault = "link"

func Starter() starter.Operator {
	return &links_mysqlStarter{}
}

type links_mysqlStarter struct {
	interfaceKey        joiner.InterfaceKey
	cleanerInterfaceKey joiner.InterfaceKey
	mysqlConfig         config.ServerAccess
	conf                config.Config
	index               config.ComponentsIndex
	tables              []config.Table
}

var _ starter.Operator = &links_mysqlStarter{}

func (lms *links_mysqlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (lms *links_mysqlStarter) Prepare(conf *config.Config, params common.Info) error {

	l = logger.zapGet()

	var errs common.Errors
	lms.mysqlConfig, errs = conf.MySQL(params.StringDefault("database", ""), errs)

	indexPath := params.StringDefault("index_path", filelib.CurrentPath())

	lms.index, errs = config.ComponentIndex(indexPath, errs)
	if len(errs) > 0 {
		return errs.Err()
	}

	lms.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(links.InterfaceKey)))
	lms.cleanerInterfaceKey = joiner.InterfaceKey(params.StringDefault("cleaner_interface_key", string(links.CleanerInterfaceKey)))

	table := params.StringDefault("table", TableDefault)

	lms.tables = []config.Table{
		{Key: "table", Name: table},
	}

	return nil
}

func (lms *links_mysqlStarter) Check() (info []joiner.Info, err error) {
	return mysqllib.CheckMySQLTables(lms.mysqlConfig, lms.index.MySQL, lms.tables)
}

func (lms *links_mysqlStarter) Setup() error {
	return mysqllib.SetupMySQLTables(lms.mysqlConfig, lms.index.MySQL, lms.tables)
}

func (lms *links_mysqlStarter) Init(joiner joiner.Operator) error {
	grpOp, ok := joiner.Interface(groups.InterfaceKey).(groups.Operator)
	if !ok {
		return errors.New("no controller interface found for links_mysql.zapInit()")
	}

	linksOp, err := New(
		lms.mysqlConfig,
		lms.index.Params["table"],
		grpOp,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "can't init links_mysql ")
	}

	err = joiner.Join(linksOp, lms.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join links_mysql ")
	}

	//err = joiner.Join(dataOp.Clean, ds.cleanerInterfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
	//}

	return nil
}
