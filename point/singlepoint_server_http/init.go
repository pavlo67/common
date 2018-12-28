package singlepoint_server_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/joiner"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/point"
	"github.com/pavlo67/punctum/server_http"
)

// the only one point description can be used for all points
// initiated with this package into some server instance (so the package is called "singlepoint")
// it can be corrected using closure for pointHandler()
var item point.Item

func Starter() starter.Operator {
	return &singlepoint_server_httpStarter{}
}

var l logger.Operator

type singlepoint_server_httpStarter struct{}

func (sps *singlepoint_server_httpStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sps *singlepoint_server_httpStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()
	item.Name = params.StringKeyDefault("name", "")

	return nil
}

func (sps *singlepoint_server_httpStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (sps *singlepoint_server_httpStarter) Setup() error {
	return nil
}

func (sps *singlepoint_server_httpStarter) Init(joinerOp joiner.Operator) error {
	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.New("no server_http.Operator interface found for singlepoint_server_http component")
	}

	errs := server_http.InitEndpoints(
		srvOp,
		endpoints,
		nil,
		restHandlers,
		nil,
		nil,
	)

	return errs.Err()
}
