package imagelib

import (
	"image"

	"golang.org/x/image/colornames"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"
)

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
