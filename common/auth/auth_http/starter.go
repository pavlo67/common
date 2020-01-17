package auth_http

import (
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/auth/auth_jwt"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

func Starter() starter.Operator {
	return &authHTTPStarter{}
}

var _ starter.Operator = &authHTTPStarter{}

var l logger.Operator
var authOps []auth.Operator
var authOpToInit auth.Operator
var authOpToSetToken auth.Operator

// var authOpToRegister auth.Operator

type authHTTPStarter struct {
	interfaceKey       joiner.InterfaceKey
	authHandlerKey     joiner.InterfaceKey
	authInitHandlerKey joiner.InterfaceKey
}

func (th *authHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (th *authHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	th.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))
	th.authHandlerKey = joiner.InterfaceKey(options.StringDefault("auth_handler_key", string(auth.AuthorizeHandlerKey)))
	th.authInitHandlerKey = joiner.InterfaceKey(options.StringDefault("auth_init_handler_key", string(auth.AuthInitHandlerKey)))

	return nil, nil
}

func (th *authHTTPStarter) Setup() error {
	return nil
}

func (th *authHTTPStarter) Run(joinerOp joiner.Operator) error {
	authOpNil := auth.Operator(nil)

	authComps := joinerOp.InterfacesAll(&authOpNil)
	for _, authComp := range authComps {
		if authOp, ok := authComp.Interface.(auth.Operator); ok {
			authOps = append(authOps, authOp)
			switch authComp.InterfaceKey {
			case auth_ecdsa.InterfaceKey:
				authOpToInit = authOp

			case auth_jwt.InterfaceKey:
				authOpToSetToken = authOp
				//case auth_users_sqlite.InterfaceKey:
				//	authOpToRegister = authOp
			}
		}
	}

	if authOpToSetToken == nil {
		return errors.New("no auth_jwt.Actor")
	}

	authEndpoint, authInitEndpoint, err := New()
	if err != nil {
		return errors.Wrap(err, "can'th init auth.Actor")
	}

	err = joinerOp.Join(authEndpoint, th.authHandlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join authEndpoint as server_http.Endpoint with key '%s'", th.authHandlerKey)
	}

	err = joinerOp.Join(authInitEndpoint, th.authInitHandlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join authInitSessionEndpoint as server_http.Endpoint with key '%s'", th.authInitHandlerKey)
	}

	return nil
}
