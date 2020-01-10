package auth_ecdsa

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
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
	l = lCommon

	// var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil, nil
}

func (ss *identity_ecdsa) Setup() error {
	return nil
}

func (ss *identity_ecdsa) Run(joiner joiner.Operator) error {
	identOp, err := New(nil)
	if err != nil {
		return errors.Wrap(err, "can't init identity_ecdsa.Actor")
	}

	err = joiner.Join(identOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join identity_ecdsa identOp as identity.Actor with key '%s'", ss.interfaceKey)
	}

	return nil
}
