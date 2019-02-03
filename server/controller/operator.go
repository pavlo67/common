package controller

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "controller"

type Key string

type Worker func(endpoint Endpoint, params basis.Params, options basis.Options, data interface{}) (*server.DataResponse, error)

type Operator interface {
	HandleWorker(endpoint Endpoint, worker Worker, allowedIDs []auth.ID)
	// RouteString(key Key, params []string) (string, error)
}

func InitEndpoints(op Operator, endpoints map[string]Endpoint, allowedIDs []auth.ID) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if ep.Worker != nil {
			op.HandleWorker(ep, ep.Worker, allowedIDs)
		} else {
			errs = append(errs, errors.Errorf("no handler for endpoint: %s", key))
		}
	}

	return errs
}
