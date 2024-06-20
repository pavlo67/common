package imagelib

import (
	"fmt"
	pnm "github.com/jbuchbinder/gopnm"
	"github.com/pavlo67/common/common/filelib"
	"image"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type PGMConfig struct {
	Width  int
	Height int
	MaxVal int
}

var reSpaces = regexp.MustCompile(`\s+`)
var reSpace = regexp.MustCompile(`\s`)

func DecodeConfigPGM(r io.ReadSeeker) (c PGMConfig, err error) {
	header := make([]byte, 20)

	n, err := r.Read(header)
	if err != nil {
		return c, fmt.Errorf("reading header failed: %s", err)
	} else if n != len(header) {
		return c, fmt.Errorf("read incomplete header: %d of %d bytes", n, len(header))
	}

	idx := reSpaces.FindIndex(header)
	if len(idx) != 2 {
		return c, fmt.Errorf("wrong header (1st field is absent): %s", header)
	} else if string(header[:idx[0]]) != "P5" {
		return c, fmt.Errorf("wrong header (1st field must be 'P5'): %s", header)
	}

	from := idx[1]
	idx = reSpaces.FindIndex(header[from:])
	if len(idx) != 2 {
		return c, fmt.Errorf("wrong header (2nd field is absent): %s", header)
	} else if c.Width, err = strconv.Atoi(string(header[from : from+idx[0]])); err != nil {
		return c, fmt.Errorf("wrong header (2nd field must be integer): %s", header)
	} else if c.Width <= 0 {
		return c, fmt.Errorf("wrong header (2nd field must be > 0): %s", header)
	}

	from += idx[1]
	idx = reSpaces.FindIndex(header[from:])
	if len(idx) != 2 {
		return c, fmt.Errorf("wrong header (3th field is absent): %s", header)
	} else if c.Height, err = strconv.Atoi(string(header[from : from+idx[0]])); err != nil {
		return c, fmt.Errorf("wrong header (3th field must be integer): %s", header)
	} else if c.Height <= 0 {
		return c, fmt.Errorf("wrong header (3th field must be > 0): %s", header)
	}

	from += idx[1]
	// TODO!!! be careful: only one space can be read here
	idx = reSpace.FindIndex(header[from:])
	if len(idx) != 2 {
		return c, fmt.Errorf("wrong header (3th field is absent): %s", header)
	} else if c.MaxVal, err = strconv.Atoi(string(header[from : from+idx[0]])); err != nil {
		return c, fmt.Errorf("wrong header (3th field must be integer): %s", header)
	} else if c.MaxVal < 0 || c.MaxVal > 255 {
		return c, fmt.Errorf("wrong header (3th field must be in the range [0..255]): %s", header)
	}

	from += idx[1]
	n1, err := r.Seek(int64(from), io.SeekStart)
	if err != nil {
		return c, fmt.Errorf("seeking at body failed: %s", err)
	} else if n1 != int64(from) {
		return c, fmt.Errorf("seeking at body failed: %d instead of %d", n1, from)
	}

	return c, nil
}

func Decode(r io.ReadSeeker) (image.Image, error) {
	c, err := DecodeConfigPGM(r)
	if err != nil {
		return nil, fmt.Errorf("pgm: parsing header failed: %s", err)
	}

	m := image.NewGray(image.Rect(0, 0, c.Width, c.Height))

	n, err := r.Read(m.Pix)
	if err != nil {
		return nil, fmt.Errorf("pgm: reading image failed: %s", err)
	} else if n != len(m.Pix) {
		return nil, fmt.Errorf("pgm: read incomplete image: %d of %d bytes", n, len(m.Pix))
	}

	return m, nil
}

const onReadPGMSpecial = "on imagelib.ReadPGMSpecial()"

func ReadPGMSpecial(filename string) (*image.Gray, error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, onReadPGMSpecial)
	}
	defer srcFile.Close()

	img, err := Decode(srcFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %s / "+onReadPGMSpecial, filename, err)
	}

	imgGray, _ := img.(*image.Gray)
	if imgGray == nil {
		return nil, fmt.Errorf("%s: imgGray == nil / "+onReadPGMSpecial, filename)
	}

	return imgGray, nil
}

const onSavePGM = "on imagelib.SavePGM()"

func SavePGM(img image.Image, filename string) error {
	if img == nil {
		return errors.New("img == nil / " + onSavePGM)
	} else if path := filepath.Dir(filename); path != "" && path != "." && path != ".." {
		if _, err := filelib.Dir(path); err != nil {
			return errors.Wrapf(err, "can't create dir '%s' / "+onSavePGM, path)
		}
	}

	resFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, onSavePGM)
	}
	defer resFile.Close()

	if err = pnm.Encode(resFile, img, pnm.PGM); err != nil {
		return errors.Wrap(err, onSavePGM)
	}
	return nil
}
