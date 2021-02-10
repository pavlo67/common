package auth_http

import (
	"fmt"

	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &authHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &authHTTPStarter{}

type authHTTPStarter struct {
	serverConfig server_http.Config
	interfaceKey common.InterfaceKey
}

func (ahs *authHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ahs *authHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {

	var access config.Access
	if err := cfg.Value("auth_http", &access); err != nil {
		return err
	}

	// TODO!!! pass for each server separately
	prefix := options.StringDefault("prefix", "")

	var ok bool
	if ahs.serverConfig, ok = options["server_config"].(server_http.Config); !ok {
		return errors.New("no server config for authHTTPStarter")
	}

	ahs.serverConfig.CompleteDirectly(auth_server_http.Endpoints, access.Host, access.Port, prefix)

	ahs.interfaceKey = common.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil
}

//} else if endpointsPtr, ok := options["endpoints"].(*server_http.Endpoints); ok {
//	ihs.endpoints = *endpointsPtr

func (ahs *authHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	authOp, err := New(ahs.serverConfig)
	if err != nil {
		return errors.Wrap(err, "can't init *authHTTP{} as auth.Operator")
	}

	if err = joinerOp.Join(authOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *authHTTP{} as auth.Operator with key '%s'", ahs.interfaceKey)
	}

	return nil
}
