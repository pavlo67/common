package ws_routes

import (
	"strconv"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/storage"
)

var endpoints = server_http.Endpoints{
	"read": {Path: "/v1/data/read", Tags: []string{"data"}, InterfaceKey: storage.ReadInterfaceKey},
	"list": {Path: "/v1/data/list", Tags: []string{"data"}, InterfaceKey: storage.ListInterfaceKey},

	"save":   {Path: "/v1/data/save", Tags: []string{"data"}, InterfaceKey: storage.SaveInterfaceKey},
	"remove": {Path: "/v1/data/remove", Tags: []string{"data"}, InterfaceKey: storage.RemoveInterfaceKey},

	"tags":   {Path: "/v1/data/tags", Tags: []string{"data"}, InterfaceKey: storage.CountTagsInterfaceKey},
	"tagged": {Path: "/v1/data/tagged", Tags: []string{"data"}, InterfaceKey: storage.ListWithTagInterfaceKey},

	"flow_read": {Path: "/v1/flow/read", Tags: []string{"flow"}, InterfaceKey: flow.ReadInterfaceKey},
	"flow_list": {Path: "/v1/flow/list", Tags: []string{"flow"}, InterfaceKey: flow.ListInterfaceKey},
}

func Init(srvOp server_http.Operator, port int) error {

	cfg := server_http.Config{
		Title:     "Pavlo's Workshop REST API",
		Version:   "0.0.1",
		Prefix:    "/storage",
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
