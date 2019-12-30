package flow_server_http_handler

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data/data_tagged"

	"github.com/pavlo67/workshop/constructions/dataflow"
)

var dataTaggedOp data_tagged.Operator
var l logger.Operator

var _ starter.Operator = &flowTaggedServerHTTPStarter{}

type flowTaggedServerHTTPStarter struct {
	// interfaceKey joiner.DataInterfaceKey
}

func Starter() starter.Operator {
	return &flowTaggedServerHTTPStarter{}
}

func (ss *flowTaggedServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *flowTaggedServerHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", ss.Name()))
	}

	return nil, errs.Err()
}

func (ss *flowTaggedServerHTTPStarter) Setup() error {
	return nil
}

func (ss *flowTaggedServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	dataTaggedOp, ok = joinerOp.Interface(dataflow.InterfaceKey).(data_tagged.Operator)
	if !ok {
		return errors.Errorf("no data_tagged.Operator with key %s", dataflow.InterfaceKey)
	}

	err := joinerOp.Join(&listEndpoint, dataflow.ListInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listEndpoint as server_http.Endpoint with key '%s'", dataflow.ListInterfaceKey)
	}

	err = joinerOp.Join(&readEndpoint, dataflow.ReadInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join readEndpoint as server_http.Endpoint with key '%s'", dataflow.ReadInterfaceKey)
	}

	return nil
}
