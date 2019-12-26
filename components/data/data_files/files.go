package data_files

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/lingua/structura"
	"github.com/pavlo67/workshop/common"

	"github.com/pavlo67/workshop/common/selectors"
)

var _ structura.Operator = &contentFiles{}

type contentFiles struct {
	path      string
	marshaler libs.Marshaler
}

const onNew = "on contentFiles.New()"

func New(path string, marshaler libs.Marshaler) (*contentFiles, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New(onNew + ": no path")
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}

	if marshaler == nil {
		return nil, errors.New(onNew + ": no basis.Marshaler")
	}

	return &contentFiles{
		path:      path,
		marshaler: marshaler,
	}, nil
}

const onDescript = "on contentFiles.Descript()"

func (cfOp contentFiles) Descript() (*structura.Description, error) {
	return nil, nil
}

const onSave = "on contentFiles.Save()"

func (cfOp contentFiles) Save(item structura.Item, options *structura.SaveOptions) (id common.Key, err error) {
	data, err := cfOp.marshaler.Marshal(item)
	if err != nil {
		return "", errors.Wrapf(err, onSave+" with native value (%#v)", item)
	}

	id = common.Key(strconv.FormatInt(time.Now().UnixNano(), 10))
	err = ioutil.WriteFile(cfOp.path+string(id), data, 0755)
	if err != nil {
		return "", errors.Wrapf(err, onSave+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return id, nil
}

const onRead = "on contentFiles.Read()"

func (cfOp contentFiles) Read(id common.Key, options *structura.GetOptions) (*structura.Item, error) {
	data, err := ioutil.ReadFile(cfOp.path + string(id))
	if err != nil {
		return nil, errors.Wrapf(err, onRead+" with path (%s) & id (%s)", cfOp.path, id)
	}

	var item structura.Item

	err = cfOp.marshaler.Unmarshal(data, &item)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+" with data (%s) to native value type 'content.Item'", data)
	}

	return &item, nil
}

const onList = "on contentFiles.ListTags()"

func (cfOp contentFiles) List(selector *selectors.Term, options *structura.GetOptions) ([]structura.Brief, error) {
	files, err := ioutil.ReadDir(cfOp.path)
	if err != nil {
		return nil, errors.Wrapf(err, onList+" with path (%s)", cfOp.path)
	}

	var briefs []structura.Brief
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		item, err := cfOp.Read(common.Key(file.Name()), nil)
		if err != nil {
			return nil, errors.Wrapf(err, onList)
		}

		briefs = append(briefs, item.Brief)
	}

	return briefs, nil
}

const onRemove = "on contentFiles.Remove()"

func (cfOp contentFiles) Remove(id common.Key, options *structura.RemoveOptions) error {
	err := os.Remove(cfOp.path + string(id))
	if err != nil {
		return errors.Wrapf(err, onRemove+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return nil
}
