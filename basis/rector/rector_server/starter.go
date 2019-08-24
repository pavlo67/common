package rector_server

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/applications/rector"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/starter"
	"github.com/pavlo67/workshop/basis/server/server_http"
)

func Starter() starter.Operator {
	return &rector_serverStarter{}
}

var _ starter.Operator = &rector_serverStarter{}

var l logger.Operator
var srvOp server_http.Operator
var endpoints []server_http.Endpoint

type rector_serverStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (ss *rector_serverStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *rector_serverStarter) Init(conf *config.Config, options common.Info) (info []common.Info, err error) {
	l = conf.Logger

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(rector.InterfaceKey)))

	return nil, nil
}

func (ss *rector_serverStarter) Setup() error {
	return nil
}

func (ss *rector_serverStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	srvOp, ok = joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.New("no server_http.Operator for rector_server.Starter")
	}

	return nil
}
