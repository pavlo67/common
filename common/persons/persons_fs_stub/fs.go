package persons_fs_stub

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/libraries/filelib"
	"github.com/pavlo67/common/common/persons"
)

var _ persons.Operator = &personsFSStub{}

type personsFSStub struct {
	path string
}

const onNew = "on personsFSStub.New() "

func New(cfg config.Access) (persons.Operator, crud.Cleaner, error) {
	path, err := filelib.Dir(cfg.Path)
	if err != nil {
		return nil, nil, errors.CommonError(err, onNew)
	}

	personsOp := personsFSStub{path: path}

	return &personsOp, &personsOp, nil
}

const onAdd = "on personsFSStub.Add()"

func (pfs *personsFSStub) Add(identity auth.Identity, data common.Map, options *crud.Options) (auth.ID, error) {

	// TODO!!! check RBAC with options

	if personOld, _ := pfs.Read(identity.ID, options); personOld != nil {
		return "", errors.KeyableError(errors.DuplicateUserErr, nil)
	}

	person := persons.Item{
		Identity:  identity,
		Data:      data,
		CreatedAt: time.Now(),
	}
	// TODO!!! generate ID & password
	// TODO!!! save on FS

	return person.ID, nil
}

const onChange = "on personsFSStub.Change()"

func (pfs *personsFSStub) Change(persons.Item, *crud.Options) error {
	return errors.NotImplemented
}

const onList = "on personsFSStub.List(): "

func (pfs *personsFSStub) List(options *crud.Options) ([]persons.Item, error) {
	return nil, errors.NotImplemented
}

const onRemove = "on personsFSStub.Remove()"

func (pfs *personsFSStub) Remove(auth.ID, *crud.Options) error {
	return errors.NotImplemented
}

const onRead = "on personsFSStub.Read()"

func (pfs *personsFSStub) Read(auth.ID, *crud.Options) (*persons.Item, error) {
	return nil, errors.NotImplemented
}
