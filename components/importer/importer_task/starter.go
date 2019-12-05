package importer_task

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

func Starter() starter.Operator {
	return &importerTaskStarter{}
}

var l logger.Operator
var _ starter.Operator = &importerTaskStarter{}

type importerTaskStarter struct {
	//config       config.Access
	//table        string
	//interfaceKey joiner.InterfaceKey
}

func (ts *importerTaskStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *importerTaskStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	//var cfgSQLite config.Access
	//err := cfg.Value("sqlite", &cfgSQLite)
	//if err != nil {
	//	return nil, err
	//}
	//
	//ts.config = cfgSQLite
	//ts.table, _ = options.String("table")
	//ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))

	return nil, nil
}

func (ts *importerTaskStarter) Setup() error {
	return nil
}

func (ts *importerTaskStarter) Run(joinerOp joiner.Operator) error {
	return nil
}
