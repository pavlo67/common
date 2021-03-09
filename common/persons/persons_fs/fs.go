package persons_fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/GehirnInc/crypt/sha256_crypt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
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

func (pfs *personsFSStub) Add(identity auth.Identity, creds auth.Creds, data common.Map, options *crud.Options) (auth.ID, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return "", errors.KeyableError(common.NoRightsKey, common.Map{"on": onAdd, "identity": identity, "data": data, "requestedRole": rbac.RoleAdmin})
	}

	idStr := strings.TrimSpace(string(identity.ID))
	if idStr == "" {
		idStr = strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(rand.Int())
	}

	path := filepath.Join(pfs.path, idStr) //  pfs.path + string(id)
	if _, err := os.Stat(path); err == nil {
		return "", errors.KeyableError(common.DuplicateUserKey, common.Map{"on": onAdd, "identity": identity, "data": data})
	}

	person := persons.Item{
		Identity:  identity,
		Data:      data,
		CreatedAt: time.Now(),
	}
	person.SetCreds(creds)

	if err := pfs.write(path, person); err != nil {
		return "", errors.Wrap(err, onAdd)
	}

	return auth.ID(idStr), nil
}

const onChange = "on personsFSStub.Change()"

func (pfs *personsFSStub) Change(item persons.Item, options *crud.Options) (*persons.Item, error) {
	if options == nil || options.Identity == nil {
		return nil, errors.KeyableError(common.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	itemOld, err := pfs.read(item.Identity.ID)
	if err != nil || itemOld == nil {
		errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
		if options.HasRole(rbac.RoleAdmin) {
			return nil, errors.KeyableError(common.WrongIDKey, common.Map{"on": onChange, "item": item, "reason": errorStr})
		} else {
			l.Error(errorStr)
			return nil, errors.KeyableError(common.NoRightsKey, common.Map{"on": onChange, "item": item, "requestedRole": rbac.RoleAdmin})
		}
	}

	// l.Infof("22222222 %s / %#v / %#v", options.Identity.ID, itemOld, itemOld.ID != options.Identity.ID)

	if itemOld.Identity.ID != options.Identity.ID && !options.Identity.Roles.Has(rbac.RoleAdmin) {
		return nil, errors.KeyableError(common.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	item.CreatedAt = itemOld.CreatedAt
	now := time.Now()
	item.UpdatedAt = &now

	path := filepath.Join(pfs.path, string(item.Identity.ID))
	if err := pfs.write(path, item); err != nil {
		return nil, errors.Wrap(err, onChange)
	}

	return &item, nil
}

const onRemove = "on personsFSStub.Remove()"

func (pfs *personsFSStub) Remove(id auth.ID, options *crud.Options) error {
	if id != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
		return errors.KeyableError(common.NoRightsKey, common.Map{"on": onRemove, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	path := filepath.Join(pfs.path, string(id)) //  pfs.path + string(id)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf(onRemove+": can't os.RemoveAll(%s), got  %s", path, err)
	}

	return nil
}

const onRead = "on personsFSStub.Read()"

func (pfs *personsFSStub) Read(id auth.ID, options *crud.Options) (*persons.Item, error) {
	if id != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.KeyableError(common.NoRightsKey, common.Map{"on": onRead, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	return pfs.read(id)
}

// read/write file ----------------------------------------------

type PersonWithCreds struct {
	persons.Item
	auth.Creds
}

func (pfs *personsFSStub) write(path string, item persons.Item) error {
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

func (pfs *personsFSStub) read(id auth.ID) (*persons.Item, error) {
	path := filepath.Join(pfs.path, string(id)) //  pfs.path + string(id)
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	var personWithCreds PersonWithCreds
	if err := json.Unmarshal(jsonBytes, &personWithCreds); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	personWithCreds.Item.Identity.ID = id
	personWithCreds.Item.SetCreds(personWithCreds.Creds)

	//for k, v := range personWithCreds.Creds {
	//	personWithCreds.SetCredsByKey(k, v)
	//}

	return &personWithCreds.Item, nil
}

const onList = "on personsFSStub.List(): "

func (pfs *personsFSStub) List(options *crud.Options) ([]persons.Item, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.KeyableError(common.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
	}

	d, err := os.Open(pfs.path)
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
		item, err := pfs.read(auth.ID(name))
		if err != nil || item == nil {
			return nil, fmt.Errorf(onList+": got %#v, %s", item, err)
		}
		// delete(item.Creds, auth.CredsPasshash)

		items = append(items, *item)
	}

	return items, nil
}
