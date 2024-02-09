package opencvlib

import (
	"regexp"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
)

const onWriteMP4 = "on videolib.WriteMP4()"

func WriteMP4(resultFilename, sourcePath string, sourceRegexp regexp.Regexp, fps float64, xWidth, yHeight int, isColor bool) error {
	fourcc, err := gocv.VideoWriterFile(resultFilename, "mp4v", fps, xWidth, yHeight, isColor)
	if err != nil {
		return errors.Wrap(err, onWriteMP4)
	}

	filenames, err := filelib.List(sourcePath, &sourceRegexp)
	if err != nil {
		return errors.Wrap(err, onWriteMP4)
	}

	for _, filename := range filenames {
		img := gocv.IMRead(filename, gocv.IMReadAnyColor)
		if err := fourcc.Write(img); err != nil {
			return errors.Wrapf(err, "on writing image from %s / "+onWriteMP4, filename)
		}
	}

	fourcc.Close()

	return nil
}
