package server_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

type Config struct {
	Title            string
	Version          string
	Host             string
	Port             string
	Prefix           string
	EndpointsSettled map[joiner.InterfaceKey]EndpointSettled
}

func (c Config) EP(endpointKey joiner.InterfaceKey, params []string, createFullURL bool) (string, string, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return "", "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if len(ep.PathParams) != len(params) {
		return "", "", fmt.Errorf("wrong params list (%#v) for endpoint (%s / %#v)", params, endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix

	for i, param := range params {
		if param == "" {
			return "", "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(param)
	}

	return ep.Method, urlStr, nil
}

type Swagger map[string]interface{}

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

func (c Config) InitSwagger(isHTTPS bool, swaggerStaticFilePath string, l logger.Operator) error {
	swaggerJSON, err := c.SwaggerV2(isHTTPS)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(swaggerStaticFilePath, swaggerJSON, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerStaticFilePath, swaggerJSON, err)
	}
	l.Infof("%d bytes are written into %s", len(swaggerJSON), swaggerStaticFilePath)

	return nil
}

const onInitPages = "on server_http.HandlePages()"

func (c Config) HandlePages(srvOp Operator, l logger.Operator) error {
	if srvOp == nil {
		return errors.New(onInitPages + ": srvOp == nil")
	}

	for key, ep := range c.EndpointsSettled {
		if err := srvOp.HandleEndpoint(key, c.Prefix+ep.Path, ep.Endpoint); err != nil {
			return fmt.Errorf(onInitPages+": handling %s, %s, %#v got %s", key, ep.Path, ep.Endpoint, err)
		}
	}

	return nil
}

// joining endpoints -----------------------------------------------------

func (c *Config) CompleteWithJoiner(joinerOp joiner.Operator, host string, port int, prefix string) error {
	if c == nil {
		return errors.New("no server_http.Config to be completed")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}
	c.Host, c.Port, c.Prefix = host, portStr, prefix

	for key, ep := range c.EndpointsSettled {
		if endpoint, ok := joinerOp.Interface(key).(Endpoint); ok {
			ep.Endpoint = endpoint
			c.EndpointsSettled[key] = ep
		} else if endpointPtr, _ := joinerOp.Interface(key).(*Endpoint); endpointPtr != nil {
			ep.Endpoint = *endpointPtr
			c.EndpointsSettled[key] = ep
		} else {
			return fmt.Errorf("no server_http.Endpoint joined with key %s", key)
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

EP_SETTLED:
	for key, epSettled := range c.EndpointsSettled {
		// TODO??? use epSettled.InternalKey to correct the main key value

		for _, ep := range endpoints {
			if ep.InternalKey == key {
				epSettled.Endpoint = ep
				c.EndpointsSettled[key] = epSettled
				continue EP_SETTLED
			}
		}
		return fmt.Errorf("no server_http.Endpoint with key %s", key)
	}

	return nil
}
