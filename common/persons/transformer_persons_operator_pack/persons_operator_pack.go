package transformer_persons_operator_pack

import (
	"fmt"

	"github.com/pavlo67/data_exchange/components/structures"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/transformer"
)

var _ transformer.Operator = &transformerOperatorPackPersons{}

type transformerOperatorPackPersons struct {
	personsOp persons.Operator
	identity  *auth.Identity

	packDescription *structures.PackDescription
}

const onNew = "on transformerOperatorPackPersons.New(): "

func New(personsOp persons.Operator, identity *auth.Identity) (transformer.Operator, error) {
	if personsOp == nil {
		return nil, errors.New(onNew + ": no persons.Operator")
	}
	//if personsOpCleaner == nil {
	//	return nil, errors.New(onNew + ": no persons.Cleaner")
	//}

	return &transformerOperatorPackPersons{
		personsOp: personsOp,
		identity:  identity,
	}, nil
}

func (transformOp *transformerOperatorPackPersons) Name() string {
	return string(InterfaceKey)
}

func (transformOp *transformerOperatorPackPersons) reset() error {
	transformOp.packDescription = nil

	// TODO!!! clean according to .PackURN

	return nil
}

const onStat = "on transformerOperatorPackPersons.Stat(): "

func (transformOp *transformerOperatorPackPersons) Stat(selector *selectors.Term, params common.Map) (interface{}, error) {
	return transformOp.personsOp.Stat(selector, transformOp.identity)
}

const onIn = "on transformerOperatorPackPersons.In(): "

func (transformOp *transformerOperatorPackPersons) In(params common.Map, data interface{}) error {
	if err := transformOp.reset(); err != nil {
		return errors.CommonError(err, onIn)
	}

	var personsPack *persons.Pack

	switch v := data.(type) {
	case persons.Pack:
		personsPack = &v
	case *persons.Pack:
		personsPack = v
	default:
		return fmt.Errorf(onIn+": wrong data to import: %#v", data)
	}

	if personsPack == nil {
		return fmt.Errorf(onIn+": nil data to import: %#v", data)
	}

	for i, personItem := range personsPack.Items {
		if _, err := transformOp.personsOp.Save(personItem, transformOp.identity); err != nil {
			return fmt.Errorf(onIn+": can't save item (%d / %#v), got %s", i, personItem, err)
		}
	}

	// TODO!!! save PackDescription persistently
	transformOp.packDescription = &personsPack.PackDescription

	return nil
}

const onOut = "on transformerOperatorPackPersons.Out(): "

func (transformOp *transformerOperatorPackPersons) Out(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	personsPack := persons.Pack{}

	// TODO!!! persistent pack descriptions storage is needed
	if transformOp.packDescription != nil {
		personsPack.PackDescription = *transformOp.packDescription
	}

	personsPack.Items, err = transformOp.personsOp.List(selector, transformOp.identity)
	if err != nil {
		return nil, fmt.Errorf(onOut+": can't list items (%#v), got %s", selector, err)
	}

	// TODO!!! set .URN for every item

	return personsPack, nil
}

func (transformOp *transformerOperatorPackPersons) Copy(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	return nil, common.ErrNotImplemented
}
