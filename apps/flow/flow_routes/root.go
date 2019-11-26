package flow_routes1

import (
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/libraries/filelib"
)

var Endpoints []server_http.Endpoint

var Prefix = "/flow/"
var PathBase = filelib.CurrentPath()
