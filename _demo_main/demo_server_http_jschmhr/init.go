package demo_server_http_jsschmhr

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/basis/starter"

	"github.com/pavlo67/punctum/server_http"
)

func Starter() starter.Operator {
	return &datacompStarter{}
}

var l *zap.SugaredLogger

var endpoints map[string]config.Endpoint

type datacompStarter struct {
}

func (dcs *datacompStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *datacompStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	index, errs := config.ComponentIndex(params.StringKeyDefault("index_path", filelib.CurrentPath()), nil)

	endpoints = index.Endpoints

	return errs.Err()
}

func (dcs *datacompStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (dcs *datacompStarter) Setup() error {
	return nil
}

func (dcs *datacompStarter) Init(joiner program.Joiner) error {

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.New("no serverhttp_jschmhr.Operator interface found for confidenter.comp component")
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
	//confidenterOp, ok := joiner.Interface(confidenter_serverhttp_jschmhr.InterfaceKey).(componenthtml.Operator)
	//if ok {
	//	opsMap["confidenter"] = confidenterOp
	//} else {
	//	errs = append(errs, errors.Errorf("no componenthtml.Operator with key %s found for datacompStarter.init()", confidenter_serverhttp_jschmhr.InterfaceKey))
	//}

	// srvOp.HandleTemplator(Templator(opsMap, joiner))

	return errs.Err()
}
