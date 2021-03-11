package transformer_pack_persons_types01

import (
	"fmt"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/data_exchange/components/types01"
	"github.com/pavlo67/data_exchange/components/vcs"
)

var _ transformer.Operator = &transformerPackPersonsTypes01{}

// DEPRECATED: use the same taype from data_exchange package
type PackDescription struct {
	Title     string               `json:",omitempty" bson:",omitempty"`
	Fields    structures.Fields    `json:",omitempty" bson:",omitempty"`
	ErrorsMap structures.ErrorsMap `json:",omitempty" bson:",omitempty"`
	History   vcs.History          `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time            `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time           `json:",omitempty" bson:",omitempty"`
}

type PersonsPack struct {
	PackDescription
	Data []persons.Item
}

type transformerPackPersonsTypes01 struct {
	personsPack *PersonsPack
}

const onNew = "on transformerPackPersonsTypes01.New(): "

func New() (transformer.Operator, error) {
	return &transformerPackPersonsTypes01{}, nil
}

func (transformOp *transformerPackPersonsTypes01) Name() string {
	return string(InterfaceKey)
}

func (transformOp *transformerPackPersonsTypes01) Reset() error {
	transformOp.personsPack = nil
	return nil
}

const onStat = "on transformerPackPersonsTypes01.Stat(): "

func (transformOp *transformerPackPersonsTypes01) Stat(selector *selectors.Term, params common.Map) (interface{}, error) {

	return structures.PackStat{
		ItemsStat: structures.ItemsStat{
			Total:    len(transformOp.personsPack.Data),
			NonEmpty: len(transformOp.personsPack.Data),
			Errored:  0, // TODO!!!
		},
		FieldsStat: transformOp.personsPack.Fields.Stat(),
		ErrorsStat: transformOp.personsPack.ErrorsMap.Stat(),
	}, nil
}

const onIn = "on transformerPackPersonsTypes01.In(): "

func (transformOp *transformerPackPersonsTypes01) In(selector *selectors.Term, params common.Map, data interface{}) error {
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

	transformOp.personsPack = &PersonsPack{
		PackDescription: PackDescription{
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
		transformOp.personsPack.Data[i] = persons.Item{
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
		transformOp.personsPack.Data[i].SetCreds(creds)
	}

	return nil
}

const onOut = "on transformerPackPersonsTypes01.Out(): "

func (transformOp *transformerPackPersonsTypes01) Out(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	dataPack := structures.Pack{
		Title:     transformOp.personsPack.Title,
		Fields:    transformOp.personsPack.Fields,
		ErrorsMap: transformOp.personsPack.ErrorsMap,
		History:   transformOp.personsPack.History,
		CreatedAt: transformOp.personsPack.CreatedAt,
		UpdatedAt: transformOp.personsPack.UpdatedAt,
	}

	persons01 := make([]types01.Person, len(transformOp.personsPack.Data))

	for i, personsItem := range transformOp.personsPack.Data {
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

func (transformOp *transformerPackPersonsTypes01) Copy(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	return transformOp.personsPack, nil
}
