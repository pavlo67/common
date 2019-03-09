package crud_file

import (
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
)

var _ crud.Operator = &crudFile{}

type crudFile struct {
	crud.Mapper

	path      string
	marshaler basis.Marshaler
}

const onNew = "on crud_file.New()"

func New(mapper crud.Mapper, path string, marshaler basis.Marshaler) (*crudFile, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New(onNew + ": no path")
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}

	if mapper == nil {
		return nil, errors.New(onNew + ": no crud.Mapper")
	}

	if marshaler == nil {
		return nil, errors.New(onNew + ": no basis.Marshaler")
	}

	return &crudFile{
		Mapper:    mapper,
		path:      path,
		marshaler: marshaler,
	}, nil
}

const onCreate = "on crud_file.Create()"

func (cfOp crudFile) Create(_ auth.ID, native interface{}) (id string, err error) {
	data, err := cfOp.marshaler.Marshal(native)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+" with native value (%#v)", native)
	}

	id = strconv.FormatInt(time.Now().UnixNano(), 10)
	err = ioutil.WriteFile(cfOp.path+id, data, 0755)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return id, nil
}

const onRead = "on crud_file.Read()"

func (cfOp crudFile) Read(_ auth.ID, id string) (interface{}, error) {
	data, err := ioutil.ReadFile(cfOp.path + id)
	if err != nil {
		return "", errors.Wrapf(err, onRead+" with path (%s) & id (%s)", cfOp.path, id)
	}

	description := cfOp.Description()
	nativeValue := reflect.New(reflect.ValueOf(description.ExemplarNative).Elem().Type()).Interface()

	err = cfOp.marshaler.Unmarshal(data, nativeValue)
	if err != nil {
		return "", errors.Wrapf(err, onRead+" with data (%s) to native value type (%T)", data, nativeValue)
	}

	return nativeValue, nil
}

const onReadList = "on crud_file.ReadList()"

func (cfOp crudFile) ReadList(userIS auth.ID, options crud.ReadOptions) ([]interface{}, *uint64, error) {
	files, err := ioutil.ReadDir(cfOp.path)
	if err != nil {
		return nil, nil, errors.Wrapf(err, onReadList+" with path (%s)", cfOp.path)
	}

	var nativeValues []interface{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		nativeValue, err := cfOp.Read("", file.Name())
		if err != nil {
			return nil, nil, errors.Wrapf(err, onReadList)
		}

		nativeValues = append(nativeValues, nativeValue)
	}

	num := uint64(len(nativeValues))
	return nativeValues, &num, nil
}

const onUpdate = "on crud_file.Update()"

func (cfOp crudFile) Update(_ auth.ID, id string, native interface{}) error {
	data, err := cfOp.marshaler.Marshal(native)
	if err != nil {
		return errors.Wrapf(err, onUpdate+" with native value (%#v)", native)
	}

	err = ioutil.WriteFile(cfOp.path+id, data, 0755)
	if err != nil {
		return errors.Wrapf(err, onUpdate+" with path (%s), id (%s) & data(%s)", cfOp.path, id, data)
	}

	return nil

}

const onDelete = "on crud_file.Update()"

func (cfOp crudFile) Delete(_ auth.ID, id string) error {
	err := os.Remove(cfOp.path + id)
	if err != nil {
		return errors.Wrapf(err, onDelete+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return nil
}
