package demo_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server/server_http"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &demo_server_http_jsschmhrStarter{}
}

var l logger.Operator

type demo_server_http_jsschmhrStarter struct{}

func (dcs *demo_server_http_jsschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *demo_server_http_jsschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
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

	srvOp.HandleTemplatorHTML(newTemplator(joinerOp))

	errs := server_http.InitEndpoints(
		srvOp,
		endpoints,
		htmlHandlers,
		nil,
		nil,
		nil,
	)

	return errs.Err()
}
