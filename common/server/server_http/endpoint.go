package server_http

import (
	"regexp"
	"strings"
)

type Endpoint struct {
	Method     string   `json:"method,omitempty"`
	ParamNames []string `json:"param_names,omitempty"`

	WorkerHTTP

	// AllowedIDs []common.ID `json:"allowed_ids,omitempty"`
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
//		path = path[:matches[nm][0]] + url.PathEscape(strings.Replace(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
//	}
//
//	return path
//}

func (ep Endpoint) PathTemplate(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ep.ParamNames) < 1 {
		return serverPath
	}

	return serverPath + "/:" + strings.Join(ep.ParamNames, "/:")
}
