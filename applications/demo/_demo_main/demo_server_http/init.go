package demo_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/config"
	"github.com/pavlo67/constructor/components/common/joiner"
	"github.com/pavlo67/constructor/components/common/logger"
	"github.com/pavlo67/constructor/components/common/starter"
	"github.com/pavlo67/constructor/components/server/server_http"
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
