package imagelib

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/filelib"
)

const onSave = "on imagelib.Save()"

func Save(img image.Image, filename string) error {

	if img == nil {
		return errors.New("img == nil / " + onSave)
	}

	if path := filepath.Dir(filename); path != "" && path != "." && path != ".." {
		if _, err := filelib.Dir(path); err != nil {
			return errors.Wrapf(err, "can't create dir '%s' / "+onSave, path)
		}
	}

	resFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, onSave)
	}
	defer resFile.Close()

	if err = png.Encode(resFile, img); err != nil {
		return errors.Wrap(err, onSave)
	}
	return nil
}
