package files

import (
	"os"

	"io/ioutil"

	"github.com/pavlo67/workshop/common/flow"
	"github.com/pkg/errors"
)

type Origin struct {
}

type Item struct {
	Path    string        `bson:",omitempty" json:",omitempty"`
	Origins []flow.Origin `bson:",omitempty" json:",omitempty"`
}

func (f Item) FilesList() ([]os.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsDir() {
		return nil, errors.Errorf("f.Path (%s) is not a directory", f.Path)
	}

	return ioutil.ReadDir(f.Path)
}
