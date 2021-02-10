package auth_jwt

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey common.InterfaceKey = "auth_jwt"

func Starter() starter.Operator {
	return &identity_jwtStarter{}
}

var l logger.Operator
var _ starter.Operator = &identity_jwtStarter{}

type identity_jwtStarter struct {
	interfaceKey common.InterfaceKey
	// interfaceSetCredsKey  joiner.InterfaceKey
	keyPath string
}

func (ss *identity_jwtStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_jwtStarter) Prepare(cfg *config.Config, options common.Map) error {

	var cfgAuthJWT common.Map
	if err := cfg.Value("auth_jwt", &cfgAuthJWT); err != nil {
		return err
	}

	// var errs basis.Errors
	ss.keyPath = strings.TrimSpace(cfgAuthJWT.StringDefault("key_path", ""))
	if ss.keyPath == "" {
		ss.keyPath = "./"
	} else if ss.keyPath[len(ss.keyPath)-1] != '/' {
		ss.keyPath += "/"
	}

	ss.interfaceKey = common.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	// ss.interfaceSetCredsKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (ss *identity_jwtStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	authOp, err := New(ss.keyPath + "jwt.key")
	if err != nil || authOp == nil {
		return errata.CommonError(err, fmt.Sprintf("got %#v", authOp))
	}

	if err = joinerOp.Join(authOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceKey)
	}

	//if ss.interfaceKey != ss.interfaceSetCredsKey {
	//	if err = joinerOp.Join(authOp, ss.interfaceSetCredsKey); err != nil {
	//		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceSetCredsKey)
	//	}
	//}

	return nil
}
