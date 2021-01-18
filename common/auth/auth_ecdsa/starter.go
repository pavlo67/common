package auth_ecdsa

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_ecdsa"

func Starter() starter.Operator {
	return &auth_ecdsa{}
}

var l logger.Operator
var _ starter.Operator = &auth_ecdsa{}

type auth_ecdsa struct {
	interfaceKey joiner.InterfaceKey
}

func (ss *auth_ecdsa) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *auth_ecdsa) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) (info []common.Map, err error) {
	if l = lCommon; lCommon == nil {
		return nil, errors.New("no logger")
	}

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil, nil
}

func (ss *auth_ecdsa) Setup() error {
	return nil
}

func (ss *auth_ecdsa) Run(joinerOp joiner.Operator) error {
	identOp, err := New()
	if err != nil {
		return err
	}

	if err = joinerOp.Join(identOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_ecdsa identOp as auth.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
