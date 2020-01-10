package transport_pop3

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/partes/connector/receiver"
)

const InterfaceKey joiner.InterfaceKey = "receiverpop3"

func Starter() starter.Operator {
	return &receiverpop3Component{}
}

type receiverpop3Component struct {
	serverAccessConfig config.ServerAccess
}

func (sc *receiverpop3Component) Name() string {
	return string(InterfaceKey)
}

func (sc *receiverpop3Component) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	var errs basis.Errors

	sc.serverAccessConfig, errs = conf.POP3(partKeys["pop3"], nil)
	return nil, errs.Err()
}

func (sc *receiverpop3Component) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (sc *receiverpop3Component) Init() error {
	receiverOp, err := NewPOP3BytBox(sc.serverAccessConfig)
	if err != nil {
		return errors.Wrap(err, "can't init receiverpop3")
	}

	err = joiner.JoinInterface(receiverOp, receiver.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join receiverpop3 as receiver.Actor")
	}
	return nil
}
