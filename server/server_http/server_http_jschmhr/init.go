package node_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/basis/starter"
)

const InterfaceKey program.InterfaceKey = "node_http_jschmhr"

func Starter() starter.Operator {
	return &node_http_jschmhrStarter{}
}

var l *zap.SugaredLogger
var _ starter.Operator = &node_http_jschmhrStarter{}

type node_http_jschmhrStarter struct {
	interfaceKey program.InterfaceKey

	config config.ServerTLS
}

func (ss *node_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *node_http_jschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	var errs basis.Errors

	ss.interfaceKey = program.InterfaceKey(params.StringKeyDefault("interface_key", string(InterfaceKey)))

	ss.config, errs = conf.Server(params.StringKeyDefault("config_key", "data"), errs)
	if ss.config.Port <= 0 {
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %d", ss.config.Port))
	}

	return errs.Err()
}

func (ss *node_http_jschmhrStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (ss *node_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *node_http_jschmhrStarter) Init(joiner program.Joiner) error {
	nodeOp, err := New(
		ss.config.Port,
		ss.config.TLSCertFile,
		ss.config.TLSKeyFile,
	)
	if err != nil {
		return errors.Wrap(err, "can't init node_http_jschmhr.Operator")
	}

	err = joiner.JoinInterface(nodeOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join node_http_jschmhr nodeOp as node.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
