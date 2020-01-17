package auth_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_stub"

func Starter() starter.Operator {
	return &auth_stubStarter{}
}

var UserStubDefault = UserStub{
	Key:          "1",
	Nickname:     "pavlo",
	PasswordHash: "123",
}

type auth_stubStarter struct {
	interfaceKey joiner.InterfaceKey

	users []UserStub
	// salt  string
}

var _ starter.Operator = &auth_stubStarter{}
var l logger.Operator

var credentialsConf = map[string]string{}

func (sc *auth_stubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sc *auth_stubStarter) Init(cfgCommon, cfg *config.Config, l logger.Operator, options common.Map) (info []common.Map, err error) {
	l = logger.Get()

	sc.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	// var ok bool
	var errs common.Errors

	//if sc.users, ok = params["users"].([]UserStub); !ok || len(sc.users) < 1 {
	//	errs = append(errs, errors.New("no users defined for auth_stub.Starter (in params['users'])"))
	//}

	sc.users = []UserStub{UserStubDefault}

	//if sc.salt, ok = credentialsConf["salt"]; !ok {
	//	errs = append(errs, errors.Wrapf(config.ErrNoValue, "no data for key 'salt' in config.credentials in %#v", credentialsConf))
	//}

	return nil, errs.Err()
}

func (sc *auth_stubStarter) Setup() error {
	return nil
}

func (sc *auth_stubStarter) Run(joiner joiner.Operator) error {
	u, err := New(sc.users, "") // sc.salt
	if err != nil {
		return errors.Wrapf(err, "can't init auth_stubStarter")
	}

	err = joiner.Join(u, sc.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join auth_stubStarter as identity.Actor interface")
	}

	return nil
}
