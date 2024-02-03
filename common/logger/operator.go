package logger

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "logger"

type Level int

type Config struct {
	Key         string
	LogLevel    Level
	BasePath    string
	OutputPaths []string
	ErrorPaths  []string
	Encoding    string
	SaveFiles   bool
}

const TraceLevel Level = -2
const DebugLevel Level = -1
const InfoLevel Level = 0
const WarnLevel Level = 1
const ErrorLevel Level = 2
const FatalLevel Level = 4

type Operator interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})

	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})

	Comment(text string)

	SetKey(key string)
	Key() string

	SetPath(basePath string)
	File(path string, data []byte)
	Image(path string, getImage imagelib.Imager)
}

// TODO!!! be careful in windows

var reRootPath = regexp.MustCompile(`^/`)

func ModifyPaths(paths []string, basePath string) ([]string, error) {
	if basePath = strings.TrimSpace(basePath); basePath == "" {
		return paths, nil
	}

	var err error
	if basePath, err = filelib.Dir(basePath); err != nil {
		return nil, errors.Wrapf(err, "on logger.ModifyPaths()")
	}

	modifiedPaths := make([]string, len(paths))

	for i, path := range paths {
		if path == "stdin" || path == "stdout" || path == "stderr" || reRootPath.MatchString(path) {
			modifiedPaths[i] = path
		} else {
			modifiedPaths[i] = filepath.Join(basePath, path)
		}
	}

	return modifiedPaths, nil
}
