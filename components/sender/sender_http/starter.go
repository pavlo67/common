package sender_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/router"
	"github.com/pavlo67/workshop/components/sender"
)

func Starter() starter.Operator {
	return &senderHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &senderHTTPStarter{}

type senderHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (sh *senderHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sh *senderHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	sh.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(sender.InterfaceKey)))
	return nil, nil
}

func (sh *senderHTTPStarter) Setup() error {
	return nil
}

func (sh *senderHTTPStarter) Run(joinerOp joiner.Operator) error {
	routerOp, ok := joinerOp.Interface(router.InterfaceKey).(router.Operator)
	if !ok {
		return errors.Errorf("no router.Operator with key %s", router.InterfaceKey)
	}

	packsOp, ok := joinerOp.Interface(packs.InterfaceKey).(packs.Operator)
	if !ok {
		return errors.Errorf("no packs.Operator with key %s", packs.InterfaceKey)
	}

	senderOp, err := New(packsOp, routerOp)
	if err != nil {
		return errors.Wrap(err, "can't init sender.Operator")
	}

	err = joinerOp.Join(senderOp, sh.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *senderHTTP as sender.Operator with key '%s'", sh.interfaceKey)
	}

	return nil
}
