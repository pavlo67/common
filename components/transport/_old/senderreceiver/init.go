package senderreceiver

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/partes/connector/receiver"
	"github.com/pavlo67/partes/connector/sender"
)

const InterfaceKey joiner.InterfaceKey = "senderreceiverstub"

// starter.Operator -------------------------------------------------------------------------------------------

var _ starter.Operator = &senderstubStarter{}

func Starter() starter.Operator {
	return &senderstubStarter{}
}

type senderstubStarter struct{}

func (sc *senderstubStarter) Name() string {
	return string(InterfaceKey)
}

func (sc *senderstubStarter) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	return nil, nil
}

func (sc *senderstubStarter) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (sc *senderstubStarter) Init() error {

	senderOp, err := New()
	if err != nil {
		return errors.Wrap(err, "can't init senderreceiverstub.Operator")
	}

	err = joiner.JoinInterface(senderOp, sender.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join senderreceiverstub as sender.Operator")
	}

	err = joiner.JoinInterface(senderOp, receiver.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join senderreceiverstub as receiver.Operator")
	}

	return nil
}
