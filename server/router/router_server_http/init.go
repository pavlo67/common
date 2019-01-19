package router_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server/router"
	"github.com/pavlo67/punctum/server/server_http"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &router_server_httpStarter{}
}

var l logger.Operator
var _ starter.Operator = &router_server_httpStarter{}

type router_server_httpStarter struct {
	interfaceKey joiner.InterfaceKey
	config       config.ServerTLS
}

func (rs *router_server_httpStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rs *router_server_httpStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	var errs basis.Errors

	rs.interfaceKey = joiner.InterfaceKey(params.StringKeyDefault("interface_key", string(router.InterfaceKey)))

	return errs.Err()
}

func (rs *router_server_httpStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (rs *router_server_httpStarter) Setup() error {
	return nil
}

func (rs *router_server_httpStarter) Init(joiner joiner.Operator) error {

	srvOp, _ := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)

	routerOp, err := New(srvOp)
	if err != nil {
		return errors.Wrap(err, "can't router_server_http.New()")
	}

	err = joiner.JoinInterface(routerOp, rs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join router_server_http.Operator as router.Operator with key '%s'", rs.interfaceKey)
	}

	return nil
}
