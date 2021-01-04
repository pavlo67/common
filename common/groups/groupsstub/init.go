package groupsstub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "groupsstub"

func Starter() starter.Operator {
	return &groupsstubStarter{}
}

type groupsstubStarter struct{}

func (sc *groupsstubStarter) Name() string {
	return string(InterfaceKey)
}

func (sc *groupsstubStarter) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	return nil, nil
}

func (sc *groupsstubStarter) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (sc *groupsstubStarter) Init() error {
	ctrlOp, err := New(nil, basis.Anyone)
	if err != nil {
		return errors.Wrap(err, "can't init groupsstub")
	}

	err = joiner.JoinInterface(ctrlOp, groups.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join groupsstub as groups.Operator")
	}
	return nil
}
