package auth_jwt

import (
	"fmt"
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
	interfaceKey joiner.InterfaceKey
	// interfaceSetCredsKey  joiner.InterfaceKey
	keyPath string
}

func (ss *identity_jwtStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_jwtStarter) Prepare(cfg *config.Config, options common.Map) error {

	var cfgServerHTTP server.Config
	if err := cfg.Value("server_http", &cfgServerHTTP); err != nil {
		return err
	}

	// var errs basis.Errors
	ss.keyPath = strings.TrimSpace(cfgServerHTTP.KeyPath)
	if ss.keyPath == "" {
		ss.keyPath = filelib.CurrentPath()
	} else if ss.keyPath[len(ss.keyPath)-1] != '/' {
		ss.keyPath += "/"
	}

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	// ss.interfaceSetCredsKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (ss *identity_jwtStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	identOp, err := New(ss.keyPath + "jwt.key")
	if err != nil || identOp == nil {
		return errors.Wrap(err, "can't init identity_jwt.ActorKey")
	}

	if err = joinerOp.Join(identOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceKey)
	}

	//if ss.interfaceKey != ss.interfaceSetCredsKey {
	//	if err = joinerOp.Join(identOp, ss.interfaceSetCredsKey); err != nil {
	//		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceSetCredsKey)
	//	}
	//}

	return nil
}
