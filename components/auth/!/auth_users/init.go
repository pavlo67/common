package auth_users

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/partes/connector/sender"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/users"
)

const InterfaceKey joiner.InterfaceKey = "authusers"

func Starter(useMessages bool) starter.Operator {
	return &authusersComponent{
		useMessages: useMessages,
	}
}

type authusersComponent struct {
	conf         config.Config
	interfaceKey joiner.InterfaceKey
	useMessages  bool
	salt         string
}

var _ starter.Operator = &authusersComponent{}

func (cm *authusersComponent) Name() string {
	return string(InterfaceKey)
}

func (cm *authusersComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	index, errs := config.ComponentIndex(indexPath, filelib.CurrentPath(), nil)
	if len(errs) > 0 {
		return nil, errs.Err()
	}

	var credentialsConf map[string]string
	credentialsConf, errs = conf.Credentials("", errs)

	var ok bool
	cm.salt, ok = credentialsConf["salt"]
	if !ok {
		errs = append(errs, errors.Wrapf(config.ErrNoValue, "no data for key 'salt' in config.credentials in %#v", credentialsConf))
	}

	if index.Params["interfaceKey"] != "" {
		cm.interfaceKey = joiner.InterfaceKey(index.Params["interfaceKey"])
	} else {
		cm.interfaceKey = auth.InterfaceKey
	}

	return nil, nil
}

func (cm *authusersComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (cm *authusersComponent) Init() error {

	usersOp, ok := joiner.Component(users.InterfaceKey).(users.Operator)
	if !ok {
		return errors.New("no users.Operator for authusers interface found :-(")
	}

	senderOp, ok := joiner.Component(sender.InterfaceKey).(sender.Operator)
	if !ok {
		return errors.New("no sender.Operator for authusers interface found :-(")
	}

	// TODO: set specific domain

	authOp, err := New(usersOp, senderOp, cm.salt, cm.useMessages)
	if err != nil {
		return errors.Wrap(err, "can't init authusers.Operator interface")
	}

	err = joiner.JoinInterface(authOp, auth.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join authusers as auth.Operator interface")
	}

	return nil
}
