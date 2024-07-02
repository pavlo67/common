package imagelib

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/filelib"
)

const onSavePNG = "on imagelib.SavePNG()"

func SavePNG(img image.Image, filename string) error {

	if img == nil {
		return errors.New("img == nil / " + onSavePNG)
	}

	if path := filepath.Dir(filename); path != "" && path != "." && path != ".." {
		if _, err := filelib.Dir(path); err != nil {
			return errors.Wrapf(err, "can't create dir '%s' / "+onSavePNG, path)
		}
	}

	resFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, onSavePNG)
	}
	defer resFile.Close()

	if err = png.Encode(resFile, img); err != nil {
		return errors.Wrap(err, onSavePNG)
	}
	return nil
}
