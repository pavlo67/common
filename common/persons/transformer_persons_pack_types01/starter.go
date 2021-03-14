package transformer_persons_pack_types01

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "transformer_pack_persons_types01"

func Starter() starter.Operator {
	return &transformerPersonsPackTypes01Starter{}
}

// ---------------------------------------------------------------------------------

var l logger.Operator
var _ starter.Operator = &transformerPersonsPackTypes01Starter{}

type transformerPersonsPackTypes01Starter struct {
	interfaceKey joiner.InterfaceKey
}

func (tppt01s *transformerPersonsPackTypes01Starter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (tppt01s *transformerPersonsPackTypes01Starter) Prepare(cfg *config.Config, options common.Map) error {
	tppt01s.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (tppt01s *transformerPersonsPackTypes01Starter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	transformOp, err := New()
	if err != nil {
		return err
	}

	if err = joinerOp.Join(transformOp, tppt01s.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *transformerPersonsPackTypes01 as transform.Operator with key '%s'", tppt01s.interfaceKey))
	}

	return nil
}
