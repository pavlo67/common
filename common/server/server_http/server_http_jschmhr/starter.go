package server_http_jschmhr

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	config server.Config
	port   int

	interfaceKey     joiner.InterfaceKey
	portInterfaceKey joiner.InterfaceKey
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors
	l = lCommon

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	ss.portInterfaceKey = joiner.InterfaceKey(options.StringDefault("port_interface_key", string(server_http.PortInterfaceKey)))

	var cfgServerHTTP server.Config
	err := cfg.Value("server_http", &cfgServerHTTP)
	if err != nil {
		return nil, err
	}

	ss.config = cfgServerHTTP
	ss.port, _ = options.Int("port")

	return nil, errs.Err()
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {

	authOpNil := auth.Operator(nil)
	authComps := joinerOp.InterfacesAll(&authOpNil)

	var authOps []auth.Operator
	for _, authComp := range authComps {
		if authOp, ok := authComp.Interface.(auth.Operator); ok {
			authOps = append(authOps, authOp)
		}
	}

	srvOp, err := New(ss.port, ss.config.TLSCertFile, ss.config.TLSKeyFile, authOps)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.Operator")
	}

	err = joinerOp.Join(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	err = joinerOp.Join(ss.port, ss.portInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
	// return srvOp.Start()

}
