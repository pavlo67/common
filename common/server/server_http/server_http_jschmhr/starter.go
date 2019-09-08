package server_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/auth"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	interfaceKey joiner.InterfaceKey
	// interfaceKeyRouter joiner.InterfaceKey
	config config.Server

	staticPaths map[string]string
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Init(conf *config.Config, options common.Info) (info []common.Info, err error) {
	var errs common.Errors
	l = conf.Logger

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	// ss.interfaceKeyRouter = joiner.InterfaceKey(options.StringDefault("interface_key_router", string(controller.InterfaceKey)))

	ss.config = conf.Server
	if ss.config.Port <= 0 {
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %#v", ss.config))
	}

	// TODO: use more then one static path
	if staticPath, ok := options["static_path"]; ok {
		ss.staticPaths = map[string]string{"static": staticPath}
	}

	return nil, errs.Err()
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {

	authOpNil := auth.Operator(nil)
	authComps := joinerOp.ComponentsAllWithInterface(&authOpNil)

	var authOps []auth.Operator
	for _, authComp := range authComps {
		if authOp, ok := authComp.Interface.(auth.Operator); ok {
			authOps = append(authOps, authOp)
		}
	}

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, authOps)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.Operator")
	}

	for path, staticPath := range ss.staticPaths {
		srvOp.HandleFiles("/"+path+"/*filepath", staticPath, nil)
	}

	err = joinerOp.Join(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
