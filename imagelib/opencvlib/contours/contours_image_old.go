package contours

import (
	"github.com/pavlo67/common/common"
	"image"
	"math"

	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/logger"
)

type ContourImage struct {
	Contour gocv.PointVector
	image.Rectangle
}

var _ logger.GetImage = &ContourImage{}

func (imageOp *ContourImage) Bounds() image.Rectangle {
	if imageOp == nil {
		return image.Rectangle{}
	}

	return imageOp.Rectangle
}

func (imageOp *ContourImage) Image(opts common.Map) (image.Image, string, error) {
	if imageOp == nil {
		return nil, "", errors.New("*ContourImage = nil")
	}

	return ContourToGrayscale(imageOp.Contour, imageOp.Rectangle)
}

const onContourToGrayscale = "on ContourToGrayscale()"

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

	return imagelib.SavePNG(img, path)
}

func ContourAreaPix(contour gocv.PointVector) (float64, float64) {
	contourArea := gocv.ContourArea(contour)
	return contourArea, math.Sqrt(4 * contourArea / math.Pi)
}
