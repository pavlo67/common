package workspace_v1

import (
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/workspace"
)

var l logger.Operator

var workspaceOp workspace.Operator

var endpoints []server_http.Endpoint

func Init(lCommon logger.Operator, workspaceOpCommon workspace.Operator) []server_http.Endpoint {
	l = lCommon

	workspaceOp = workspaceOpCommon

	return endpoints
}
