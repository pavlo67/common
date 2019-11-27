package workspace_routes_starter

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/workspace"

	"github.com/pavlo67/workshop/apps/workspace/workspace_routes"
	"github.com/pavlo67/workshop/apps/workspace/workspace_routes/v1"
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

func (ss *workspaceStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
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

	workspaceOp, ok := joinerOp.Interface(workspace.InterfaceKey).(workspace.Operator)
	if !ok {
		return errors.Errorf("no workspace.Operator with key %s", workspace.InterfaceKey)
	}

	var endpoints []server_http.Endpoint
	endpoints = append(endpoints, workspace_v1.Init(l, workspaceOp)...)
	for _, ep := range endpoints {
		srvOp.HandleEndpoint(ep)
	}
	srvOp.HandleFiles(workspace_routes.Prefix+"api-docs/*filepath", filelib.CurrentPath()+"../api-docs/", nil)

	return nil
}