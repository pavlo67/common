package scheduler

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "scheduler"

func Starter() starter.Operator {
	return &schedulerStarter{}
}

var l logger.Operator
var _ starter.Operator = &schedulerStarter{}

type schedulerStarter struct {
	interfaceKey joiner.InterfaceKey
	//config       server.Config
}

func (ss *schedulerStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *schedulerStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	return nil, nil
}

func (ss *schedulerStarter) Setup() error {
	return nil
}

func (ss *schedulerStarter) Run(_ joiner.Operator) error {
	return nil
}
