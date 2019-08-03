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
	return &demoServerHTTPStarter{}
}

var l logger.Operator

var _ starter.Operator = &demoServerHTTPStarter{}

type demoServerHTTPStarter struct{}

func (dcs *demoServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *demoServerHTTPStarter) Init(conf *config.Config, options basis.Info) ([]basis.Info, error) {
	l = logger.Get()

	return nil, nil
}

func (dcs *demoServerHTTPStarter) Setup() error {
	return nil
}

func (dcs *demoServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.New("no server_http_jschmhr.Operator interface found for demo_server_http component")
	}

	errs := server_http.InitEndpoints(
		srvOp,
		nil,
	)

	return errs.Err()
}
