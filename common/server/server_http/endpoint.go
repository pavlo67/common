package server_http

import (
	"regexp"
	"strings"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pkg/errors"
)

func InitEndpoints(cfg Config, srvOp Operator, l logger.Operator) error {
	if srvOp == nil {
		return errors.New("on .InitEndpoints(): srvOp == nil")
	}

	for _, ep := range cfg.Endpoints {
		err := srvOp.HandleEndpoint(ep.Key, cfg.Prefix+ep.Path, ep.Endpoint)
		if err != nil {
			return errors.Errorf("on .srvOp.HandleEndpoint(%s, %s, %#v): %s", ep.Key, ep.Path, ep.Endpoint, err)
		}
	}

	return nil
}

type Endpoint struct {
	Method      string   `json:"method,omitempty"`
	PathParams  []string `json:"path_params,omitempty"`
	QueryParams []string `json:"query_params,omitempty"`

	WorkerHTTP

	// AllowedIDs []common.Key `json:"allowed_ids,omitempty"`
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
//		path = path[:matches[nm][0]] + url.PathEscape(strings.ReplaceTags(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
//	}
//
//	return path
//}

func (ep Endpoint) PathTemplate(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ep.PathParams) < 1 {
		return serverPath
	}

	return serverPath + "/:" + strings.Join(ep.PathParams, "/:")
}

func (ep Endpoint) PathTemplateBraced(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ep.PathParams) < 1 {
		return serverPath
	}

	return serverPath + "/{" + strings.Join(ep.PathParams, "}/{") + "}"
}
