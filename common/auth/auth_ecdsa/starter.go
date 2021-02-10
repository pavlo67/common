package auth_ecdsa

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey common.InterfaceKey = "auth_ecdsa"

func Starter() starter.Operator {
	return &auth_ecdsa{}
}

var l logger.Operator
var _ starter.Operator = &auth_ecdsa{}

type auth_ecdsa struct {
	interfaceKey common.InterfaceKey
}

func (ss *auth_ecdsa) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *auth_ecdsa) Prepare(cfg *config.Config, options common.Map) error {
	ss.interfaceKey = common.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (ss *auth_ecdsa) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	identOp, err := New()
	if err != nil {
		return err
	}

	if err = joinerOp.Join(identOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_ecdsa identOp as auth.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
