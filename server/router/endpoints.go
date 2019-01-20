package router

import (
	"net/url"
	"regexp"
	"strings"
)

type Endpoint struct {
	Method     string   `json:"method"`
	ServerPath string   `json:"server_path,omitempty"`
	ParamNames []string `json:"param_names,omitempty"`
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
	return ep.ParamNames
}
