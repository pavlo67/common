package transportrouter_stub

import (
	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/components/transportrouter"
)

var _ transportrouter.Operator = &routesStub{}

// var _ crud.Cleaner = &tagsSQLite{}

type routesStub struct {
	routes transportrouter.Routes
}

const onNew = "on routesStub.New(): "

func New(routes transportrouter.Routes) (transportrouter.Operator, crud.Cleaner, error) {
	routerOp := routesStub{
		routes: routes,
	}

	return &routerOp, nil, nil
}

const onRoutes = "on routesStub.Routes(): "

func (routerOp *routesStub) Routes() (transportrouter.Routes, error) {
	return routerOp.routes, nil
}
