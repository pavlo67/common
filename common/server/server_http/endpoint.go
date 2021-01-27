package server_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/logger"
)

const swaggerFile = "swagger.json"

func InitEndpointsWithSwaggerV2(cfg Config, host string, noHTTPS bool, srvOp Operator, swaggerPath, swaggerSubpath string, l logger.Operator) error {
	swaggerFilePath := swaggerPath + swaggerFile

	swagger, err := cfg.SwaggerV2(host, noHTTPS)
	if err != nil {
		l.Errorf("%#v", cfg)
		return fmt.Errorf("on .SwaggerV2(): %s", err) //
	}

	if err = ioutil.WriteFile(swaggerFilePath, swagger, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerFilePath, swagger, err)
	}
	l.Infof("%d bytes are written into %s", len(swagger), swaggerFilePath)

	if err = InitEndpoints(cfg, srvOp, l); err != nil {
		return err
	}
	return srvOp.HandleFiles("swagger", cfg.Prefix+"/"+swaggerSubpath+"/*filepath", StaticPath{LocalPath: swaggerPath})
}

func InitEndpoints(cfg Config, srvOp Operator, l logger.Operator) error {
	if srvOp == nil {
		return errata.New("on .InitEndpoints(): srvOp == nil")
	}

	for key, ep := range cfg.Endpoints {
		//if ep.Skip {
		//	continue
		//}
		if ep.Handler == nil {
			return fmt.Errorf("on InitEndpoints: no .Handler %#v", ep)
		}

		err := srvOp.HandleEndpoint(key, cfg.Prefix+ep.Path, *ep.Handler)
		if err != nil {
			return fmt.Errorf("on srvOp.HandleEndpoint(%s, %s, %#v): %s", key, ep.Path, ep.Handler, err)
		}
	}

	return nil
}

type Endpoint struct {
	Method      string          `json:",omitempty"`
	PathParams  []string        `json:",omitempty"`
	QueryParams []string        `json:",omitempty"`
	BodyParams  json.RawMessage `json:",omitempty"`

	WorkerHTTP

	// AllowedIDs []common.PbxID `json:"allowed_ids,omitempty"`
	// DataItem   interface{} `json:"data_item,omitempty"` // for Interface
	// SwaggerDescription string
}

var rePathParam = regexp.MustCompile(":[^/]+")

//func (ep Endpoint) PathWithParams(params ...string) string {
//	matches := rePathParam.FindAllStringSubmatchIndex(ep.Path, -1)
//
//	numMatches := len(matches)
//	if len(params) < numMatches {
//		numMatches = len(params)
//	}
//
//	path := ep.Path
//	for nm := numMatches - 1; nm >= 0; nm-- {
//		path = path[:matches[nm][0]] + url.PathEscape(strings.ReplaceTags(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
//	}
//
//	return path
//}

func (ep Endpoint) PathTemplate(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ep.PathParams) < 1 {
		return serverPath
	}

	var pathParams []string
	for _, pp := range ep.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp)
		} else {
			pathParams = append(pathParams, ":"+pp)
		}

	}

	return serverPath + "/" + strings.Join(pathParams, "/")
}

func (ep Endpoint) PathTemplateBraced(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ep.PathParams) < 1 {
		return serverPath
	}

	var pathParams []string
	for _, pp := range ep.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp[1:])
		} else {
			pathParams = append(pathParams, pp)
		}

	}

	return serverPath + "/{" + strings.Join(ep.PathParams, "}/{") + "}"
}
