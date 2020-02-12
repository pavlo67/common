package storage_server_http

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/exporter"
	"github.com/pavlo67/workshop/components/storage"
)

var dataOp data.Operator
var exporterOp exporter.Operator
var l logger.Operator

var _ starter.Operator = &dataServerHTTPStarter{}

type dataServerHTTPStarter struct {
	dataKey     joiner.InterfaceKey
	exporterKey joiner.InterfaceKey
	// interfaceKey joiner.DataInterfaceKey
}

func Starter() starter.Operator {
	return &dataServerHTTPStarter{}
}

func (dtsh *dataServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dtsh *dataServerHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", dtsh.Name())
	}

	dtsh.dataKey = joiner.InterfaceKey(options.StringDefault("data_key", string(storage.InterfaceKey)))
	dtsh.exporterKey = joiner.InterfaceKey(options.StringDefault("exporter_key", string(exporter.InterfaceKey)))

	return nil, nil
}

func (dtsh *dataServerHTTPStarter) Setup() error {
	return nil
}

func (dtsh *dataServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	dataOp, _ = joinerOp.Interface(dtsh.dataKey).(data.Operator)
	if dataOp == nil {
		return errors.Errorf("no data_tagged.Operator with key %s", dtsh.dataKey)
	}

	exporterOp, _ = joinerOp.Interface(dtsh.exporterKey).(exporter.Operator)
	if exporterOp == nil {
		return errors.Errorf("no exporter.Operator with key %s", dtsh.exporterKey)
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

	err = joinerOp.Join(&exportEndpoint, storage.ExportInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join exportEndpoint as server_http.Endpoint with key '%s'", storage.ExportInterfaceKey)
	}

	return nil
}
