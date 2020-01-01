package router_stub

import (
	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/components/router"
)

var _ router.Operator = &routesStub{}

// var _ crud.Cleaner = &tagsSQLite{}

type routesStub struct {
	routes router.Routes
}

const onNew = "on routesStub.New(): "

func New(routes router.Routes) (router.Operator, crud.Cleaner, error) {
	routerOp := routesStub{
		routes: routes,
	}

	return &routerOp, nil, nil
}

const onRoutes = "on routesStub.Routes(): "

func (routerOp *routesStub) Routes() (router.Routes, error) {
	return routerOp.routes, nil
}
