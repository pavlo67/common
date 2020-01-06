package transport_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transportrouter"
)

func Starter() starter.Operator {
	return &transportHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &transportHTTPStarter{}

type transportHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
	handlerKey   joiner.InterfaceKey

	domain identity.Domain
}

func (th *transportHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (th *transportHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	th.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(transport.InterfaceKey)))
	th.handlerKey = joiner.InterfaceKey(options.StringDefault("handler_key", string(transport.HandlerInterfaceKey)))
	th.domain = identity.Domain(options.StringDefault("domain", ""))

	return nil, nil
}

func (th *transportHTTPStarter) Setup() error {
	return nil
}

func (th *transportHTTPStarter) Run(joinerOp joiner.Operator) error {
	routerOp, ok := joinerOp.Interface(transportrouter.InterfaceKey).(transportrouter.Operator)
	if !ok {
		return errors.Errorf("no router.Operator with key %s", transportrouter.InterfaceKey)
	}

	packsOp, ok := joinerOp.Interface(packs.InterfaceKey).(packs.Operator)
	if !ok {
		return errors.Errorf("no packs.Operator with key %s", packs.InterfaceKey)
	}

	transpOp, receiveEndpoint, err := New(packsOp, routerOp, th.domain)
	if err != nil {
		return errors.Wrap(err, "can'th init transport.Operator")
	}

	err = joinerOp.Join(transpOp, th.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can'th join *transportHTTP as transport.Operator with key '%s'", th.interfaceKey)
	}

	err = joinerOp.Join(receiveEndpoint, th.handlerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join receiveEndpoint as server_http.Endpoint with key '%s'", th.handlerKey)
	}

	return nil
}
