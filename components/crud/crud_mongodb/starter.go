package crud_mongodb

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/crud"
)

const Name = "crud_mongodb"

func Starter() starter.Operator {
	return &crudMongoDBStarter{}
}

var l logger.Operator
var _ starter.Operator = &crudMongoDBStarter{}

type crudMongoDBStarter struct {
	config       config.Access
	interfaceKey joiner.InterfaceKey
}

func (cm *crudMongoDBStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cm *crudMongoDBStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Options) ([]common.Options, error) {
	l = lCommon

	cfgMongoDB := config.Access{}
	err := cfg.Value("mongodb", &cfgMongoDB)
	if err != nil {
		return nil, err
	}

	cm.config = cfgMongoDB
	cm.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(crud.InterfaceKey)))

	return nil, nil
}

func (cm *crudMongoDBStarter) Setup() error {
	return nil
}

func (cm *crudMongoDBStarter) Run(joinerOp joiner.Operator) error {

	// TODO!!!
	crudOp, _, _, err := NewCRUD(&cm.config, 5*time.Second, "crud", crud.Item{})

	err = joinerOp.Join(crudOp, cm.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *flowSQLite as flow.Operator with key '%s'", cm.interfaceKey)
	}

	return nil
}
