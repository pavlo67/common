package persons_fs_stub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/libraries/filelib"
	"github.com/pavlo67/common/common/libraries/strlib"
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

func hashCreds(creds auth.Creds) (auth.Creds, error) {
	crypt := crypt.SHA256.New()

	password := strings.TrimSpace(creds.StringDefault(auth.CredsPassword, ""))
	if password == "" {
		return nil, errors.KeyableError(errors.NoCredsKey, common.Map{"creds": creds, "reason": "no '" + auth.CredsPassword + "' key"})
	}

	hash, err := crypt.Generate([]byte(password), []byte(strlib.RandomString(10)))
	if err != nil {
		return nil, errors.Wrap(err, onAdd)
	}

	creds[auth.CredsPasshash] = hash
	delete(creds, auth.CredsPassword)

	return creds, nil
}

const onAdd = "on personsFSStub.Add()"

func (pfs *personsFSStub) Add(identity auth.Identity, data common.Map, options *crud.Options) (auth.ID, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return "", errors.KeyableError(errors.NoRightsKey, common.Map{"on": onAdd, "identity": identity, "data": data, "requestedRole": rbac.RoleAdmin})
	}

	authIDStr := strings.TrimSpace(string(identity.ID))
	if authIDStr == "" {
		return "", errors.KeyableError(errors.WrongIDKey, common.Map{"on": onAdd, "identity": identity, "data": data})
	}

	path := filepath.Join(pfs.path, authIDStr) //  pfs.path + string(authID)
	if _, err := os.Stat(path); err == nil {
		return "", errors.KeyableError(errors.DuplicateUserKey, common.Map{"on": onAdd, "identity": identity, "data": data})
	}

	var err error
	if identity.Creds, err = hashCreds(identity.Creds); err != nil {
		return "", errors.CommonError(err, onAdd)
	}

	if err := pfs.write(path, persons.Item{
		Identity:  identity,
		Data:      data,
		CreatedAt: time.Now(),
	}); err != nil {
		return "", errors.Wrap(err, onAdd)
	}

	return auth.ID(authIDStr), nil
}

func (pfs *personsFSStub) write(path string, item persons.Item) error {
	jsonBytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, jsonBytes, 0644)
}

const onChange = "on personsFSStub.Change()"

func (pfs *personsFSStub) Change(item persons.Item, options *crud.Options) error {
	if options == nil || options.Identity == nil {
		return errors.KeyableError(errors.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	itemOld, err := pfs.read(item.ID)
	if err != nil || itemOld == nil {
		errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
		if options.HasRole(rbac.RoleAdmin) {
			return errors.KeyableError(errors.WrongIDKey, common.Map{"on": onChange, "item": item, "reason": errorStr})
		} else {
			l.Error(errorStr)
			return errors.KeyableError(errors.NoRightsKey, common.Map{"on": onChange, "item": item, "requestedRole": rbac.RoleAdmin})
		}
	}

	if itemOld.ID != options.Identity.ID && !options.Identity.Roles.Has(rbac.RoleAdmin) {
		return errors.KeyableError(errors.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	if newPassword := strings.TrimSpace(item.Creds.StringDefault(auth.CredsPassword, "")); newPassword != "" {
		var err error
		if item.Creds, err = hashCreds(item.Creds); err != nil {
			return errors.CommonError(err, onChange)
		}
	} else if item.Creds != nil {
		item.Creds[auth.CredsPasshash] = itemOld.Creds[auth.CredsPasshash]
	} else {
		item.Creds = common.Map{auth.CredsPasshash: itemOld.Creds[auth.CredsPasshash]}
	}
	item.CreatedAt = itemOld.CreatedAt
	now := time.Now()
	item.UpdatedAt = &now

	path := filepath.Join(pfs.path, string(item.ID))
	if err := pfs.write(path, item); err != nil {
		return errors.Wrap(err, onChange)
	}

	return nil
}

const onList = "on personsFSStub.List(): "

func (pfs *personsFSStub) List(options *crud.Options) ([]persons.Item, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.KeyableError(errors.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
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
			return nil, errors.Errorf(onList+": got %#v, %s", item, err)
		}
		delete(item.Creds, auth.CredsPasshash)

		items = append(items, *item)
	}

	return items, nil
}

const onRemove = "on personsFSStub.Remove()"

func (pfs *personsFSStub) Remove(authID auth.ID, options *crud.Options) error {
	if !options.HasRole(rbac.RoleAdmin) {
		return errors.KeyableError(errors.NoRightsKey, common.Map{"on": onRemove, "authID": authID, "requestedRole": rbac.RoleAdmin})
	}

	path := filepath.Join(pfs.path, string(authID)) //  pfs.path + string(authID)
	if err := os.RemoveAll(path); err != nil {
		return errors.Errorf(onRemove+": can't os.RemoveAll(%s), got  %s", path, err)
	}

	return nil
}

const onRead = "on personsFSStub.Read()"

func (pfs *personsFSStub) Read(authID auth.ID, options *crud.Options) (*persons.Item, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.KeyableError(errors.NoRightsKey, common.Map{"on": onRead, "authID": authID, "requestedRole": rbac.RoleAdmin})
	}

	return pfs.read(authID)
}

func (pfs *personsFSStub) read(authID auth.ID) (*persons.Item, error) {
	path := filepath.Join(pfs.path, string(authID)) //  pfs.path + string(authID)
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}
	var item persons.Item
	if err := json.Unmarshal(jsonBytes, &item); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &item, nil
}
