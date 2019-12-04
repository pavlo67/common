package workspace_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/workspace"
)

var workspaceOp workspace.Operator
var l logger.Operator

const Name = "workspace_starter"

var _ starter.Operator = &workspaceServerHTTPStarter{}

type workspaceServerHTTPStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &workspaceServerHTTPStarter{}
}

func (ss *workspaceServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName + "/" + Name
}

func (ss *workspaceServerHTTPStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *workspaceServerHTTPStarter) Setup() error {
	return nil
}

func (ss *workspaceServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	workspaceOp, ok = joinerOp.Interface(workspace.InterfaceKey).(workspace.Operator)
	if !ok {
		return errors.Errorf("no workspace.Operator with key %s", workspace.InterfaceKey)
	}

	return nil
}
