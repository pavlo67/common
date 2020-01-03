package router_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/router"
)

func Starter() starter.Operator {
	return &routerStubStarter{}
}

var l logger.Operator
var _ starter.Operator = &routerStubStarter{}

type routerStubStarter struct {
	interfaceKey joiner.InterfaceKey
	// cleanerInterfaceKey joiner.HandlerKey

	routes router.Routes
}

func (rs *routerStubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rs *routerStubStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	err := cfgCommon.Value("routes", &rs.routes)
	if err != nil {
		return nil, err
	}

	rs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(router.InterfaceKey)))
	// rs.cleanerInterfaceKey = joiner.HandlerKey(options.StringDefault("cleaner_interface_key", string(router.CleanerInterfaceKey)))

	return nil, nil
}

func (rs *routerStubStarter) Setup() error {
	return nil
}

func (rs *routerStubStarter) Run(joinerOp joiner.Operator) error {
	routerOp, _, err := New(rs.routes)
	if err != nil {
		return errors.Wrap(err, "can't init router.Operator")
	}

	err = joinerOp.Join(routerOp, rs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *tagsSQLite as router.Operator with key '%s'", rs.interfaceKey)
	}

	return nil
}
