package ws_routes

import (
	"github.com/pavlo67/workshop/common/server/server_http"

	"io/ioutil"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/components/flow/flow_server_http"
	"github.com/pavlo67/workshop/components/storage/storage_server_http"
	"github.com/pkg/errors"
)

var srvCfg = server_http.Config{
	Title:   "Pavlo's Workshop REST API",
	Version: "0.0.1",
	Prefix:  "/storage",
	Endpoints: []server_http.EndpointConfig{
		{"save", "/v1/data/save", []string{"data"}, nil, storage_server_http.SaveEndpoint},
		{"read", "/v1/data/read", []string{"data"}, nil, storage_server_http.ReadEndpoint},
		{"list", "/v1/data/list", []string{"data"}, nil, storage_server_http.ListEndpoint},
		{"remove", "/v1/data/remove", []string{"data"}, nil, storage_server_http.RemoveEndpoint},

		{"tags", "/v1/data/tags", []string{"data"}, nil, storage_server_http.CountTagsEndpoint},
		{"tagged", "/v1/data/tagged", []string{"data"}, nil, storage_server_http.ListWithTagEndpoint},

		{"flow", "/v1/flow/list", []string{"flow"}, nil, flow_server_http.FlowEndpoint},
		{"flow_read", "/v1/flow/read", []string{"flow"}, nil, flow_server_http.FlowReadEndpoint},
	},
}

func InitEndpoints(host string, srvOp server_http.Operator) error {
	swaggerPath := filelib.CurrentPath() + "api-docs/"
	swaggerFile := swaggerPath + "swagger.json"

	swagger, err := srvCfg.Swagger2(host)
	if err != nil {
		return errors.Errorf("on .Swagger2(%#v): %s", srvCfg, err)
	}

	err = ioutil.WriteFile(swaggerFile, swagger, 0644)
	if err != nil {
		return errors.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerFile, swagger, err)
	}
	l.Infof("%d bytes are written into %s", len(swagger), swaggerFile)

	err = server_http.InitEndpoints(srvCfg, srvOp, l)
	if err != nil {
		return err
	}
	return srvOp.HandleFiles("swagger", srvCfg.Prefix+"/api-docs/*filepath", server_http.StaticPath{LocalPath: swaggerPath, MIMEType: nil})
}
