package persons_fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/GehirnInc/crypt/sha256_crypt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/selectors"
)

var _ persons.Operator = &personsFSStub{}

type personsFSStub struct {
	randInt int
	path    string
}

const onNew = "on personsFSStub.New() "

func New(cfg config.Access) (persons.Operator, db.Cleaner, error) {
	path, err := filelib.Dir(cfg.Path)
	if err != nil {
		return nil, nil, errors.CommonError(err, onNew)
	}

	personsOp := personsFSStub{path: path}
	return &personsOp, &personsOp, nil
}

const onSave = "on personsFSStub.Save()"

func (personsOp *personsFSStub) Save(item persons.Item, identity *auth.Identity) (auth.ID, error) {
	if identity == nil || (item.ID != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item})
	}

	var path string

	if item.ID != "" {
		path = filepath.Join(personsOp.path, string(item.ID))

		itemOld, err := personsOp.read(item.ID)
		if err != nil || itemOld == nil {
			errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
			l.Error(errorStr)
			return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item, "requestedRole": rbac.RoleAdmin})
		}

		item.CreatedAt = itemOld.CreatedAt
		now := time.Now()
		item.UpdatedAt = &now

	} else {
		personsOp.randInt++
		item.ID = auth.ID(strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(personsOp.randInt))

		path = filepath.Join(personsOp.path, string(item.ID))
		if _, err := os.Stat(path); err == nil {
			return "", errors.CommonError(common.DuplicateUserKey, common.Map{"on": onSave, "item": item})
		}
	}

	if err := personsOp.write(path, item); err != nil {
		return "", errors.Wrap(err, onSave)
	}

	return item.ID, nil
}

const onRemove = "on personsFSStub.Remove()"

func (personsOp *personsFSStub) Remove(id auth.ID, identity *auth.Identity) error {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return errors.CommonError(common.NoRightsKey, common.Map{"on": onRemove, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	path := filepath.Join(personsOp.path, string(id)) //  personsOp.path + string(id)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf(onRemove+": can't os.RemoveAll(%s), got  %s", path, err)
	}

	return nil
}

const onRead = "on personsFSStub.Read()"

func (personsOp *personsFSStub) Read(id auth.ID, identity *auth.Identity) (*persons.Item, error) {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onRead, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	return personsOp.read(id)
}

// read/write file ----------------------------------------------

type PersonWithCreds struct {
	persons.Item
	auth.Creds
}

func (personsOp *personsFSStub) write(path string, item persons.Item) error {
	personWithCreds := PersonWithCreds{
		item,
		item.Creds(),
	}

	jsonBytes, err := json.Marshal(personWithCreds)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, jsonBytes, 0644)
}

func (personsOp *personsFSStub) read(id auth.ID) (*persons.Item, error) {
	path := filepath.Join(personsOp.path, string(id)) //  personsOp.path + string(id)
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	var personWithCreds PersonWithCreds
	if err := json.Unmarshal(jsonBytes, &personWithCreds); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	personWithCreds.Item.ID = id
	personWithCreds.Item.SetCreds(personWithCreds.Creds)

	//for k, v := range personWithCreds.Creds {
	//	personWithCreds.SetCredsByKey(k, v)
	//}

	return &personWithCreds.Item, nil
}

const onList = "on personsFSStub.List(): "

func (personsOp *personsFSStub) List(Selector *selectors.Term, identity *auth.Identity) ([]persons.Item, error) {
	if !identity.HasRole(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
	}

	d, err := os.Open(personsOp.path)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	var items []persons.Item
	for _, name := range names {
		item, err := personsOp.read(auth.ID(name))
		if err != nil || item == nil {
			return nil, fmt.Errorf(onList+": got %#v, %s", item, err)
		}
		// delete(item.Creds, auth.CredsPasshash)

		items = append(items, *item)
	}

	return items, nil
}

func (personsOp *personsFSStub) Stat(*selectors.Term, *auth.Identity) (db.StatMap, error) {
	return nil, common.ErrNotImplemented
}
