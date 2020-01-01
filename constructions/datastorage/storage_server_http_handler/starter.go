package storage_server_http_handler

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data_tagged"

	"github.com/pavlo67/workshop/constructions/datastorage"
)

var dataTaggedOp data_tagged.Operator
var l logger.Operator

var _ starter.Operator = &dataTaggedServerHTTPStarter{}

type dataTaggedServerHTTPStarter struct {
	// interfaceKey joiner.DataInterfaceKey
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
		errs = append(errs, fmt.Errorf("no logger for %s:-(", ss.Name()))
	}

	// interfaceKey = joiner.DataInterfaceKey(options.StringDefault("interface_key", string(server_http.DataInterfaceKey)))

	return nil, errs.Err()
}

func (ss *dataTaggedServerHTTPStarter) Setup() error {
	return nil
}

func (ss *dataTaggedServerHTTPStarter) Run(joinerOp joiner.Operator) error {

	var ok bool
	dataTaggedOp, ok = joinerOp.Interface(datastorage.InterfaceKey).(data_tagged.Operator)
	if !ok {
		return errors.Errorf("no data_tagged.Operator with key %s", datastorage.InterfaceKey)
	}

	err := joinerOp.Join(&listEndpoint, datastorage.ListInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listEndpoint as server_http.Endpoint with key '%s'", datastorage.ListInterfaceKey)
	}

	err = joinerOp.Join(&readEndpoint, datastorage.ReadInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join readEndpoint as server_http.Endpoint with key '%s'", datastorage.ReadInterfaceKey)
	}

	err = joinerOp.Join(&saveEndpoint, datastorage.SaveInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join saveEndpoint as server_http.Endpoint with key '%s'", datastorage.SaveInterfaceKey)
	}

	err = joinerOp.Join(&removeEndpoint, datastorage.RemoveInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join removeEndpoint as server_http.Endpoint with key '%s'", datastorage.RemoveInterfaceKey)
	}

	err = joinerOp.Join(&countTagsEndpoint, datastorage.CountTagsInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join countTagsEndpoint as server_http.Endpoint with key '%s'", datastorage.CountTagsInterfaceKey)
	}

	err = joinerOp.Join(&listWithTagEndpoint, datastorage.ListWithTagInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listWithTagEndpoint as server_http.Endpoint with key '%s'", datastorage.ListWithTagInterfaceKey)
	}

	return nil
}
