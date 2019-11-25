package workspace_starter

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

var L logger.Operator
var WorkspaceOp workspace.Operator

var _ starter.Operator = &workspaceStarter{}

type workspaceStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *workspaceStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Options) ([]common.Options, error) {
	var errs common.Errors

	L = lCommon
	if L == nil {
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

	WorkspaceOp, ok = joinerOp.Interface(workspace.InterfaceKey).(workspace.Operator)
	if !ok {
		return errors.Errorf("no workspace.Operator with key %s", workspace.InterfaceKey)
	}

	srvOp.HandleFiles("/workspace/api-docs/*filepath", filelib.CurrentPath()+"../docs/", nil)
	workspace_v1.Init()
	for _, ep := range workspace_routes.Endpoints {
		srvOp.HandleEndpoint(ep)
	}

	return nil
}
