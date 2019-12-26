package filesloader_http

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/constructions/filesloader"
	"github.com/pkg/errors"
)

func Starter() starter.Operator {
	return &loaderHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &loaderHTTPStarter{}

type loaderHTTPStarter struct {
	interfaceKey joiner.InterfaceKey

	// TODO: use proxies
}

// ------------------------------------------------

func (lh *loaderHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (lh *loaderHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	lh.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(filesloader.InterfaceKey)))

	return nil, nil
}

func (lh *loaderHTTPStarter) Setup() error {
	return nil
}

func (lh *loaderHTTPStarter) Run(joinerOp joiner.Operator) error {

	dataOp, _, err := New(ts.config, ts.table, ts.interfaceKey, taggerOp, cleanerOp)
	if err != nil {
		return errors.Wrap(err, "can't init data.Operator")
	}

	err = joinerOp.Join(dataOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *dataSQLite as data.Operator with key '%s'", ts.interfaceKey)
	}

	return nil
}
