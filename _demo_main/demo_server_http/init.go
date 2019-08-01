package demo_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/server/server_http"
	"github.com/pavlo67/constructor/starter"
	"github.com/pavlo67/constructor/starter/config"
	"github.com/pavlo67/constructor/starter/joiner"
	"github.com/pavlo67/constructor/starter/logger"
)

func Starter() starter.Operator {
	return &demo_server_http_jsschmhrStarter{}
}

var l logger.Operator

var _ starter.Operator = &demo_server_http_jsschmhrStarter{}

type demo_server_http_jsschmhrStarter struct{}

func (dcs *demo_server_http_jsschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *demo_server_http_jsschmhrStarter) Prepare(conf *config.Config, params, options basis.Info) error {
	l = logger.Get()

	return nil
}

func (dcs *demo_server_http_jsschmhrStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (dcs *demo_server_http_jsschmhrStarter) Setup() error {
	return nil
}

func (dcs *demo_server_http_jsschmhrStarter) Init(joinerOp joiner.Operator) error {
	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.New("no server_http_jschmhr.Operator interface found for demo_server_http component")
	}

	errs := server_http.InitEndpoints(
		srvOp,
		nil,
		nil,
	)

	return errs.Err()
}
