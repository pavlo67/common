package server_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/config"
	"github.com/pavlo67/constructor/components/common/joiner"
	"github.com/pavlo67/constructor/components/common/logger"
	"github.com/pavlo67/constructor/components/common/starter"
	"github.com/pavlo67/constructor/components/server/server_http"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	interfaceKey joiner.InterfaceKey
	// interfaceKeyRouter joiner.InterfaceKey
	config config.ServerTLS

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
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %d", ss.config.Port))
	}

	// TODO: use more then one static path
	if staticPath, ok := options.String("static_path"); ok {
		ss.staticPaths = map[string]string{"static": staticPath}
	}

	return nil, errs.Err()
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {
	identOpsMap := map[auth.CredsType][]auth.Operator{}

	// ???
	// authOpNil := auth.Operator(nil)
	// identOpsPtr := joinerOp.ComponentsAllWithInterface(&authOpNil)

	identOpsPtr := joinerOp.ComponentsAllWithInterface((*auth.Operator)(nil))
	for _, identOpIntf := range identOpsPtr {
		if identOp, ok := identOpIntf.Interface.(auth.Operator); ok {
			credsTypes, err := identOp.Accepts()
			if err != nil {
				l.Error(err)
			}
			for _, credsType := range credsTypes {
				identOpsMap[credsType] = append(identOpsMap[credsType], identOp)
			}
		}
	}

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, identOpsMap)
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
