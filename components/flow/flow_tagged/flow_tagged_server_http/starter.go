package flow_tagged_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data/data_tagged"
	"github.com/pavlo67/workshop/components/flow"
)

var flowTaggedOp data_tagged.Operator
var l logger.Operator

const Name = "flow_tagged_server_http"

var _ starter.Operator = &flowTaggedServerHTTPStarter{}

type flowTaggedServerHTTPStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &flowTaggedServerHTTPStarter{}
}

func (ss *flowTaggedServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *flowTaggedServerHTTPStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	return nil, errs.Err()
}

func (ss *flowTaggedServerHTTPStarter) Setup() error {
	return nil
}

func (ss *flowTaggedServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	flowTaggedOp, ok = joinerOp.Interface(flow.TaggedInterfaceKey).(data_tagged.Operator)
	if !ok {
		return errors.Errorf("no workspace.Operator with key %s", flow.TaggedInterfaceKey)
	}

	return nil
}
