package server_http_jschmhr

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
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
	onRequest, _ := joinerOp.Interface(server_http.OnRequestInterfaceKey).(server_http.OnRequest)
	if onRequest == nil {
		return fmt.Errorf("no server_http.OnRequest with key %s", server_http.OnRequestInterfaceKey)
	}

	// TODO!!! customize it
	var secretENVs []string

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, onRequest, secretENVs)
	if err != nil {
		return errata.Wrap(err, "on server_http_jschmhr.New()")
	}

	if err = joinerOp.Join(srvOp, ss.interfaceKey); err != nil {
		return errata.Wrapf(err, "can't join *serverHTTPJschmhr{} as server_http.Operator with key '%s'", ss.interfaceKey)
	}

	if err = joinerOp.Join(ss.config.TLSCertFile != "" && ss.config.TLSKeyFile != "", server_http.HTTPSInterfaceKey); err != nil {
		return errata.Wrapf(err, "can't join HTTPS info with key '%s'", server_http.HTTPSInterfaceKey)
	}

	if err = joinerOp.Join(ss.config.Port, server_http.PortInterfaceKey); err != nil {
		return errata.Wrapf(err, "can't join port with key '%s'", server_http.PortInterfaceKey)
	}

	return nil

}
