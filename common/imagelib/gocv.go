package imagelib

import (
	"fmt"
	"image"
	"math"

	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/errors"
)

const onContourToGrayscale = "on imagelib.ContourToGrayscale()"

type ContourImage struct {
	Contour gocv.PointVector
	image.Rectangle
}

var _ GetImage = &ContourImage{}

func (imageOp *ContourImage) Bounds() image.Rectangle {
	if imageOp == nil {
		return image.Rectangle{}
	}

	return imageOp.Rectangle
}

func (imageOp *ContourImage) Image() (image.Image, string, error) {
	if imageOp == nil {
		return nil, "", errors.New("*ContourImage = nil")
	}

	return ContourToGrayscale(imageOp.Contour, imageOp.Rectangle)
}

func ContourToGrayscale(contour gocv.PointVector, rect image.Rectangle) (image.Image, string, error) {
	mat := gocv.NewMatWithSize(rect.Max.Y-rect.Min.Y, rect.Max.X-rect.Min.X, gocv.MatTypeCV8UC1)
	defer mat.Close()

	contours := gocv.NewPointsVector()
	defer contours.Close()

	contours.Append(contour)
	gocv.DrawContours(&mat, contours, 0, colornames.White, 1)

	img, err := mat.ToImage()
	if err != nil {
		return nil, "", errors.Wrap(err, onContourToGrayscale)
	}

	return img, "", nil
}

func ContourToGrayscalePng(contour gocv.PointVector, rect image.Rectangle, path string) error {
	img, _, err := ContourToGrayscale(contour, rect)
	if err != nil {
		return err
	}

	return SavePNG(img, path)
}

func ContourAreaPix(contour gocv.PointVector) (float64, float64) {
	contourArea := gocv.ContourArea(contour)
	return contourArea, math.Sqrt(4 * contourArea / math.Pi)
}

const onResize = "on imagelib.Resize()"

func Resize(rgb image.RGBA, scale float64) (*image.RGBA, float64, error) {
	if scale == 1 || scale == 0 {
		return &rgb, 1, nil
	} else if scale < 0 {
		return nil, 0, fmt.Errorf("wrong resize scale (%f) / "+onResize, scale)
	}

	mat, err := gocv.ImageToMatRGB(&rgb)
	if err != nil {
		return nil, 0, errors.Wrap(err, onResize)
	}
	defer mat.Close()

	matForResize := gocv.NewMat()
	defer matForResize.Close()

	gocv.Resize(mat, &matForResize, image.Point{}, scale, scale, gocv.InterpolationDefault)

	imgResized, err := matForResize.ToImage()
	if err != nil {
		return nil, 0, errors.Wrap(err, onResize)
	}

	rgbaResized, ok := imgResized.(*image.RGBA)
	if !ok {
		return nil, 0, fmt.Errorf("resized image has wrong type: %T / "+onResize, rgbaResized)
	}

	return rgbaResized, scale, nil
}

const onTranspose = "on imagelib.Transpose()"

func Transpose(rgb image.RGBA) (*image.RGBA, error) {

	mat, err := gocv.ImageToMatRGB(&rgb)
	if err != nil {
		return nil, errors.Wrap(err, onTranspose)
	}
	defer mat.Close()

	matForTranspose := gocv.NewMat()
	defer matForTranspose.Close()

	gocv.Transpose(mat, &matForTranspose)

	imgTransposed, err := matForTranspose.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onTranspose)
	}

	rgbaTransposed, ok := imgTransposed.(*image.RGBA)
	if !ok {
		return nil, fmt.Errorf("transposed image has wrong type: %T / "+onTranspose, rgbaTransposed)
	}

	return rgbaTransposed, nil
}
