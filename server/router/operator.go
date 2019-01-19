package router

import (
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "router"

type Key string

type Worker interface {
	Do(route string, routeParams server.RouteParams, data []byte) (server.BinaryResponse, error)
}

type Operator interface {
	SetRoute(key Key, route string, paramNames []string, worker Worker) error
	// RouteString(key Key, params []string) (string, error)
}
