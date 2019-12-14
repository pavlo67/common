package gatherer_routes

import (
	"fmt"

	"github.com/pkg/errors"

	"strconv"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
)

const Name = "workspace_starter"

func Starter() starter.Operator {
	return &workspaceStarter{}
}

var l logger.Operator

var _ starter.Operator = &workspaceStarter{}

type workspaceStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName + "/" + Name
}

func (ss *workspaceStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *workspaceStarter) Setup() error {
	return nil
}

func (ss *workspaceStarter) Run(joinerOp joiner.Operator) error {
	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	//wsOp, ok := joinerOp.Interface(storage.InterfaceKey).(storage.Operator)
	//if !ok {
	//	return errors.Errorf("no storage.Operator with key %s", storage.InterfaceKey)
	//}

	srvPort, ok := joinerOp.Interface(server_http.PortInterfaceKey).(int)
	if !ok {
		return errors.Errorf("no server_http.Port with key %s", server_http.PortInterfaceKey)
	}

	return InitEndpoints(":"+strconv.Itoa(srvPort), srvOp)
}
