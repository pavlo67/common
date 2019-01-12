package founts_leveldb

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/processor/founts"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &founts_leveldbStarter{}
}

var l logger.Operator

type founts_leveldbStarter struct {
	path string
}

func (fl *founts_leveldbStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fl *founts_leveldbStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	fl.path = strings.TrimSpace(params.StringKeyDefault("path", ""))
	if fl.path == "" {
		return errors.New("no path to leveldb files described")
	}

	return nil
}

func (fl *founts_leveldbStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (fl *founts_leveldbStarter) Setup() error {
	return nil
}

func (fl *founts_leveldbStarter) Init(joinerOp joiner.Operator) error {
	fountsOp, err := New(fl.path)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = joinerOp.JoinInterface(fountsOp, founts.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
