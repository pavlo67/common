package transformer_persons_operator_pack

import (
	"fmt"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "transformer_operator_persons_pack"

func Starter() starter.Operator {
	return &transformerPackPersonsOperatorStarter{}
}

// ---------------------------------------------------------------------------------

var l logger.Operator
var _ starter.Operator = &transformerPackPersonsOperatorStarter{}

type transformerPackPersonsOperatorStarter struct {
	personsKey   joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
}

func (tppos *transformerPackPersonsOperatorStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (tppos *transformerPackPersonsOperatorStarter) Prepare(cfg *config.Config, options common.Map) error {
	tppos.personsKey = joiner.InterfaceKey(options.StringDefault("persons_key", string(persons.InterfaceKey)))
	tppos.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (tppos *transformerPackPersonsOperatorStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	personsOp, _ := joinerOp.Interface(tppos.personsKey).(persons.Operator)
	if personsOp == nil {
		return fmt.Errorf("no persons.Operator with key %s", tppos.personsKey)
	}

	transformOp, err := New(personsOp, auth.IdentityWithRoles(rbac.RoleAdmin))
	if err != nil {
		return err
	}

	if err = joinerOp.Join(transformOp, tppos.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *transformerOperatorPackPersons as transform.Operator with key '%s'", tppos.interfaceKey))
	}

	return nil
}
