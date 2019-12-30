package receiver_server_http

import (
	"fmt"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/receiver"
)

var l logger.Operator

const Name = "receiver_server_http"

var _ starter.Operator = &receiverServerHTTPStarter{}

type receiverServerHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &receiverServerHTTPStarter{}
}

func (rs *receiverServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rs *receiverServerHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {

	rs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(receiver.InterfaceKey)))

	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", Name)
	}

	return nil, nil
}

func (rs *receiverServerHTTPStarter) Setup() error {
	return nil
}

func (rs *receiverServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	return nil
}
