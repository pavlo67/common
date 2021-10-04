package server_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
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

func (ss *server_http_jschmhrStarter) Prepare(cfg *config.Config, options common.Map) error {
	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	configKey := options.StringDefault("config_key", "server_http")
	if err := cfg.Value(configKey, &ss.config); err != nil {
		return err
	}

	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	// TODO!!! customize it
	var secretENVs []string

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, secretENVs)
	if err != nil {
		return errors.Wrap(err, "on server_http_jschmhr.New()")
	}

	if err = joinerOp.Join(srvOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *serverHTTPJschmhr{} as server_http.Operator with key '%s'", ss.interfaceKey)
	}

	return nil

}
