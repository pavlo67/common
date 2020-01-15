package ws_routes

import (
	"strconv"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/storage"
)

var endpoints = server_http.Endpoints{
	"read": {Path: "/data/read", Tags: []string{"data"}, HandlerKey: storage.ReadInterfaceKey},
	"list": {Path: "/data/list", Tags: []string{"data"}, HandlerKey: storage.ListInterfaceKey},

	"save":   {Path: "/data/save", Tags: []string{"data"}, HandlerKey: storage.SaveInterfaceKey},
	"remove": {Path: "/data/remove", Tags: []string{"data"}, HandlerKey: storage.RemoveInterfaceKey},

	"tags":   {Path: "/data/tags", Tags: []string{"data"}, HandlerKey: storage.CountTagsInterfaceKey},
	"tagged": {Path: "/data/tagged", Tags: []string{"data"}, HandlerKey: storage.ListWithTagInterfaceKey},

	"flow_read": {Path: "/flow/read", Tags: []string{"flow"}, HandlerKey: flow.ReadInterfaceKey},
	"flow_list": {Path: "/flow/list", Tags: []string{"flow"}, HandlerKey: flow.ListInterfaceKey},
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
