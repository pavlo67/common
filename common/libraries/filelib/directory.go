package filelib

import (
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"
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

func GetDir(path, pathDefault string) (string, error) {

	path = strings.TrimSpace(path)
	if path == "" {
		path = pathDefault

	}

	// converting Windows-backslashed pathes to the normal ones
	path = reBackslash.ReplaceAllString(path, "/")

	fi, err := os.Stat(path)
	if err != nil {
		return "", errors.Wrapf(err, onGetDir+"can't os.Stat(%s)", path)
	}

	if !fi.IsDir() {
		return "", errors.Wrapf(err, onGetDir+"path (%s) isn't a directory", path)
	}

	if path[len(path)-1] != '/' {
		path += "/"
	}

	return path, nil
}

func Dir(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("can't create dir for empty path")
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "can't create dir '%s'", path)
			}
			return nil
		}
		return errors.Wrapf(err, "can't get stat for file '%s'", path)
	}
	return nil
}
