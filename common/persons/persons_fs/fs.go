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
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
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
		return nil, nil, errata.CommonError(err, onNew)
	}

	personsOp := personsFSStub{path: path}

	return &personsOp, &personsOp, nil
}

const onAdd = "on personsFSStub.Add()"

func (pfs *personsFSStub) Add(identity auth.Identity, data common.Map, options *crud.Options) (auth.ID, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return "", errata.KeyableError(errata.NoRightsKey, common.Map{"on": onAdd, "identity": identity, "data": data, "requestedRole": rbac.RoleAdmin})
	}

	authIDStr := strings.TrimSpace(string(identity.ID))
	if authIDStr == "" {
		authIDStr = strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(rand.Int())
	}

	path := filepath.Join(pfs.path, authIDStr) //  pfs.path + string(authID)
	if _, err := os.Stat(path); err == nil {
		return "", errata.KeyableError(errata.DuplicateUserKey, common.Map{"on": onAdd, "identity": identity, "data": data})
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

func (pfs *personsFSStub) Change(item persons.Item, options *crud.Options) (*persons.Item, error) {
	if options == nil || options.Identity == nil {
		return nil, errata.KeyableError(errata.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	// l.Info(1111111111, " ", item.ID)

	itemOld, err := pfs.read(item.ID)
	if err != nil || itemOld == nil {
		errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
		if options.HasRole(rbac.RoleAdmin) {
			return nil, errata.KeyableError(errata.WrongIDKey, common.Map{"on": onChange, "item": item, "reason": errorStr})
		} else {
			l.Error(errorStr)
			return nil, errata.KeyableError(errata.NoRightsKey, common.Map{"on": onChange, "item": item, "requestedRole": rbac.RoleAdmin})
		}
	}

	// l.Infof("22222222 %s / %#v / %#v", options.Identity.ID, itemOld, itemOld.ID != options.Identity.ID)

	if itemOld.ID != options.Identity.ID && !options.Identity.Roles.Has(rbac.RoleAdmin) {
		return nil, errata.KeyableError(errata.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	item.CreatedAt = itemOld.CreatedAt
	now := time.Now()
	item.UpdatedAt = &now

	path := filepath.Join(pfs.path, string(item.ID))
	if err := pfs.write(path, item); err != nil {
		return nil, errors.Wrap(err, onChange)
	}

	return &item, nil
}

const onList = "on personsFSStub.List(): "

func (pfs *personsFSStub) List(options *crud.Options) ([]persons.Item, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return nil, errata.KeyableError(errata.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
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

const onRemove = "on personsFSStub.Remove()"

func (pfs *personsFSStub) Remove(authID auth.ID, options *crud.Options) error {
	if authID != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
		return errata.KeyableError(errata.NoRightsKey, common.Map{"on": onRemove, "authID": authID, "requestedRole": rbac.RoleAdmin})
	}

	path := filepath.Join(pfs.path, string(authID)) //  pfs.path + string(authID)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf(onRemove+": can't os.RemoveAll(%s), got  %s", path, err)
	}

	return nil
}

const onRead = "on personsFSStub.Read()"

func (pfs *personsFSStub) Read(authID auth.ID, options *crud.Options) (*persons.Item, error) {
	if authID != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
		return nil, errata.KeyableError(errata.NoRightsKey, common.Map{"on": onRead, "authID": authID, "requestedRole": rbac.RoleAdmin})
	}

	return pfs.read(authID)
}

func (pfs *personsFSStub) read(authID auth.ID) (*persons.Item, error) {
	// l.Info(10000000, " ", authID)

	path := filepath.Join(pfs.path, string(authID)) //  pfs.path + string(authID)
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}
	var item persons.Item
	if err := json.Unmarshal(jsonBytes, &item); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	// l.Infof("readed: %#v", item)

	item.ID = authID

	return &item, nil
}
