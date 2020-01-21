package storage_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/storage"
)

var dataTaggedOp datatagged.Operator
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
	dataTaggedOp, ok = joinerOp.Interface(storage.InterfaceKey).(datatagged.Operator)
	if !ok {
		return errors.Errorf("no data_tagged.ActorKey with key %s", storage.InterfaceKey)
	}

	err := joinerOp.Join(&listEndpoint, storage.ListInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listEndpoint as server_http.Endpoint with key '%s'", storage.ListInterfaceKey)
	}

	err = joinerOp.Join(&readEndpoint, storage.ReadInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join readEndpoint as server_http.Endpoint with key '%s'", storage.ReadInterfaceKey)
	}

	err = joinerOp.Join(&saveEndpoint, storage.SaveInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join saveEndpoint as server_http.Endpoint with key '%s'", storage.SaveInterfaceKey)
	}

	err = joinerOp.Join(&removeEndpoint, storage.RemoveInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join removeEndpoint as server_http.Endpoint with key '%s'", storage.RemoveInterfaceKey)
	}

	err = joinerOp.Join(&countTagsEndpoint, storage.CountTagsInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join countTagsEndpoint as server_http.Endpoint with key '%s'", storage.CountTagsInterfaceKey)
	}

	err = joinerOp.Join(&listWithTagEndpoint, storage.ListWithTagInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listWithTagEndpoint as server_http.Endpoint with key '%s'", storage.ListWithTagInterfaceKey)
	}

	return nil
}
