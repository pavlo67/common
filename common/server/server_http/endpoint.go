package server_http

import (
	"encoding/json"
	"fmt"
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
			return errors.CommonError(err, fmt.Sprintf("can't join %#v as server_http.Endpoint with key '%s'", ep, ep.InternalKey))
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
