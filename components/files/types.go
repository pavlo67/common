package files

import (
	"io/ioutil"
	"os"
)

type Item struct {
	Path string
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
