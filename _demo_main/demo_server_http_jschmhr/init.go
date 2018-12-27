package demo_server_http_jsschmhr

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/joiner"
	"github.com/pavlo67/punctum/basis/libs/filelib"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/server_http"
)

func Starter() starter.Operator {
	return &demo_server_http_jsschmhrStarter{}
}

var l logger.Operator
var endpoints map[string]config.Endpoint

type demo_server_http_jsschmhrStarter struct {
}

func (dcs *demo_server_http_jsschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *demo_server_http_jsschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	index, errs := config.ComponentIndex(params.StringKeyDefault("index_path", filelib.CurrentPath()), nil)

	endpoints = index.Endpoints

	return errs.Err()
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
		return errors.New("no server_http_jschmhr.Operator interface found for demo_server_http_jschmhr component")
	}

	errs := server_http.InitEndpoints(
		srvOp,
		endpoints,
		htmlHandlers,
		nil,
		nil,
		nil,
	)

	//opsMap := map[string]componenthtml.Operator{}
	//
	//confidenterOp, ok := joinerOp.Component(confidenter_serverhttp_jschmhr.InterfaceKey).(componenthtml.Operator)
	//if ok {
	//	opsMap["confidenter"] = confidenterOp
	//} else {
	//	errs = append(errs, errors.Errorf("no componenthtml.Operator with key %s found for datacompStarter.init()", confidenter_serverhttp_jschmhr.InterfaceKey))
	//}

	// srvOp.HandleTemplator(Templator(opsMap, joinerOp))

	return errs.Err()
}
