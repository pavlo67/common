package imagelib

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
)

const onSavePNG = "on SavePNG()"

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

type PointsImage struct {
	Points []image.Point
	image.Rectangle
}

var _ Imager = &PointsImage{}

func (imageOp *PointsImage) Bounds() image.Rectangle {
	if imageOp == nil {
		return image.Rectangle{}
	}

	return imageOp.Rectangle
}

func (imageOp *PointsImage) Image() (image.Image, string, error) {
	if imageOp == nil {
		return nil, "", errors.New("*PointsImage = nil")
	}

	return PointsToGrayscale(imageOp.Points, imageOp.Rectangle), "", nil
}

// PointsToGrayscale returns Gray_ structure instead of image.Gray because the structure implements Imager interface required for show.Demo()
func PointsToGrayscale(points []image.Point, rect image.Rectangle) image.Image {
	xWidth := rect.Max.X - rect.Min.X
	yHeight := rect.Max.Y - rect.Min.Y

	gray := image.Gray{
		Pix:    make([]uint8, xWidth*yHeight),
		Stride: xWidth,
		Rect:   rect,
	}

	for _, p := range points {
		gray.Set(p.X, p.Y, color.White)
	}

	return &gray
}

func PointsToGrayscalePng(points []image.Point, rect image.Rectangle, path string) error {
	img := PointsToGrayscale(points, rect)
	return SavePNG(img, path)
}

const onImageGray = "on ImageGray()"

func ImageGray(srcFilename, resFilename string, rect *image.Rectangle) (imgArea *image.Gray, original *image.Rectangle, err error) {
	srcFile, err := os.Open(srcFilename)
	if err != nil {
		// return nil, nil, errors.Wrapf(err, "on opening %s / "+onImageGray, srcFilename)
		return nil, nil, errors.Wrap(err, onImageGray)
	}
	defer srcFile.Close()

	src, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "on decoding %s / "+onImageGray, srcFilename)
	}

	switch img := src.(type) {
	//case *img.Gray16:
	//	t.Logf("%T", img)
	case *image.Gray:
		if rect == nil {
			return img, &img.Rect, nil
		}

		imgArea := img.SubImage(*rect).(*image.Gray)

		if resFilename != "" {
			if err = SavePNG(imgArea, resFilename); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}

		return imgArea, &img.Rect, nil
	}

	return nil, nil, fmt.Errorf("wrong format (%T) of %s / "+onImageGray, src, srcFilename)
}

const onImageRGBA = "on ImageRGBA()"

func ImageRGBA(filename string, rect *image.Rectangle) (imgRGBA *image.RGBA, original *image.Rectangle, err error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		// return nil, nil, errors.Wrapf(err, "on opening %s / "+onImageRGBA, filename)
		return nil, nil, errors.Wrap(err, onImageRGBA)
	}
	defer srcFile.Close()

	src, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "on decoding %s / "+onImageRGBA, filename)
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
		return nil, nil, fmt.Errorf("wrong format (%T) of %s / "+onImageRGBA, src, filename)

	}

	if rect != nil {
		imgRGBA = imgRGBA.SubImage(*rect).(*image.RGBA)
	}

	//if resFilename != "" {
	//	if err = SavePNG(imgRGBA, resFilename); err != nil {
	//		fmt.Fprint(os.Stderr, err)
	//	}
	//}

	return imgRGBA, &imgRGBA.Rect, nil

}
