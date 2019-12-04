package flow_starter

import (
	"fmt"

	"github.com/pavlo67/workshop/apps/gatherer/flow_routes/v1"
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components.old/flow/flow_routes"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/libraries/filelib"
	"github.com/pkg/errors"
)

const Name = "flow_starter"

func Starter() starter.Operator {
	return &flowStarter{}
}

var L logger.Operator
var DataOp data.Operator

var _ starter.Operator = &flowStarter{}

type flowStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *flowStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *flowStarter) Init(cfg *config.Config, options common.Map) (info []common.Map, err error) {
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

	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	DataOp, ok = joinerOp.Interface(data.InterfaceKey).(data.Operator)
	if !ok {
		return errors.Errorf("no data.Operator with key %s", data.InterfaceKey)
	}

	srvOp.HandleFiles("/flow/api-docs/*filepath", filelib.CurrentPath()+"../docs/", nil)
	flow_v1.Init()
	for _, ep := range flow_routes.Endpoints {
		srvOp.HandleEndpoint(ep)
	}

	return nil
}
