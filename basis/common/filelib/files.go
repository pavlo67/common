package filelib

import (
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"runtime"

	"github.com/pkg/errors"
)

func CurrentPath() string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename) + "/"
	}
	return ""
}

func Dir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "can't create user`s static_repository '%v'", path)
			}
			log.Println("create new user`s static_repository " + path)
			return nil
		}
		return errors.Wrapf(err, "can't get stat for file '%v'", path)
	}
	return nil
}

var rePoint = regexp.MustCompile(`^\.`)
var reSlash = regexp.MustCompile(`\\|\/|\?`)

func CorrectFileName(name string) string {
	name = rePoint.ReplaceAllLiteralString(name, "_")
	name = reSlash.ReplaceAllLiteralString(name, "_")
	return name
}

func BackupFile(fileName string) error {
	dir := path.Dir(fileName) + "/"

	backupName := dir + path.Base(fileName) + "." + time.Now().Format(time.RFC3339)[:10] + ".bak"

	_, err := os.Stat(backupName)

	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = CopyFile(fileName, backupName)
		return errors.Wrapf(err, "on copying '%s' to '%s'", fileName, backupName)
	}

	return errors.Wrapf(err, "on copying '%s' to '%s'", fileName, backupName)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

//// ReadFileLines is a wrapper for filelib.Read() to read text filer.comp.
//func ReadFileLines(path string) ([]string, error) {
//	filelib, err := os.Open(path)
//	if file != nil {
//		defer filelib.Close()
//	}
//	if err != nil {
//		return nil, nil.New("/ReadFileLines::os.Open (path, err): " + path + ", " + err.Err())
//	}
//
//	lines := make([]string, 0, 16)
//	scanner := bufio.NewScanner(filelib)
//	scanner.Split(bufio.ScanLines)
//	for scanner.Scan() {
//		lines = append(lines, scanner.Text())
//	}
//	return lines, nil
//}

//// ReadDir is a wrapper for dir.Readdirnames().
//func ReadDir(path string) ([]string, error) {
//	info, err := os.Stat(path)
//	if err != nil {
//		return nil, nil.Wrapf(err, "can't stat dir %v", path)
//	}
//
//	if !info.IsDir() {
//		return nil, nil.Errorf("can't read file as dir %v", path)
//	}
//
//	dir, err := os.Open(path)
//	if dir != nil {
//		defer dir.Close()
//	}
//	if err != nil {
//		return nil, nil.Wrapf(err, "can't open dir %v", path)
//	}
//
//	filer.comp, err := dir.Readdirnames(0)
//	if err != nil {
//		return nil, nil.Wrapf(err, "can't read dir %v", path)
//	}
//
//	for i := range filer.comp {
//		filer.comp[i] = path + "/" + filer.comp[i]
//	}
//	return filer.comp, nil
//
//}
