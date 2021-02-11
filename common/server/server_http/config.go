package server_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

type Endpoint struct {
	Method      string          `json:",omitempty"`
	PathParams  []string        `json:",omitempty"`
	QueryParams []string        `json:",omitempty"`
	BodyParams  json.RawMessage `json:",omitempty"`
	Produces    []string        `json:",omitempty"`

	WorkerHTTP
}

type EndpointSettled struct {
	Path string
	Tags []string

	Endpoint
	EndpointInternalKey joiner.InterfaceKey
}

type Config struct {
	Title            string
	Version          string
	Host             string
	Port             string
	Prefix           string
	EndpointsSettled map[joiner.InterfaceKey]EndpointSettled
}

type Swagger map[string]interface{}

// type SwaggerEndpoint struct {}

func (c Config) SwaggerV2(isHTTPS bool) ([]byte, error) {
	paths := map[string]common.Map{} // map[string]map[string]map[string]interface{}{}

	for key, ep := range c.EndpointsSettled {

		path := c.Prefix + ep.Endpoint.PathTemplateBraced(ep.Path)
		method := strings.ToLower(ep.Endpoint.Method)

		epDescr := common.Map{
			"operationId": key,
			"tags":        ep.Tags,
		}

		if len(ep.Produces) >= 1 {
			epDescr["produces"] = ep.Produces
		} else {
			epDescr["produces"] = []string{"application/json"}
		}

		var parameters []interface{} // []map[string]interface{}

		for _, pp := range ep.Endpoint.PathParams {
			if len(pp) > 0 && pp[0] == '*' {
				pp = pp[1:]
			}

			parameters = append(
				parameters,
				common.Map{
					"in":          "path",
					"required":    true,
					"name":        pp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}
		for _, qp := range ep.Endpoint.QueryParams {
			parameters = append(
				parameters,
				common.Map{
					"in":          "query",
					"required":    false, // TODO!!!
					"name":        qp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}

		if method == "post" {
			if len(ep.Endpoint.BodyParams) > 0 {
				parameters = append(parameters, ep.Endpoint.BodyParams)
			} else {
				parameters = append(parameters, common.Map{
					"in":       "body",
					"required": true,
					"name":     "body_item",
					"type":     "string",
				})
			}
		}

		if len(parameters) > 0 {
			epDescr["parameters"] = parameters
		}

		if epDescrPrev, ok := paths[path][method]; ok {
			return nil, fmt.Errorf("duplicate endpoint description (%s %s): \n%#v\nvs.\n%#v", method, path, epDescrPrev, epDescr)
		}
		if _, ok := paths[path]; ok { // pathPrev
			paths[path][method] = epDescr
		} else {
			paths[path] = common.Map{method: epDescr} // map[string]map[string]interface{}
		}
	}

	var schemes []string
	if isHTTPS {
		schemes = []string{"https", "http"}
	} else {
		schemes = []string{"http"}
	}

	swagger := Swagger{
		"swagger": "2.0",
		"info": map[string]string{
			"title":   c.Title,
			"version": c.Version,
		},
		// "basePath": c.Prefix,
		"schemes": schemes,
		"port":    c.Port,
		"paths":   paths,
	}

	return json.MarshalIndent(swagger, "", " ")
}

const swaggerFile = "swagger.json"

const onInitEndpointsWithSwaggerV2 = "on server_http.InitEndpointsWithSwaggerV2()"

func InitEndpointsWithSwaggerV2(srvOp Operator, cfg Config, isHTTPS bool, swaggerPath, swaggerSubpath string, l logger.Operator) error {
	if srvOp == nil {
		return errors.New(onInitEndpointsWithSwaggerV2 + ": srvOp == nil")
	}

	swaggerFilePath := swaggerPath + swaggerFile

	swagger, err := cfg.SwaggerV2(isHTTPS)
	if err != nil {
		return fmt.Errorf(onInitEndpointsWithSwaggerV2+": %s", err) //
	}

	if err = ioutil.WriteFile(swaggerFilePath, swagger, 0644); err != nil {
		return fmt.Errorf(onInitEndpointsWithSwaggerV2+": on ioutil.WriteFile(%s, %s, 0755): %s", swaggerFilePath, swagger, err)
	}
	l.Infof(onInitEndpointsWithSwaggerV2+": %d bytes are written into %s", len(swagger), swaggerFilePath)

	for key, ep := range cfg.EndpointsSettled {
		if err := srvOp.HandleEndpoint(key, cfg.Prefix+ep.Path, ep.Endpoint); err != nil {
			return fmt.Errorf(onInitEndpointsWithSwaggerV2+": handling endpoint(%s, %s, %#v) got %s", key, ep.Path, ep.Endpoint, err)
		}
	}

	return srvOp.HandleFiles("swagger", cfg.Prefix+"/"+swaggerSubpath+"/*filepath", StaticPath{LocalPath: swaggerPath})
}

const onInitPages = "on server_http.InitPages()"

func InitPages(srvOp Operator, pagesCfg Config, l logger.Operator) error {
	if srvOp == nil {
		return errors.New(onInitPages + ": srvOp == nil")
	}

	for key, ep := range pagesCfg.EndpointsSettled {
		if err := srvOp.HandleEndpoint(key, pagesCfg.Prefix+ep.Path, ep.Endpoint); err != nil {
			return fmt.Errorf(onInitPages+": handling %s, %s, %#v got %s", key, ep.Path, ep.Endpoint, err)
		}
	}

	return nil
	// return srvOp.HandleFiles("swagger", pagesCfg.Prefix+"/"+swaggerSubpath+"/*filepath", StaticPath{LocalPath: swaggerPath})
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

// joining endpoints -----------------------------------------------------

type Endpoints map[joiner.InterfaceKey]Endpoint

func JoinEndpoints(joinerOp joiner.Operator, eps Endpoints) error {
	for key, ep := range eps {
		if err := joinerOp.Join(&ep, key); err != nil {
			return errata.CommonError(err, fmt.Sprintf("can't join %#v as server_http.Endpoint with key '%s'", ep, key))
		}
	}

	return nil
}

func (c *Config) CompleteWithJoiner(joinerOp joiner.Operator, host string, port int, prefix string) error {
	if c == nil {
		return errors.New("no server_http.Config to be completed")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}
	c.Host, c.Port, c.Prefix = host, portStr, prefix

	var ok bool
	for key, ep := range c.EndpointsSettled {
		if ep.Endpoint, ok = joinerOp.Interface(ep.EndpointInternalKey).(Endpoint); ok {
			c.EndpointsSettled[key] = ep
		} else if handlerPtr, _ := joinerOp.Interface(ep.EndpointInternalKey).(*Endpoint); handlerPtr != nil {
			ep.Endpoint = *handlerPtr
			c.EndpointsSettled[key] = ep
		} else {
			return fmt.Errorf("no server_http.Endpoint with key %s", ep.EndpointInternalKey)
		}
	}

	return nil
}

func (c *Config) CompleteDirectly(endpoints Endpoints, host string, port int, prefix string) error {
	if c == nil {
		return errors.New("no server_http.Config to be completed")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}
	c.Host, c.Port, c.Prefix = host, portStr, prefix

	var ok bool
	for key, ep := range c.EndpointsSettled {
		if ep.Endpoint, ok = endpoints[ep.EndpointInternalKey]; ok {
			c.EndpointsSettled[key] = ep
		} else {
			return fmt.Errorf("no server_http.Endpoint with key %s", ep.EndpointInternalKey)
		}
	}

	return nil
}
