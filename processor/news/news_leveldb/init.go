package news_leveldb

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/processor/news"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &news_leveldbStarter{}
}

var l logger.Operator

type news_leveldbStarter struct {
	path string
}

func (fl *news_leveldbStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fl *news_leveldbStarter) Prepare(cfg *config.PunctumConfig, params basis.Options) error {
	l = logger.Get()

	var errs basis.Errors
	fl.path, errs = cfg.Path("news.leveldb", nil)
	if fl.path == "" {
		return errors.New("no path to leveldb files described")
	}

	return errs.Err()
}

func (fl *news_leveldbStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (fl *news_leveldbStarter) Setup() error {
	return nil
}

func (fl *news_leveldbStarter) Init(joinerOp joiner.Operator) error {
	newsOp, err := New(fl.path)
	if err != nil {
		return errors.Wrap(err, "can't init newsOp")
	}

	err = joinerOp.JoinInterface(newsOp, news.InterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join newsOp as %s operator", news.InterfaceKey)
	}

	return nil
}
