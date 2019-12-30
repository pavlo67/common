package sendersmtp

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/connector/sender"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

func Starter() starter.Operator {
	return &sendergomailComponent{}
}

type sendergomailComponent struct {
	smtpConfig   config.ServerAccess
	senderConfig map[string]string
}

const InterfaceKey joiner.InterfaceKey = "sendergomail"

func (sc *sendergomailComponent) Name() string {
	return string(InterfaceKey)
}

func (sc *sendergomailComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	var errs basis.Errors

	sc.smtpConfig, errs = conf.SMTP("", errs)
	sc.senderConfig, errs = conf.Sender("", errs)
	return nil, errs.Err()
}

func (sc *sendergomailComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (sc *sendergomailComponent) Init() error {

	senderOp, err := New(sc.smtpConfig, sc.senderConfig)
	if err != nil {
		return errors.Wrap(err, "can't init sendergomail")
	}

	err = joiner.JoinInterface(senderOp, sender.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join sendergomail as sender.Operator")
	}
	return nil
}
