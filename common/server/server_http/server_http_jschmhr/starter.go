package server_http_jschmhr

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	config server.Config

	interfaceKey joiner.InterfaceKey
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	configKey := options.StringDefault("config_key", "server_http")
	if err := cfg.Value(configKey, &ss.config); err != nil {
		return nil, err
	}

	return nil, nil
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

	var requestOptions server_http.RequestOptions
	var secretENVs []string

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, requestOptions, secretENVs)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.UserKey")
	}

	err = joinerOp.Join(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.UserKey with key '%s'", ss.interfaceKey)
	}

	err = joinerOp.Join(ss.config.Port, server_http.PortInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join port with key '%s'", server_http.PortInterfaceKey)
	}

	err = joinerOp.Join(ss.config.NoHTTPS, server_http.NoHTTPSInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join no_https with key '%s'", server_http.NoHTTPSInterfaceKey)
	}

	return nil

}
