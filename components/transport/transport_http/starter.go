package transport_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/router"
	"github.com/pavlo67/workshop/components/transport"
)

func Starter() starter.Operator {
	return &transportHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &transportHTTPStarter{}

type transportHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
	handlerKey   joiner.InterfaceKey
}

func (th *transportHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (th *transportHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	th.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(transport.InterfaceKey)))
	th.handlerKey = joiner.InterfaceKey(options.StringDefault("handler_key", string(transport.HandlerInterfaceKey)))

	return nil, nil
}

func (th *transportHTTPStarter) Setup() error {
	return nil
}

func (th *transportHTTPStarter) Run(joinerOp joiner.Operator) error {
	routerOp, ok := joinerOp.Interface(router.InterfaceKey).(router.Operator)
	if !ok {
		return errors.Errorf("no router.Operator with key %s", router.InterfaceKey)
	}

	packsOp, ok := joinerOp.Interface(packs.InterfaceKey).(packs.Operator)
	if !ok {
		return errors.Errorf("no packs.Operator with key %s", packs.InterfaceKey)
	}

	transpOp, err := New(packsOp, routerOp)
	if err != nil {
		return errors.Wrap(err, "can'th init transport.Operator")
	}

	err = joinerOp.Join(transpOp, th.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can'th join *transportHTTP as transport.Operator with key '%s'", th.interfaceKey)
	}

	return nil
}
