package server_http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
)

type Endpoint struct {
	InternalKey joiner.InterfaceKey
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
}

var rePathParam = regexp.MustCompile(":[^/]+")

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

type Endpoints []Endpoint

func (eps Endpoints) Join(joinerOp joiner.Operator) error {
	for i, ep := range eps {
		if err := joinerOp.Join(&eps[i], ep.InternalKey); err != nil {
			return errors.CommonError(err, fmt.Sprintf("can't join %#v as Endpoint with key '%s'", ep, ep.InternalKey))
		}
	}

	return nil
}

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

// this trick allows to prevent run-time errors with wrong endpoint parameters number
// using CheckGet...() functions we move parameter number checks to initiation stage

type Get1 func(string) (string, error)
type Get2 func(string, string) (string, error)
type Get3 func(string, string, string) (string, error)
type Get4 func(string, string, string, string) (string, error)

func CheckGet0(c Config, endpointKey joiner.InterfaceKey, createFullURL bool) (string, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if strings.ToUpper(ep.Method) != "GET" {
		return "", fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	return urlStr, nil
}

func CheckGet1(c Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get1, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1 string) (string, error) {
		p1 = strings.TrimSpace(p1)
		if p1 == "" {
			return "", fmt.Errorf("empty param %s for endpoint (%s / %#v)", p1, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(p1)
		return urlStr, nil
	}, nil
}

func CheckGet2(c Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get2, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2 string) (string, error) {
		params := [2]string{p1, p2}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet3(c Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get3, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3 string) (string, error) {
		params := [3]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet4(c Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get4, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3, p4 string) (string, error) {
		params := [4]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}
