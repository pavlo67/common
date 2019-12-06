package data_mongodb

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
)

const Name = "data_mongodb"

func Starter() starter.Operator {
	return &dataMongoDBStarter{}
}

var l logger.Operator
var _ starter.Operator = &dataMongoDBStarter{}

type dataMongoDBStarter struct {
	config       config.Access
	interfaceKey joiner.InterfaceKey
}

func (cm *dataMongoDBStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cm *dataMongoDBStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	cfgMongoDB := config.Access{}
	err := cfg.Value("mongodb", &cfgMongoDB)
	if err != nil {
		return nil, err
	}

	cm.config = cfgMongoDB
	cm.interfaceKey = joiner.InterfaceKey(options.StringDefault(joiner.InterfaceKeyFld, string(data.InterfaceKey)))

	return nil, nil
}

func (cm *dataMongoDBStarter) Setup() error {
	return nil
}

func (cm *dataMongoDBStarter) Run(joinerOp joiner.Operator) error {

	// TODO!!!
	dataOp, _, _, err := NewData(&cm.config, 5*time.Second, cm.config.Path, "data", data.Item{})

	err = joinerOp.Join(dataOp, cm.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *flowSQLite as flow.Operator with key '%s'", cm.interfaceKey)
	}

	return nil
}
