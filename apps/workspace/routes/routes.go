package routes

import (
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/workspace/workspace_server_http"
)

var Prefix = "/workspace/"

func InitEndpoints(srvOp server_http.Operator) {
	srvOp.HandleEndpoint(Prefix+"v1/save", workspace_server_http.SaveEndpoint)
	srvOp.HandleEndpoint(Prefix+"v1/read", workspace_server_http.ReadEndpoint)
	srvOp.HandleEndpoint(Prefix+"v1/list", workspace_server_http.ListEndpoint)
	srvOp.HandleEndpoint(Prefix+"v1/remove", workspace_server_http.RemoveEndpoint)
	srvOp.HandleFiles(Prefix+"api-docs/*filepath", server_http.StaticPath{LocalPath: filelib.CurrentPath() + "api-docs/", MIMEType: nil})
}
