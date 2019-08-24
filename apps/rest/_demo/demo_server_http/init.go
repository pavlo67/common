package demo_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/starter"
	"github.com/pavlo67/workshop/basis/server/server_http"
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

func (dcs *demoServerHTTPStarter) Init(conf *config.Config, options common.Info) ([]common.Info, error) {
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
