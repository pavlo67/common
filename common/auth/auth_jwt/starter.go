package auth_jwt

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pkg/errors"

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

func (ss *identity_jwtStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) (info []common.Map, err error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon

	// var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil, nil
}

func (ss *identity_jwtStarter) Setup() error {
	return nil
}

func (ss *identity_jwtStarter) Run(joinerOp joiner.Operator) error {
	identOp, err := New(filelib.CurrentPath() + "jwt.key")
	if err != nil {
		return errors.Wrap(err, "can't init identity_jwt.ActorKey")
	}

	err = joinerOp.Join(identOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join identity_jwt identOp as identity.ActorKey with key '%s'", ss.interfaceKey)
	}

	return nil
}
