package auth_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_jwt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_http"

func Starter() starter.Operator {
	return &authServerHTTPStarter{}
}

var _ starter.Operator = &authServerHTTPStarter{}

type authServerHTTPStarter struct {
	authKey    joiner.InterfaceKey
	authJWTKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var authOp auth.Operator

func (ashs *authServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ashs *authServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {

	ashs.authKey = joiner.InterfaceKey(options.StringDefault("auth_key", string(auth.InterfaceKey)))
	ashs.authJWTKey = joiner.InterfaceKey(options.StringDefault("auth_jwt_key", string(auth_jwt.InterfaceKey)))
	ashs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (ashs *authServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	// middleware -------------------------------------------------------

	authJWTOp, _ := joinerOp.Interface(ashs.authJWTKey).(auth.Operator)
	if authJWTOp == nil {
		return fmt.Errorf("no auth.Operator with key %s", ashs.authJWTKey)
	}

	middleware, err := OnRequestMiddleware(authJWTOp)
	if err != nil || middleware == nil {
		return fmt.Errorf("can't create server_http.OnRequestMiddleware(authJWTOp), got %#v, %s", middleware, err)
	}

	if err := joinerOp.Join(middleware, server_http.OnRequestMiddlewareInterfaceKey); err != nil {
		return errors.Wrapf(err, "can't join RequestOptions as server_http.onRequestMiddleware with key '%s'", server_http.OnRequestMiddlewareInterfaceKey)
	}

	// endpoints --------------------------------------------------------

	if authOp, _ = joinerOp.Interface(ashs.authKey).(auth.Operator); authOp == nil {
		return fmt.Errorf("no auth.Operator with key %s", ashs.authKey)
	}

	return server_http.JoinEndpoints(joinerOp, Endpoints)
}
