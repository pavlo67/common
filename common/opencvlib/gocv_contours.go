package opencvlib

import (
	"image"
	"math"

	"github.com/pavlo67/common/common/imagelib"

	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/errors"
)

const onContourToGrayscale = "on ContourToGrayscale()"

type ContourImage struct {
	Contour gocv.PointVector
	image.Rectangle
}

var _ imagelib.Imager = &ContourImage{}

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

	return imagelib.SavePNG(img, path)
}

func ContourAreaPix(contour gocv.PointVector) (float64, float64) {
	contourArea := gocv.ContourArea(contour)
	return contourArea, math.Sqrt(4 * contourArea / math.Pi)
}

const onFillOutsideContours = "on imagelib.GrayOutsideContours()"

func GrayWhitedOutsideContours(imgGray image.Gray, psv gocv.PointsVector) (*image.Gray, error) {

	matImg, err := gocv.ImageGrayToMatGray(&imgGray)
	if err != nil {
		return nil, errors.Wrap(err, onFillOutsideContours)
	}

	matMaskCntrs := gocv.NewMatWithSize(imgGray.Rect.Dy(), imgGray.Rect.Dx(), gocv.MatTypeCV8U)
	gocv.FillPoly(&matMaskCntrs, psv, colornames.White)

	matMaskCntrsOutside := gocv.NewMat()
	gocv.BitwiseNot(matMaskCntrs, &matMaskCntrsOutside)

	matWhitedOutside := gocv.NewMat()
	gocv.BitwiseOr(matImg, matMaskCntrsOutside, &matWhitedOutside)

	imgWhitedOutside, err := matWhitedOutside.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onFillOutsideContours)
	}

	imgGrayWhitedOutside, _ := imgWhitedOutside.(*image.Gray)
	if imgGrayWhitedOutside == nil {
		return nil, errors.New("imgGrayWhitedOutside == nil / " + onFillOutsideContours)
	}

	return imgGrayWhitedOutside, nil
}

func GrayBlackOutsideContours(imgGray image.Gray, psv gocv.PointsVector) (*image.Gray, error) {

	matImg, err := gocv.ImageGrayToMatGray(&imgGray)
	if err != nil {
		return nil, errors.Wrap(err, onFillOutsideContours)
	}

	matMaskCntrs := gocv.NewMatWithSize(imgGray.Rect.Dy(), imgGray.Rect.Dx(), gocv.MatTypeCV8U)
	gocv.FillPoly(&matMaskCntrs, psv, colornames.White)

	matImgMasked := gocv.NewMat()
	matImg.CopyToWithMask(&matImgMasked, matMaskCntrs)

	imgBlackedOutside, err := matImgMasked.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onFillOutsideContours)
	}

	imgGrayBlackedOutside, _ := imgBlackedOutside.(*image.Gray)
	if imgGrayBlackedOutside == nil {
		return nil, errors.New("imgGrayBlackedOutside == nil / " + onFillOutsideContours)
	}

	return imgGrayBlackedOutside, nil
}
