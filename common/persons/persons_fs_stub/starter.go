package persons_fs_stub

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &personsFSStubStarter{}
}

const configKeyDefault = "persons_fs"

type personsFSStubStarter struct {
	interfaceKey        joiner.InterfaceKey
	interfaceCleanerKey joiner.InterfaceKey

	cfg config.Access
}

var _ starter.Operator = &personsFSStubStarter{}
var l logger.Operator

func (uks *personsFSStubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (uks *personsFSStubStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	if lCommon == nil {
		return nil, errata.New("no logger.Operator")
	}
	l = lCommon

	configKey := options.StringDefault("config_key", configKeyDefault)

	if err := cfg.Value(configKey, &uks.cfg); err != nil {
		return nil, err
	}

	uks.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(persons.InterfaceKey)))
	uks.interfaceCleanerKey = joiner.InterfaceKey(options.StringDefault("interface_cleaner_key", string(persons.InterfaceCleanerKey)))

	return nil, nil
}

func (uks *personsFSStubStarter) Run(joiner joiner.Operator) error {
	personsOp, personsCleanerOp, err := New(uks.cfg)
	if err != nil {
		return errata.Wrapf(err, "can't personsFSStub.New()")
	}

	if err = joiner.Join(personsOp, uks.interfaceKey); err != nil {
		return errata.Wrap(err, "can't join *personsFSStub{} as persons.Operator interface")
	}

	if err = joiner.Join(personsCleanerOp, uks.interfaceCleanerKey); err != nil {
		return errata.Wrap(err, "can't join *personsFSStub{} as crud.Cleaner interface")
	}

	return nil
}
