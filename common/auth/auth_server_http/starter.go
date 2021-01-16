package auth_server_http

import (
	"github.com/pavlo67/workshop/common/data"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/auth/auth_jwt"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_http"

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
	interfaceKey        joiner.InterfaceKey
	authorizeHandlerKey joiner.InterfaceKey
	setCredsHandlerKey  joiner.InterfaceKey
	getCredsHandlerKey  joiner.InterfaceKey
}

func (ah *authHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ah *authHTTPStarter) Init(cfg *config.Config, lCommon logger.Operator, options data.Map) ([]data.Map, error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon
	ah.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	ah.authorizeHandlerKey = joiner.InterfaceKey(options.StringDefault("authorize_handler_key", string(auth.AuthorizeHandlerKey)))
	ah.setCredsHandlerKey = joiner.InterfaceKey(options.StringDefault("set_creds_handler_key", string(auth.SetCredsHandlerKey)))
	ah.getCredsHandlerKey = joiner.InterfaceKey(options.StringDefault("get_creds_handler_key", string(auth.GetCredsHandlerKey)))

	return nil, nil
}

func (ah *authHTTPStarter) Setup() error {
	return nil
}

func (ah *authHTTPStarter) Run(joinerOp joiner.Operator) error {
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
		return errors.New("no auth_jwt.ActorKey")
	}

	authorizeEndpoint, setCredsEndpoint, getCredsEndpoint, err := New()
	if err != nil {
		return errors.Wrap(err, "can'ah init auth.ActorKey")
	}

	err = joinerOp.Join(authorizeEndpoint, ah.authorizeHandlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join authorizeEndpoint as server_http.Endpoint with key '%s'", ah.authorizeHandlerKey)
	}

	err = joinerOp.Join(setCredsEndpoint, ah.setCredsHandlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join setCredsEndpoint as server_http.Endpoint with key '%s'", ah.setCredsHandlerKey)
	}

	err = joinerOp.Join(getCredsEndpoint, ah.getCredsHandlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join getCredsEndpoint as server_http.Endpoint with key '%s'", ah.getCredsHandlerKey)
	}

	return nil
}
