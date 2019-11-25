package workspace_routes

import (
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"
)

var Endpoints []server_http.Endpoint

var Prefix = "/ws/"
var PathBase = filelib.CurrentPath()
