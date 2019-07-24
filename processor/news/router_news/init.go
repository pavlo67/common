package router_news

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"

	"github.com/pavlo67/punctum/processor/news"
	"github.com/pavlo67/punctum/server/controller"
)

func Starter() starter.Operator {
	return &news_routerStarter{}
}

var l logger.Operator
var newsOp news.Operator

type news_routerStarter struct{}

func (dcs *news_routerStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dcs *news_routerStarter) Prepare(conf *config.Config, options, runtimeOptions basis.Info) error {
	l = logger.Get()

	return nil
}

func (dcs *news_routerStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (dcs *news_routerStarter) Setup() error {
	return nil
}

func (dcs *news_routerStarter) Init(joinerOp joiner.Operator) error {
	routerInterfaceKey := controller.InterfaceKey
	routerOp, ok := joinerOp.Interface(routerInterfaceKey).(controller.Operator)
	if !ok {
		return errors.Errorf("no router.Operator interface with key %s found for rss_router component", routerInterfaceKey)
	}

	newsInterfaceKey := news.InterfaceKey
	newsOp, ok = joinerOp.Interface(newsInterfaceKey).(news.Operator)
	if !ok {
		return errors.Errorf("no news.Operator interface with key %s found for rss_router component", newsInterfaceKey)
	}

	errs := controller.InitEndpoints(
		routerOp,
		endpoints,
		nil,
	)

	return errs.Err()
}
