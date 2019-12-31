package gatherer_routes

import (
	"strconv"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/constructions/dataflow"
)

var endpoints = server_http.Endpoints{
	"flow": {Path: "/v1/export", Tags: []string{"flow"}, InterfaceKey: dataflow.ExportInterfaceKey},
}

func Init(srvOp server_http.Operator, port int) error {

	cfg := server_http.Config{
		Title:     "Pavlo's Gatherer REST API",
		Version:   "0.0.1",
		Prefix:    "/gatherer",
		Endpoints: endpoints,
	}

	return server_http.InitEndpointsWithSwaggerV2(
		cfg,
		":"+strconv.Itoa(port),
		srvOp,
		filelib.CurrentPath()+"api-docs/",
		"swagger.json",
		"api-docs",
		l,
	)

}
