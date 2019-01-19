package server_http

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
)

type Endpoint struct {
	Method     string   `json:"method"`
	ServerPath string   `json:"server_path,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

var rePathParam = regexp.MustCompile(":[^/]+")

func (ep Endpoint) Path(params ...string) string {
	matches := rePathParam.FindAllStringSubmatchIndex(ep.ServerPath, -1)

	numMatches := len(matches)
	if len(params) < numMatches {
		numMatches = len(params)
	}

	path := ep.ServerPath
	for nm := numMatches - 1; nm >= 0; nm-- {
		path = path[:matches[nm][0]] + url.PathEscape(strings.Replace(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
	}

	return path
}

func (ep Endpoint) Params() []string {
	return ep.Parameters
}

func InitEndpoints(op Operator, endpoints map[string]Endpoint, htmlHandlers map[string]HTMLHandler, restHandlers map[string]RESTHandler, binaryHandlers map[string]BinaryHandler, allowedIDs []auth.ID) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if htmlHandler, ok := htmlHandlers[key]; ok {
			op.HandleFuncHTML(ep.Method, ep.ServerPath, htmlHandler, allowedIDs...)
		} else if restHandler, ok := restHandlers[key]; ok {
			op.HandleFuncREST(ep.Method, ep.ServerPath, restHandler, allowedIDs...)
		} else if binaryHandler, ok := binaryHandlers[key]; ok {
			op.HandleFuncBinary(ep.Method, ep.ServerPath, binaryHandler, allowedIDs...)
		} else {
			errs = append(errs, errors.New("no handler for endpoint: "+key))
		}
	}

	return errs
}
