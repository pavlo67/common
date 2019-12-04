package server

import (
	"reflect"
	"strings"

	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

type Endpoint struct {
	Method     string `json:"method,omitempty"`
	ServerPath string `json:"server_path,omitempty"`
	LocalPath  string `json:"static_path,omitempty"`

	//Title       string                       `json:"title,omitempty"`
	//Parameters  []string                     `json:"parameters,omitempty"`
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

// TODO: process `*params`-templates also

func readEndpoint(e0 interface{}, localPath string) (*Endpoint, error) {
	var e1 []string
	if e, ok := e0.([]string); ok {
		e1 = e
	} else if e, ok := e0.([]interface{}); ok {
		eTmp, err := stringifySlice(e)
		if err != nil {
			return nil, errors.Wrapf(err, "bad Endpoint JSON: %#v", e0)
		}
		e1 = eTmp
	} else {
		return nil, errors.Errorf("bad Endpoint JSON type: %#v", reflect.TypeOf(e0))
	}

	if len(e1) == 1 {
		return &Endpoint{Method: "GET", ServerPath: e1[0]}, nil
	}

	if len(e1) >= 2 {
		ep := Endpoint{
			Method:     strings.ToUpper(e1[1]),
			ServerPath: e1[0],
		}
		if ep.Method == "FILE" {
			if len(e1) < 3 || strings.TrimSpace(e1[2]) == "" {
				return nil, errors.Errorf("bad FILE endpoint static path: %#v", e1)
			} else {
				ep.LocalPath = localPath + strings.TrimSpace(e1[2])
			}
		}
		return &ep, nil
	}
	return nil, errors.Errorf("bad Endpoint JSON length: %s", len(e1))
}

func stringifySlice(s0 []interface{}) ([]string, error) {
	var s1 []string
	for _, v0 := range s0 {
		if v, ok := v0.(string); ok {
			s1 = append(s1, v)
			//} else if v, ok := v0.(float64); ok {
			//	s1 = append(s1, strconv.FormatFloat(v, 10))
		} else {
			return nil, errors.Errorf("bad string value %v type %#v", v0, reflect.TypeOf(v0))
		}
	}
	return s1, nil
}
