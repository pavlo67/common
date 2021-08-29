package auth_stub

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &authstubStarter{}
}

var l logger.Operator
var _ starter.Operator = &authstubStarter{}

type authstubStarter struct {
	defaultUser  config.Access
	interfaceKey joiner.InterfaceKey
}

func (ahs *authstubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ahs *authstubStarter) Prepare(cfg *config.Config, options common.Map) error {

	cfg.Value("auth_stub", &ahs.defaultUser)

	ahs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil
}

func (ahs *authstubStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	authOp, err := New(ahs.defaultUser)
	if err != nil {
		return errors.Wrap(err, "can't init *authstub{} as auth.Operator")
	}

	if err = joinerOp.Join(authOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *authstub{} as auth.Operator with key '%s'", ahs.interfaceKey)
	}

	return nil
}
