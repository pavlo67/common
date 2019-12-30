package server_http

import (
	"encoding/json"
	"strings"

	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pkg/errors"
)

type EndpointConfig struct {
	Path     string
	Tags     []string
	Produces []string

	Handler      *Endpoint
	InterfaceKey joiner.InterfaceKey // for init purposes only
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

func (c Config) SwaggerV2(port string) ([]byte, error) {
	paths := map[string]map[string]map[string]interface{}{}

	for key, ep := range c.Endpoints {
		if ep.Handler == nil {
			continue
		}

		path := c.Prefix + ep.Handler.PathTemplateBraced(ep.Path)
		method := strings.ToLower(ep.Handler.Method)

		epDescr := map[string]interface{}{
			"operationId": key,
			"tags":        ep.Tags,
		}

		if len(ep.Produces) >= 1 {
			epDescr["produces"] = ep.Produces
		} else {
			epDescr["produces"] = []string{"application/json"}
		}

		var parameters []map[string]interface{}

		for _, pp := range ep.Handler.PathParams {
			parameters = append(
				parameters,
				map[string]interface{}{
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
				map[string]interface{}{
					"in":          "query",
					"required":    false, // TODO!!!
					"name":        qp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}

		if method == "post" {
			parameters = append(
				parameters,
				map[string]interface{}{
					"in":       "body",
					"required": true,
					"name":     "body_item",
					"type":     "string",
				},
			)
		}

		if len(parameters) > 0 {
			epDescr["parameters"] = parameters
		}

		if epDescrPrev, ok := paths[path][method]; ok {
			return nil, errors.Errorf("duplicate endpoint description (%s/%s): %#v vs. %#v", path, method, epDescrPrev, epDescr)
		}
		if _, ok := paths[path]; ok { // pathPrev
			paths[path][method] = epDescr
		} else {
			paths[path] = map[string]map[string]interface{}{method: epDescr}
		}
	}

	swagger := Swagger{
		"swagger": "2.0",
		"info": map[string]string{
			"title":   c.Title,
			"version": c.Version,
		},
		// "basePath": c.Prefix,
		"schemes": []string{"http", "https"},
		"port":    port,
		"paths":   paths,
	}

	return json.MarshalIndent(swagger, "", " ")
}
