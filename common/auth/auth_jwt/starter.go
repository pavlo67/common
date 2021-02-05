package auth_jwt

import (
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_jwt"

func Starter() starter.Operator {
	return &identity_jwtStarter{}
}

var l logger.Operator
var _ starter.Operator = &identity_jwtStarter{}

type identity_jwtStarter struct {
	interfaceKey         joiner.InterfaceKey
	interfaceSetCredsKey joiner.InterfaceKey
	keyPath              string
}

func (ss *identity_jwtStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_jwtStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) (info []common.Map, err error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon

	var cfgServerHTTP server.Config
	err = cfg.Value("server_http", &cfgServerHTTP)
	if err != nil {
		return nil, err
	}

	// var errs basis.Errors
	ss.keyPath = strings.TrimSpace(cfgServerHTTP.KeyPath)
	if ss.keyPath == "" {
		ss.keyPath = filelib.CurrentPath()
	} else if ss.keyPath[len(ss.keyPath)-1] != '/' {
		ss.keyPath += "/"
	}

	// ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceJWTKey)))
	// ss.interfaceSetCredsKey = joiner.InterfaceKey(options.StringDefault("interface_set_creds_key", string(auth.InterfaceJWTKey)))

	return nil, nil
}

func (ss *identity_jwtStarter) Setup() error {
	return nil
}

func (ss *identity_jwtStarter) Run(joinerOp joiner.Operator) error {
	identOp, err := New(ss.keyPath + "jwt.key")
	if err != nil || identOp == nil {
		return errors.Wrap(err, "can't init identity_jwt.ActorKey")
	}

	if err = joinerOp.Join(identOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceKey)
	}

	if ss.interfaceKey != ss.interfaceSetCredsKey {
		if err = joinerOp.Join(identOp, ss.interfaceSetCredsKey); err != nil {
			return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceSetCredsKey)
		}
	}

	return nil
}
