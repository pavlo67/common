package flow_routes

import (
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"
)

var Endpoints []server_http.Endpoint

var Prefix = "/flow/"
var PathBase = filelib.CurrentPath()
