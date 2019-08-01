package server_http

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pavlo67/constructor/auth"
)

type Endpoint struct {
	Method     string    `json:"method,omitempty"`
	Path       string    `json:"path"`
	ParamNames []string  `json:"param_names,omitempty"`
	AllowedIDs []auth.ID `json:"allowed_ids,omitempty"`

	DataItem interface{} `json:"data_item,omitempty"` // for Worker

	// Shortcut string `json:"shortcut,omitempty"`

}

var rePathParam = regexp.MustCompile(":[^/]+")

func (ep Endpoint) WithParams(params ...string) string {
	matches := rePathParam.FindAllStringSubmatchIndex(ep.Path, -1)

	numMatches := len(matches)
	if len(params) < numMatches {
		numMatches = len(params)
	}

	path := ep.Path
	for nm := numMatches - 1; nm >= 0; nm-- {
		path = path[:matches[nm][0]] + url.PathEscape(strings.Replace(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
	}

	return path
}
