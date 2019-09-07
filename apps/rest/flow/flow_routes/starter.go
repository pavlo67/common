package flow_routes

import (
	"fmt"
	"log"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/server/server_http"
	"github.com/pavlo67/workshop/basis/starter"
	"github.com/pavlo67/workshop/components/data"
)

const Name = "flow_starter"

func Starter() starter.Operator {
	return &flowStarter{}
}

var L logger.Operator
var DataOp data.Operator
var Endpoints []server_http.Endpoint
var Prefix = "/flow/"

var _ starter.Operator = &flowStarter{}

type flowStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *flowStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *flowStarter) Init(cfg *config.Config, options common.Info) (info []common.Info, err error) {
	var errs common.Errors

	L = cfg.Logger
	if L == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *flowStarter) Setup() error {
	return nil
}

func (ss *flowStarter) Run(joinerOp joiner.Operator) error {

	//srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	//if !ok {
	//	log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	//}

	var ok bool
	DataOp, ok = joinerOp.Interface(data.InterfaceKey).(data.Operator)
	if !ok {
		log.Fatalf("no data.Operator with key %s", data.InterfaceKey)
	}

	return nil
}
