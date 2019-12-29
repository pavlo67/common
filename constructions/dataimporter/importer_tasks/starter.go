package importer_tasks

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_cleaner"
)

func Starter() starter.Operator {
	return &importerTasksStarter{}
}

var l logger.Operator
var _ starter.Operator = &importerTasksStarter{}

type importerTasksStarter struct {
	//config       config.Access
	//table        string
	//interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------

var fcOp crud.Cleaner

func (ts *importerTasksStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *importerTasksStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	//var cfgSQLite config.Access
	//err := cfg.Value("sqlite", &cfgSQLite)
	//if err != nil {
	//	return nil, err
	//}
	//
	//ts.config = cfgSQLite
	//ts.table, _ = options.Key("table")
	//ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))

	return nil, nil
}

func (ts *importerTasksStarter) Setup() error {
	return nil
}

func (ts *importerTasksStarter) Run(joinerOp joiner.Operator) error {

	fcOp, _ = joinerOp.Interface(flow_cleaner.InterfaceKey).(crud.Cleaner)
	if fcOp == nil {
		l.Fatalf("no flow_cleaner.Operator with key %s", flow_cleaner.InterfaceKey)
	}

	return nil
}
