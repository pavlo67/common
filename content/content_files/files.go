package content_files

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/selectors"
	"github.com/pavlo67/punctum/content"
)

var _ content.Operator = &contentFiles{}

type contentFiles struct {
	path      string
	marshaler basis.Marshaler
}

const onNew = "on content_file.New()"

func New(path string, marshaler basis.Marshaler) (*contentFiles, error) {
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

const onSave = "on content_files.Save()"

func (cfOp contentFiles) Save(item content.Item, options *content.SaveOptions) (id basis.ID, err error) {
	data, err := cfOp.marshaler.Marshal(item)
	if err != nil {
		return "", errors.Wrapf(err, onSave+" with native value (%#v)", item)
	}

	id = basis.ID(strconv.FormatInt(time.Now().UnixNano(), 10))
	err = ioutil.WriteFile(cfOp.path+string(id), data, 0755)
	if err != nil {
		return "", errors.Wrapf(err, onSave+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return id, nil
}

const onRead = "on content_files.Read()"

func (cfOp contentFiles) Read(id basis.ID, options *content.ReadOptions) (*content.Item, error) {
	data, err := ioutil.ReadFile(cfOp.path + string(id))
	if err != nil {
		return nil, errors.Wrapf(err, onRead+" with path (%s) & id (%s)", cfOp.path, id)
	}

	var item content.Item

	err = cfOp.marshaler.Unmarshal(data, &item)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+" with data (%s) to native value type 'content.Item'", data)
	}

	return &item, nil
}

const onList = "on content_files.List()"

func (cfOp contentFiles) List(selector *selectors.Term, options *content.ListOptions) ([]content.Brief, *content.Description, error) {
	files, err := ioutil.ReadDir(cfOp.path)
	if err != nil {
		return nil, nil, errors.Wrapf(err, onList+" with path (%s)", cfOp.path)
	}

	var briefs []content.Brief
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		item, err := cfOp.Read(basis.ID(file.Name()), nil)
		if err != nil {
			return nil, nil, errors.Wrapf(err, onList)
		}

		briefs = append(briefs, item.Brief)
	}

	return briefs, nil, nil
}

const onRemove = "on content_files.Remove()"

func (cfOp contentFiles) Remove(id basis.ID, options *content.RemoveOptions) error {
	err := os.Remove(cfOp.path + string(id))
	if err != nil {
		return errors.Wrapf(err, onRemove+" with path (%s) & id (%s)", cfOp.path, id)
	}

	return nil
}
