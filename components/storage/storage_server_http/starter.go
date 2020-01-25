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
	dataKey joiner.InterfaceKey
	// interfaceKey joiner.DataInterfaceKey
}

func Starter() starter.Operator {
	return &dataTaggedServerHTTPStarter{}
}

func (dtsh *dataTaggedServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dtsh *dataTaggedServerHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", dtsh.Name())
	}

	dtsh.dataKey = joiner.InterfaceKey(options.StringDefault("data_key", string(storage.InterfaceKey)))
	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.DataInterfaceKey)))

	return nil, nil
}

func (dtsh *dataTaggedServerHTTPStarter) Setup() error {
	return nil
}

func (dtsh *dataTaggedServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	var ok bool
	dataTaggedOp, ok = joinerOp.Interface(dtsh.dataKey).(datatagged.Operator)
	if !ok {
		return errors.Errorf("no data_tagged.ActorKey with key %s", storage.InterfaceKey)
	}

	err := joinerOp.Join(&recentEndpoint, storage.RecentInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join recentEndpoint as server_http.Endpoint with key '%s'", storage.RecentInterfaceKey)
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

	err = joinerOp.Join(&listTagsEndpoint, storage.ListTagsInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listTagsEndpoint as server_http.Endpoint with key '%s'", storage.ListTagsInterfaceKey)
	}

	err = joinerOp.Join(&listTaggedEndpoint, storage.ListTaggedInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join listTaggedEndpoint as server_http.Endpoint with key '%s'", storage.ListTaggedInterfaceKey)
	}

	return nil
}
