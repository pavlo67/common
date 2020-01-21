package exporter_data

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/exporter"
)

func Starter() starter.Operator {
	return &exporterDataStarter{}
}

var l logger.Operator
var _ starter.Operator = &exporterDataStarter{}

type exporterDataStarter struct {
	dataKey      joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
}

func (ed *exporterDataStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ed *exporterDataStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	if lCommon == nil {
		return nil, errors.New("no logger")
	}
	l = lCommon

	ed.dataKey = joiner.InterfaceKey(options.StringDefault("data_key", string(data.InterfaceKey)))
	ed.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(exporter.InterfaceKey)))

	return nil, nil
}

func (ed *exporterDataStarter) Setup() error {
	return nil
}

func (ed *exporterDataStarter) Run(joinerOp joiner.Operator) error {
	dataOp, ok := joinerOp.Interface(ed.dataKey).(data.Operator)
	if !ok {
		return errors.Errorf("no data.Operator with key %s", ed.dataKey)
	}

	exporterOp, err := New(dataOp, ed.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init *exporterData as exporter.Operator")
	}

	err = joinerOp.Join(exporterOp, ed.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *exporterData as exporter.Operator with key '%s'", ed.interfaceKey)
	}

	return nil
}
