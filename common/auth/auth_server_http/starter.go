package auth_server_http

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_http"

func Starter() starter.Operator {
	return &authHTTPStarter{}
}

var _ starter.Operator = &authHTTPStarter{}

var l logger.Operator
var authOps []auth.Operator
var authOpPersons auth.Operator

//var authOpToSetToken auth.Operator

// var authOpToRegister auth.Operator

type authHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
	// setTokenKey  joiner.InterfaceKey
}

func (ah *authHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ah *authHTTPStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon
	ah.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	// ah.setTokenKey = joiner.InterfaceKey(options.StringDefault("set_token_key", string(auth.InterfaceJWTKey)))

	return nil, nil
}

func (ah *authHTTPStarter) Setup() error {
	return nil
}

func (ah *authHTTPStarter) Run(joinerOp joiner.Operator) error {
	authOpNil := auth.Operator(nil)

	authComps := joinerOp.InterfacesAll(&authOpNil)
	for _, authComp := range authComps {
		if authOp, _ := authComp.Interface.(auth.Operator); authOp != nil {
			authOps = append(authOps, authOp)
			if authComp.InterfaceKey == auth.InterfaceKey {
				authOpPersons = authOp
			}
		}
	}

	//if authOpToSetToken, _ = joinerOp.Interface(ah.setTokenKey).(auth.Operator); authOpToSetToken == nil {
	//	return errors.New("no authOpToSetToken")
	//}
	if err := joinerOp.Join(&authEndpoint, auth.AuthHandlerKey); err != nil {
		return errors.Wrapf(err, "can't join authEndpoint as server_http.Endpoint with key '%s'", auth.AuthHandlerKey)
	}

	return nil
}
