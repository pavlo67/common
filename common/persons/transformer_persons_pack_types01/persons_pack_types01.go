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

func (transformOp *transformerPersonsPackTypes01) Reset() error {
	transformOp.packPersons = nil
	return nil
}

const onStat = "on transformerPersonsPackTypes01.Stat(): "

func (transformOp *transformerPersonsPackTypes01) Stat(selector *selectors.Term, params common.Map) (interface{}, error) {
	return &structures.PackStat{
		ItemsStat: structures.ItemsStat{
			Total:    len(transformOp.packPersons.Data),
			NonEmpty: len(transformOp.packPersons.Data),
			Errored:  0, // TODO!!!
		},
		FieldsStat: transformOp.packPersons.Fields.Stat(),
		ErrorsStat: transformOp.packPersons.ErrorsMap.Stat(),
	}, nil
}

const onIn = "on transformerPersonsPackTypes01.In(): "

func (transformOp *transformerPersonsPackTypes01) In(selector *selectors.Term, params common.Map, data interface{}) error {
	if err := transformOp.Reset(); err != nil {
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

	switch v := dataPack.Data.(type) {
	case []types01.Person:
		persons01 = v
	case *[]types01.Person:
		if v == nil {
			return fmt.Errorf("nil dataPack.Data to import: %#v", dataPack)
		}
		persons01 = *v
	default:
		return fmt.Errorf("wrong dataPack.Data to import: %#v", dataPack.Data)
	}

	transformOp.packPersons = &persons.Pack{
		PackDescription: structures.PackDescription{
			Title:     dataPack.Title,
			Fields:    dataPack.Fields,
			ErrorsMap: dataPack.ErrorsMap,
			History:   dataPack.History,
			CreatedAt: dataPack.CreatedAt,
			UpdatedAt: dataPack.UpdatedAt,
		},
		Data: make([]persons.Item, len(persons01)),
	}

	for i, person01 := range persons01 {
		transformOp.packPersons.Data[i] = persons.Item{
			Identity: auth.Identity{
				URN:      person01.URN,
				Nickname: person01.Nickname,
				Roles:    person01.Roles,
			},
			Data:      person01.Data,
			History:   person01.History, // TODO??? modify here
			CreatedAt: person01.CreatedAt,
			UpdatedAt: person01.UpdatedAt,
		}

		creds := auth.Creds{}
		for c := range person01.Creds {
			creds[auth.CredsType(c)] = person01.Creds.StringDefault(c, "")
		}
		transformOp.packPersons.Data[i].SetCreds(creds)
	}

	return nil
}

const onOut = "on transformerPersonsPackTypes01.Out(): "

func (transformOp *transformerPersonsPackTypes01) Out(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	dataPack := structures.Pack{
		PackDescription: structures.PackDescription{
			Title:     transformOp.packPersons.Title,
			Fields:    transformOp.packPersons.Fields,
			ErrorsMap: transformOp.packPersons.ErrorsMap,
			History:   transformOp.packPersons.History,
			CreatedAt: transformOp.packPersons.CreatedAt,
			UpdatedAt: transformOp.packPersons.UpdatedAt,
		},
	}

	persons01 := make([]types01.Person, len(transformOp.packPersons.Data))

	for i, personsItem := range transformOp.packPersons.Data {
		creds := common.Map{}

		for c, v := range personsItem.Creds() {
			creds[string(c)] = v
		}

		persons01[i] = types01.Person{
			URN:       personsItem.URN,
			Nickname:  personsItem.Nickname,
			Roles:     personsItem.Roles,
			Creds:     creds,
			Data:      personsItem.Data,
			History:   personsItem.History,
			CreatedAt: personsItem.CreatedAt,
			UpdatedAt: personsItem.UpdatedAt,
		}
	}
	dataPack.Data = persons01

	return dataPack, nil
}

func (transformOp *transformerPersonsPackTypes01) Copy(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	return transformOp.packPersons, nil
}
