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

type Worker func(user *auth.User, params basis.Params, data interface{}) (server.Response, error)

type Operator interface {
	HandleWorker(endpoint Endpoint, worker Worker)

	Start()
}

func InitEndpoints(op Operator, endpoints map[string]Endpoint, workers map[string]Worker) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if worker, ok := workers[key]; ok {
			op.HandleWorker(ep, worker)
		} else {
			errs = append(errs, errors.Errorf("no handler for endpoint: %s", key))
		}
	}

	return errs
}
