package confidence_routes

import (
	"fmt"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/auth"
	"github.com/pavlo67/workshop/components/auth/auth_jwt"
	"github.com/pavlo67/workshop/components/auth/auth_users_sqlite"
	"github.com/pkg/errors"
)

const Name = "confidence_starter"

func Starter() starter.Operator {
	return &confidenceStarter{}
}

var Cfg *config.Config
var L logger.Operator
var AuthOps []auth.Operator

var AuthOpToSetToken auth.Operator
var AuthOpToRegister auth.Operator

var Endpoints []server_http.Endpoint

var BasePath = filelib.CurrentPath()
var Prefix = "/confidence/"

var _ starter.Operator = &confidenceStarter{}

type confidenceStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *confidenceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *confidenceStarter) Init(cfg *config.Config, options common.Info) (info []common.Info, err error) {
	Cfg = cfg

	var errs common.Errors

	L = cfg.Logger
	if L == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *confidenceStarter) Setup() error {
	return nil
}

func (ss *confidenceStarter) Run(joinerOp joiner.Operator) error {

	authOpNil := auth.Operator(nil)

	AuthOps = nil
	AuthOpToSetToken = nil

	authComps := joinerOp.ComponentsAllWithInterface(&authOpNil)
	for _, authComp := range authComps {
		if authOp, ok := authComp.Interface.(auth.Operator); ok {
			AuthOps = append(AuthOps, authOp)
			switch authComp.InterfaceKey {
			case auth_jwt.InterfaceKey:
				AuthOpToSetToken = authOp
			case auth_users_sqlite.InterfaceKey:
				AuthOpToRegister = authOp
			}
		}
	}

	if AuthOpToSetToken == nil {
		return errors.New("no auth_jwt.Operator")
	}

	//_, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	//if !ok {
	//	return errors.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	//}
	//for _, ep := range Endpoints {
	//	srvOp.HandleEndpoint(ep)
	//}

	return nil
}
