package filelib

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"strconv"
	"time"

	"github.com/pavlo67/common/common/errata"
)

func RelativePath(pathFull, pathBase, pathPrefix string) string {
	path := ""

	if pathFull[:len(pathBase)] == pathBase {
		path = pathFull[len(pathBase):]
	} else {
		log.Printf(".RelativePath(%s, %s, %s)????", pathFull, pathBase, pathPrefix)
	}

	return pathPrefix + path
}

func CurrentPath() string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename) + "/"
	}
	return ""
}

const onGetDir = "on filelib.GetDir(): "

func GetDir(path string) (string, error) {

	// converting Windows-backslashed pathes to the normal ones
	path = reBackslash.ReplaceAllString(path, "/")

	fi, err := os.Stat(path)
	if err != nil {
		return "", errata.Wrapf(err, onGetDir+"can't os.Stat(%s)", path)
	}

	if !fi.IsDir() {
		return "", errata.Wrapf(err, onGetDir+"path (%s) isn't a directory", path)
	}

	if path[len(path)-1] != '/' {
		path += "/"
	}

	return path, nil
}

func Dir(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errata.New("can't create dir for empty path")
	}

	// converting Windows-backslashed pathes to the normal ones
	path = reBackslash.ReplaceAllString(path, "/")
	if path[len(path)-1] != '/' {
		path += "/"
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return "", errata.Wrapf(err, "can't create dir '%s'", path)
			}
			return path, nil
		}
		return "", errata.Wrapf(err, "can't get stat for dir '%s'", path)
	}

	return path, nil
}

func ClearDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err = os.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}
	return nil
}

const maxRetries = 10

func SubDirUnique(path string) (string, error) {
	path, err := Dir(path)
	if err != nil {
		return "", err
	}

	var subpath string

	for i := 0; i < maxRetries; i++ {
		subpath, err = Dir(path + CorrectFileName(time.Now().Format(time.RFC3339)) + "_" + strconv.Itoa(i))
		if err == nil {
			return subpath, nil
		}
	}

	return "", fmt.Errorf("can't create unique subpath %d times, last try was '%s'", maxRetries, subpath)
}
