package auth_ecdsa

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_ecdsa"

func Starter() starter.Operator {
	return &identity_ecdsa{}
}

var l logger.Operator
var _ starter.Operator = &identity_ecdsa{}

type identity_ecdsa struct {
	interfaceKey joiner.InterfaceKey
}

func (ss *identity_ecdsa) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_ecdsa) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) (info []common.Map, err error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon

	// var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil, nil
}

func (ss *identity_ecdsa) Setup() error {
	return nil
}

func (ss *identity_ecdsa) Run(joinerOp joiner.Operator) error {
	identOp, err := New(1000, time.Second*2, nil)
	if err != nil {
		return errors.Wrap(err, "can't init identity_ecdsa.ActorKey")
	}

	err = joinerOp.Join(identOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join identity_ecdsa identOp as identity.ActorKey with key '%s'", ss.interfaceKey)
	}

	return nil
}
