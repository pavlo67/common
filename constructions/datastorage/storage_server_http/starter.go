package storage_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data/data_tagged"
	"github.com/pavlo67/workshop/constructions/datastorage"
)

var dataTaggedOp data_tagged.Operator
var l logger.Operator

const Name = "data_tagged_server_http"

var _ starter.Operator = &dataTaggedServerHTTPStarter{}

type dataTaggedServerHTTPStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &dataTaggedServerHTTPStarter{}
}

func (ss *dataTaggedServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *dataTaggedServerHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *dataTaggedServerHTTPStarter) Setup() error {
	return nil
}

func (ss *dataTaggedServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	dataTaggedOp, ok = joinerOp.Interface(datastorage.TaggedInterfaceKey).(data_tagged.Operator)
	if !ok {
		return errors.Errorf("no storage.Operator with key %s", datastorage.TaggedInterfaceKey)
	}

	return nil
}
