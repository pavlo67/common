package auth_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &identity_login_stubStarter{}
}

type UserStub struct {
	ID       auth.ID
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

func (sc *identity_login_stubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sc *identity_login_stubStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	sc.interfaceKey = joiner.InterfaceKey(params.StringKeyDefault("interface_key", string(auth.InterfaceKey)))

	credentialsConf, errs := conf.Credentials(params.StringKeyDefault("config_credentials_key", "default"), nil)

	var ok bool

	if sc.users, ok = params["users"].([]UserStub); !ok || len(sc.users) < 1 {
		errs = append(errs, errors.New("no users defined for identity_login_stub.Starter (in params['users'])"))
	}

	if sc.salt, ok = credentialsConf["salt"]; !ok {
		errs = append(errs, errors.Wrapf(config.ErrNoValue, "no data for key 'salt' in config.credentials in %#v", credentialsConf))
	}

	return errs.Err()
}

func (sc *identity_login_stubStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (sc *identity_login_stubStarter) Setup() error {
	return nil
}

func (sc *identity_login_stubStarter) Init(joiner joiner.Operator) error {
	u, err := New(sc.users, sc.salt)
	if err != nil {
		return errors.Wrapf(err, "can't init identity_login_stubStarter")
	}

	err = joiner.JoinInterface(u, sc.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join identity_login_stubStarter as identity.Operator interface")
	}

	return nil
}
