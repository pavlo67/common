package auth_users

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/common/users"
)

const InterfaceKey joiner.InterfaceKey = "auth_users"

var l logger.Operator

var _ starter.Operator = &authPassUsersStarter{}

type authPassUsersStarter struct {
	interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &authPassUsersStarter{}
}

func (apu *authPassUsersStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (apu *authPassUsersStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon

	apu.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil, nil
}

func (apu *authPassUsersStarter) Setup() error {
	return nil
}

const onRun = "on authPassUsersStarter.Run(): "

func (apu *authPassUsersStarter) Run(joinerOp joiner.Operator) error {

	usersOp, ok := joinerOp.Interface(users.InterfaceKey).(users.Operator)
	if !ok {
		return errors.Errorf(onRun+"no users.Operator with key %s", users.InterfaceKey)
	}

	authOp, err := New(usersOp, 10, "")
	if err != nil {
		return errors.Wrap(err, onRun)
	}

	err = joinerOp.Join(authOp, apu.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, onRun+"can't join authPassUsers{} as auth.Operator with key '%s'", apu.interfaceKey)
	}

	return nil
}
