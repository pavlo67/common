package auth_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/config"
	"github.com/pavlo67/constructor/components/common/joiner"
	"github.com/pavlo67/constructor/components/common/logger"
	"github.com/pavlo67/constructor/components/common/starter"
)

func Starter() starter.Operator {
	return &identity_login_stubStarter{}
}

type UserStub struct {
	ID       common.ID
	Login    string
	Password string
}

type identity_login_stubStarter struct {
	interfaceKey joiner.InterfaceKey

	users []UserStub
	salt  string
}

var _ starter.Operator = &identity_login_stubStarter{}
var l logger.Operator

var credentialsConf = map[string]string{}

func (sc *identity_login_stubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sc *identity_login_stubStarter) Init(conf *config.Config, params common.Info) (info []common.Info, err error) {
	l = logger.Get()

	sc.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(auth.InterfaceKey)))

	var ok bool
	var errs common.Errors

	if sc.users, ok = params["users"].([]UserStub); !ok || len(sc.users) < 1 {
		errs = append(errs, errors.New("no users defined for identity_login_stub.Starter (in params['users'])"))
	}

	if sc.salt, ok = credentialsConf["salt"]; !ok {
		errs = append(errs, errors.Wrapf(config.ErrNoValue, "no data for key 'salt' in config.credentials in %#v", credentialsConf))
	}

	return nil, errs.Err()
}

func (sc *identity_login_stubStarter) Setup() error {
	return nil
}

func (sc *identity_login_stubStarter) Run(joiner joiner.Operator) error {
	u, err := New(sc.users, sc.salt)
	if err != nil {
		return errors.Wrapf(err, "can't init identity_login_stubStarter")
	}

	err = joiner.Join(u, sc.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join identity_login_stubStarter as identity.Operator interface")
	}

	return nil
}
