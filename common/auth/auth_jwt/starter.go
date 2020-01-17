package auth_jwt

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_jwt"

func Starter() starter.Operator {
	return &identity_jwtStarter{}
}

var l logger.Operator
var _ starter.Operator = &identity_jwtStarter{}

type identity_jwtStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (ss *identity_jwtStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_jwtStarter) Init(cfgCommon, cfg *config.Config, l logger.Operator, options common.Map) (info []common.Map, err error) {
	l = logger.Get()

	// var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil, nil
}

func (ss *identity_jwtStarter) Setup() error {
	return nil
}

func (ss *identity_jwtStarter) Run(joiner joiner.Operator) error {
	identOp, err := New()
	if err != nil {
		return errors.Wrap(err, "can't init identity_jwt.Actor")
	}

	err = joiner.Join(identOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join identity_jwt identOp as identity.Actor with key '%s'", ss.interfaceKey)
	}

	return nil
}
