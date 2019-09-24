package auth_users_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/auth"
)

const InterfaceKey joiner.InterfaceKey = "auth_user_sqlite"

func Starter() starter.Operator {
	return &auth_user_sqliteStarter{}
}

type UserSQLite struct {
	ID       common.ID
	Login    string
	Password string
}

var UserStubDefault = UserSQLite{
	ID:       "1",
	Login:    "aaa",
	Password: "bbb",
}

type auth_user_sqliteStarter struct {
	interfaceKey joiner.InterfaceKey

	users []UserSQLite
	// salt  string
}

var _ starter.Operator = &auth_user_sqliteStarter{}
var l logger.Operator

var credentialsConf = map[string]string{}

func (sc *auth_user_sqliteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sc *auth_user_sqliteStarter) Init(conf *config.Config, params common.Map) (info []common.Map, err error) {
	l = logger.Get()

	sc.interfaceKey = joiner.InterfaceKey(params.StringDefault("interface_key", string(auth.InterfaceKey)))

	// var ok bool
	var errs common.Errors

	//if sc.users, ok = params["users"].([]UserSQLite); !ok || len(sc.users) < 1 {
	//	errs = append(errs, errors.New("no users defined for auth_stub.Starter (in params['users'])"))
	//}

	sc.users = []UserSQLite{UserStubDefault}

	//if sc.salt, ok = credentialsConf["salt"]; !ok {
	//	errs = append(errs, errors.Wrapf(config.ErrNoValue, "no data for key 'salt' in config.credentials in %#v", credentialsConf))
	//}

	return nil, errs.Err()
}

func (sc *auth_user_sqliteStarter) Setup() error {
	return nil
}

func (sc *auth_user_sqliteStarter) Run(joiner joiner.Operator) error {
	u, err := New(sc.users, "") // sc.salt
	if err != nil {
		return errors.Wrapf(err, "can't init auth_user_sqliteStarter")
	}

	err = joiner.Join(u, sc.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join auth_user_sqliteStarter as identity.Operator interface")
	}

	return nil
}
