package imagelib

import (
	"fmt"
	"regexp"

	"github.com/pavlo67/common/common/errors"

	// pnm "github.com/jbuchbinder/gopnm"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var rePGM = regexp.MustCompile(`(?i)\.pgm$`)

const onRead = "on imagelib.Read()"

func Read(filename string) (image.Image, error) {

	if rePGM.MatchString(filename) {
		return ReadPGMSpecial(filename)
	}

	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, errors.Wrapf(err, "on decoding %s / "+onRead, filename)
	}

	return img, nil
}

const onReadGray = "on imagelib.ReadGray()"

func ReadGray(filename string) (*image.Gray, error) {
	if rePGM.MatchString(filename) {
		return ReadPGMSpecial(filename)
	}

	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, onReadGray)
	}
	defer srcFile.Close()

	src, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, errors.Wrapf(err, "on decoding %s / "+onReadGray, filename)
	}

	imgGray, err := ImageToGray(src)
	if err != nil {
		return nil, errors.Wrap(err, onReadGray)
	}

	return imgGray, nil
}

const onReadRGBA = "on imagelib.ReadRGBA()"

func ReadRGBA(filename string) (imgRGBA *image.RGBA, _ error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, onReadRGBA)
	}
	defer srcFile.Close()

	src, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, errors.Wrapf(err, "on decoding %s / "+onReadRGBA, filename)
	}

	switch img := src.(type) {
	case *image.RGBA:
		imgRGBA = img

	case *image.NRGBA:
		imgRGBA = &image.RGBA{
			Pix:    img.Pix,
			Stride: img.Stride,
			Rect:   img.Rect,
		}

	case *image.Gray:
		imgRGBA = image.NewRGBA(img.Rect)
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				imgRGBA.Set(x, y, img.At(x, y))
			}
		}

	case *image.YCbCr:
		imgRGBA = image.NewRGBA(img.Rect)
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				imgRGBA.Set(x, y, img.At(x, y))
			}
		}

	case *image.CMYK:
		imgRGBA = image.NewRGBA(img.Rect)
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				imgRGBA.Set(x, y, img.At(x, y))
			}
		}

	default:
		return nil, fmt.Errorf("wrong format (%T) of %s / "+onReadRGBA, src, filename)

	}

	return imgRGBA, nil
}
