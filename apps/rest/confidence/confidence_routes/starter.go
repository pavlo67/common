package confidence_routes

import (
	"fmt"
	"log"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/server/server_http"
	"github.com/pavlo67/workshop/basis/starter"
)

const Name = "confidence_starter"

func Starter() starter.Operator {
	return &confidenceStarter{}
}

var L logger.Operator
var AuthOp auth.Operator
var Endpoints []server_http.Endpoint
var Prefix = "/confidence/"

var _ starter.Operator = &confidenceStarter{}

type confidenceStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *confidenceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *confidenceStarter) Init(cfg *config.Config, options common.Info) (info []common.Info, err error) {
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

	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	AuthOp, ok = joinerOp.Interface(auth.InterfaceKey).(auth.Operator)
	if !ok {
		log.Fatalf("no auth.Operator with key %s", auth.InterfaceKey)
	}

	for _, ep := range Endpoints {
		srvOp.HandleEndpoint(ep)
	}

	return nil
}
