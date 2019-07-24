package founts_leveldb

import (
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

func (fl *founts_leveldbStarter) Prepare(cfg *config.Config, options, runtimeOptions basis.Info) error {
	l = logger.Get()

	var errs basis.Errors
	fl.path, errs = cfg.Path("founts.leveldb", nil)
	if fl.path == "" {
		return errors.New("no path to leveldb files described")
	}

	return errs.Err()
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
		return errors.Wrap(err, "can't init fountsOp")
	}

	err = joinerOp.JoinInterface(fountsOp, founts.InterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join fountsOp as %s operator", founts.InterfaceKey)
	}

	return nil
}
