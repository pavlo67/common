package gatherer_routes

import (
	"github.com/pavlo67/workshop/common/server/server_http"

	"io/ioutil"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_server_http_handler"
	"github.com/pkg/errors"
)

var srvCfg = server_http.Config{
	Title:   "Pavlo's StorageIndex Gatherer REST API",
	Version: "0.0.1",
	Prefix:  "/gatherer",
	Endpoints: []server_http.EndpointConfig{
		{"flow", "/v1/export", []string{"flow"}, nil, flow_server_http_handler.ExportFlowEndpoint},
	},
}

func InitEndpoints(port string, srvOp server_http.Operator) error {
	swaggerPath := filelib.CurrentPath() + "api-docs/"
	swaggerFile := swaggerPath + "swagger.json"

	swagger, err := srvCfg.SwaggerV2(port)
	if err != nil {
		return errors.Errorf("on .SwaggerV2(%#v): %s", srvCfg, err)
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
