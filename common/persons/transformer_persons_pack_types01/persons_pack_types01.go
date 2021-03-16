package transformer_persons_pack_types01

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/data_exchange/components/types/types01"
)

var _ transformer.Operator = &transformerPersonsPackTypes01{}

type transformerPersonsPackTypes01 struct {
	packPersons *persons.Pack
}

const onNew = "on transformerPersonsPackTypes01.New(): "

func New() (transformer.Operator, error) {
	return &transformerPersonsPackTypes01{}, nil
}

func (transformOp *transformerPersonsPackTypes01) Name() string {
	return string(InterfaceKey)
}

func (transformOp *transformerPersonsPackTypes01) reset() error {
	transformOp.packPersons = nil
	return nil
}

const onStat = "on transformerPersonsPackTypes01.Stat(): "

func (transformOp *transformerPersonsPackTypes01) Stat(selector *selectors.Term, params common.Map) (interface{}, error) {
	return &structures.PackStat{
		ItemsStat: structures.ItemsStat{
			Total:    len(transformOp.packPersons.Items),
			NonEmpty: len(transformOp.packPersons.Items),
			Errored:  0, // TODO!!!
		},
		FieldsStat: transformOp.packPersons.Fields.Stat(),
		ErrorsStat: transformOp.packPersons.ErrorsMap.Stat(),
	}, nil
}

const onIn = "on transformerPersonsPackTypes01.In(): "

func (transformOp *transformerPersonsPackTypes01) In(params common.Map, data interface{}) error {
	if err := transformOp.reset(); err != nil {
		return errors.CommonError(err, onIn)
	}

	var dataPack *structures.Pack

	if data != nil {
		switch v := data.(type) {
		case structures.Pack:
			dataPack = &v
		case *structures.Pack:
			dataPack = v
		default:
			return fmt.Errorf("wrong data to import: %#v", data)
		}
	}

	if dataPack == nil {
		return fmt.Errorf("nil data to import: %#v", data)
	}

	var persons01 []types01.Person

	switch v := dataPack.Items.(type) {
	case []types01.Person:
		persons01 = v
	case *[]types01.Person:
		if v == nil {
			return fmt.Errorf("nil dataPack.Items to import: %#v", dataPack)
		}
		persons01 = *v
	default:
		return fmt.Errorf("wrong dataPack.Items to import: %#v", dataPack.Items)
	}

	transformOp.packPersons = &persons.Pack{
		PackDescription: dataPack.PackDescription,
		Items:           make([]persons.Item, len(persons01)),
	}

	for i, person01 := range persons01 {
		transformOp.packPersons.Items[i] = persons.Item{
			Identity: auth.Identity{
				Nickname: person01.Nickname,
				Roles:    person01.Roles,
			},
			ItemDescription: person01.ItemDescription,
		}
		transformOp.packPersons.Items[i].SetCreds(person01.Creds)
	}

	return nil
}

const onOut = "on transformerPersonsPackTypes01.Out(): "

func (transformOp *transformerPersonsPackTypes01) Out(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	dataPack := structures.Pack{
		PackDescription: transformOp.packPersons.PackDescription,
	}

	persons01 := make([]types01.Person, len(transformOp.packPersons.Items))

	for i, personsItem := range transformOp.packPersons.Items {

		persons01[i] = types01.Person{
			Nickname:        personsItem.Nickname,
			Roles:           personsItem.Roles,
			Creds:           personsItem.Creds(),
			ItemDescription: personsItem.ItemDescription,
		}
	}
	dataPack.Items = persons01

	return dataPack, nil
}

func (transformOp *transformerPersonsPackTypes01) Copy(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	return transformOp.packPersons, nil
}
