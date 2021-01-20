package server_http

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavlo67/workshop/common"

	"github.com/pavlo67/workshop/common/joiner"
)

type EndpointConfig struct {
	Path     string
	Tags     []string
	Produces []string

	Handler    *Endpoint
	HandlerKey joiner.InterfaceKey // for init purposes only
}

type Endpoints map[string]EndpointConfig

type Config struct {
	Title   string
	Version string
	Prefix  string
	Endpoints
}

type Swagger map[string]interface{}

// type SwaggerEndpoint struct {}

func (c Config) SwaggerV2(port string, noHTTPS bool) ([]byte, error) {
	paths := map[string]common.Map{} // map[string]map[string]map[string]interface{}{}

	for key, ep := range c.Endpoints {
		if ep.Handler == nil {
			continue
		}

		path := c.Prefix + ep.Handler.PathTemplateBraced(ep.Path)
		method := strings.ToLower(ep.Handler.Method)

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

		for _, pp := range ep.Handler.PathParams {
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
		for _, qp := range ep.Handler.QueryParams {
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
			if len(ep.Handler.BodyParams) > 0 {
				parameters = append(parameters, ep.Handler.BodyParams)
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
	if noHTTPS {
		schemes = []string{"http"}
	} else {
		schemes = []string{"https", "http"}
	}

	swagger := Swagger{
		"swagger": "2.0",
		"info": map[string]string{
			"title":   c.Title,
			"version": c.Version,
		},
		// "basePath": c.Prefix,
		"schemes": schemes,
		"port":    port,
		"paths":   paths,
	}

	return json.MarshalIndent(swagger, "", " ")
}
