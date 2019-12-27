package files

import (
	"io/ioutil"
	"os"

	"github.com/pavlo67/workshop/common/flow"
)

type Item struct {
	Path   string      `bson:",omitempty" json:",omitempty"`
	Origin flow.Origin `bson:",omitempty" json:",omitempty"`
}

func (f Item) IsDir() (bool, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}

func (f Item) FilesList() ([]os.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return nil, err
	}

	if fi.Mode().IsDir() {
		return ioutil.ReadDir(f.Path)
	}

	return []os.FileInfo{fi}, nil
}
