package router

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "router"

type Key string

type WorkerFunc func(params Params, data []byte) (server.DataResponse, error)

type Operator interface {
	HandleWorker(endpoint Endpoint, worker WorkerFunc, allowedIDs []auth.ID)
	// RouteString(key Key, params []string) (string, error)
}

func InitEndpoints(op Operator, endpoints map[string]Endpoint, workers map[string]WorkerFunc, allowedIDs []auth.ID) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if worker, ok := workers[key]; ok {
			op.HandleWorker(ep, worker, allowedIDs)
		} else {
			errs = append(errs, errors.Errorf("no handler for endpoint: %s", key))
		}
	}

	return errs
}
