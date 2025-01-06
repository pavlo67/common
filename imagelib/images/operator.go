package images

import (
	"image"
	"path/filepath"
	"regexp"

	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/imagelib/sources"
)

const InterfaceKey joiner.InterfaceKey = "images"
const InterfaceCleanerKey joiner.InterfaceKey = "images_cleaner"

type Key string

type Operator interface {
	Get(Key) (image.Image, *sources.Description, error)
	Save(image.Image, sources.Description, Key) (filepath string, err error)
	Check(Key) (bool, error)
	ListPaths(keyRegexStr string) ([]string, error)
}

func KeyPath(key Key, colored bool) string {
	path := filepath.Clean(string(key))
	if path == "" {
		return ""
	}

	if colored {
		return path + ".png"
	}
	return path + ".pgm"
}

func KeyPathRegex(keyRegexStr string, colored bool) *regexp.Regexp {
	if keyRegexStr = filepath.Clean(keyRegexStr); keyRegexStr == "" {
		return nil
	}

	if colored {
		return regexp.MustCompile(keyRegexStr + `\.png$`)
	}
	return regexp.MustCompile(keyRegexStr + `\.pgm$`)
}
