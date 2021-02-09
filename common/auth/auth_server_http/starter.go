package auth_server_http

import (
	"fmt"

	"github.com/pkg/errors"

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
	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator

func (ah *authServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ah *authServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {
	ah.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (ah *authServerHTTPStarter) Setup() error {
	return nil
}

func (ah *authServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if err := joinerOp.Join(&onRequestMiddleware{}, server_http.OnRequestInterfaceKey); err != nil {
		return errors.Wrapf(err, "can't join RequestOptions as server_http.onRequestMiddleware with key '%s'", server_http.OnRequestInterfaceKey)
	}

	return server_http.JoinEndpoints(joinerOp, Endpoints)
}
