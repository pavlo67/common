package server_http

import (
	"encoding/json"
)

type EndpointDescription struct {
	Method      string          `json:",omitempty"`
	PathParams  []string        `json:",omitempty"`
	QueryParams []string        `json:",omitempty"`
	BodyParams  json.RawMessage `json:",omitempty"`
}

type EndpointKey = string
type EndpointsSettled map[EndpointKey]EndpointSettled

type EndpointSettled struct {
	Path     string
	Tags     []string `json:",omitempty"`
	Produces []string `json:",omitempty"`
	Endpoint
}

type Endpoint struct {
	EndpointDescription
	WorkerHTTP
}
