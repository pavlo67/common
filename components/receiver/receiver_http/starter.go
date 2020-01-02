package receiver_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/receiver"
)

func Starter() starter.Operator {
	return &receiverHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &receiverHTTPStarter{}

type receiverHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (sh *receiverHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sh *receiverHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	sh.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(receiver.InterfaceKey)))

	return nil, nil
}

func (sh *receiverHTTPStarter) Setup() error {
	return nil
}

func (sh *receiverHTTPStarter) Run(joinerOp joiner.Operator) error {
	packsOp, ok := joinerOp.Interface(packs.InterfaceKey).(packs.Operator)
	if !ok {
		return errors.Errorf("no packs.Operator with key %s", packs.InterfaceKey)
	}

	err := joinerOp.Join(&receiveEndpoint, receiver.ActionInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join receiveEndpoint as server_http.Endpoint with key '%s'", receiver.ActionInterfaceKey)
	}

	receiverOp, err := New(packsOp)
	if err != nil {
		return errors.Wrap(err, "can't init receiver.Operator")
	}

	err = joinerOp.Join(receiverOp, sh.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *receiverHTTP as receiver.Operator with key '%s'", sh.interfaceKey)
	}

	return nil
}
