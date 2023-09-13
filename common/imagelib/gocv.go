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

var _ Imager = &ContourImage{}

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

const onResizeToRange = "on imagelib.ResizeToRange()"

func ResizeToRange(imgRGB image.RGBA, dpm float64, dpmRange [2]float64) (*image.RGBA, float64, error) {
	if !(dpm > 0 && !math.IsInf(dpm, 1)) {
		return nil, 0, fmt.Errorf("wrong dpm: %f / "+onResizeToRange, dpm)
	}
	if dpm >= dpmRange[0] && dpm <= dpmRange[1] {
		return &imgRGB, dpm, nil
	}

	imgRGBResized, resizeRatio, err := Resize(imgRGB, 0.5*(dpmRange[0]+dpmRange[1])/dpm)
	if err != nil {
		return nil, 0, errors.Wrap(err, onResizeToRange)
	} else if imgRGBResized == nil {
		return nil, 0, errors.New("resized img == nil / " + onResizeToRange)
	}

	return imgRGBResized, dpm * resizeRatio, nil
}

const onResize = "on imagelib.Resize()"

func Resize(imgRGB image.RGBA, ratio float64) (*image.RGBA, float64, error) {
	if ratio == 1 || ratio == 0 {
		return &imgRGB, 1, nil
	} else if ratio < 0 || math.IsNaN(ratio) || math.IsInf(ratio, 0) {
		return nil, 0, fmt.Errorf("wrong resize ratio (%f) / "+onResize, ratio)
	}

	mat, err := gocv.ImageToMatRGB(&imgRGB)
	if err != nil {
		return nil, 0, errors.Wrap(err, onResize)
	}
	defer mat.Close()

	matForResize := gocv.NewMat()
	defer matForResize.Close()

	gocv.Resize(mat, &matForResize, image.Point{}, ratio, ratio, gocv.InterpolationDefault)

	imgResized, err := matForResize.ToImage()
	if err != nil {
		return nil, 0, errors.Wrap(err, onResize)
	}

	rgbaResized, ok := imgResized.(*image.RGBA)
	if !ok {
		return nil, 0, fmt.Errorf("resized image has wrong type: %T / "+onResize, rgbaResized)
	}

	return rgbaResized, ratio, nil
}

const onRotateResized = "on imagelib.RotateResized()"

func RotateResized(imgRGB image.RGBA, angle float64, targetSide int) (*image.RGBA, float64, error) {

	if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return nil, 0, fmt.Errorf("wrong rotation angle (%f) / "+onRotateResized, angle)
	}

	dx, dy := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y

	sideMin := dx
	if dy < sideMin {
		sideMin = dy
	}

	if sideMin <= 0 {
		return nil, 0, fmt.Errorf("wrong image rectangle: %v / "+onRotateResized, imgRGB.Rect)
	}

	if targetSide <= 0 {
		return nil, 0, fmt.Errorf("wrong target side: %d / "+onRotateResized, targetSide)
	}

	scale1 := float64(targetSide) / float64(sideMin)

	mat, err := gocv.ImageToMatRGB(&imgRGB)
	if err != nil {
		return nil, 0, errors.Wrap(err, onRotateResized)
	}
	defer mat.Close()

	var matForResize gocv.Mat
	var center image.Point

	if scale1 == 1 {
		matForResize = mat
		center = image.Point{dx / 2, dy / 2}

	} else {
		matForResize = gocv.NewMat()
		defer matForResize.Close()
		gocv.Resize(mat, &matForResize, image.Point{}, scale1, scale1, gocv.InterpolationDefault)
		center = image.Point{int(float64(dx)*scale1) / 2, int(float64(dy)*scale1) / 2}
	}

	diag := scale1 * math.Sqrt(float64(dx*dx+dy*dy))

	// log.Fatal(targetSideMin, sideMin, scale1, float64(dx)*scale1, float64(dy)*scale1, diag, center)

	scale2 := 1.
	if diag != float64(targetSide) {
		scale2 = float64(targetSide) / diag
	}

	matForRotate := gocv.NewMat()
	defer matForRotate.Close()

	angleDegrees := angle * 180 / math.Pi

	m := gocv.GetRotationMatrix2D(center, angleDegrees, scale2)

	sideX, sideY := int(math.Round(float64(dx)*scale1)), int(math.Round(float64(dy)*scale1))

	gocv.WarpAffine(matForResize, &matForRotate, m, image.Point{sideX, sideY})

	imgRotated, err := matForRotate.ToImage()
	if err != nil {
		return nil, 0, errors.Wrap(err, onRotateResized)
	}

	imgRGBRotated, _ := imgRotated.(*image.RGBA)
	if imgRGBRotated == nil {
		return nil, 0, fmt.Errorf("wrong resized image: %T / "+onRotateResized, imgRGBRotated)
	}

	delta2 := (sideX - sideY) / 2
	var imgRGBFinal *image.RGBA
	if sideX > sideY {
		imgRGBFinal, _ = imgRGBRotated.SubImage(image.Rectangle{image.Point{delta2, 0}, image.Point{delta2 + targetSide, sideY}}).(*image.RGBA)
	} else {
		imgRGBFinal, _ = imgRGBRotated.SubImage(image.Rectangle{image.Point{0, -delta2}, image.Point{sideX, -delta2 + targetSide}}).(*image.RGBA)
	}

	imgRGBFinal.Rect = Normalize(imgRGBFinal.Rect)

	return imgRGBFinal, scale1 * scale2, nil

}

const onRotate = "on imagelib.Rotate()"

func Rotate(imgRGB image.RGBA, angle float64) (*image.RGBA, error) {

	if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return nil, fmt.Errorf("wrong rotation angle (%f) / "+onRotate, angle)
	}

	dx, dy := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y

	mat, err := gocv.ImageToMatRGB(&imgRGB)
	if err != nil {
		return nil, errors.Wrap(err, onRotate)
	}
	defer mat.Close()
	matForRotate := gocv.NewMat()
	defer matForRotate.Close()

	center := image.Point{dx / 2, dy / 2}
	angleDegrees := angle * 180 / math.Pi

	m := gocv.GetRotationMatrix2D(center, angleDegrees, 1)

	gocv.WarpAffine(mat, &matForRotate, m, image.Point{dx, dy})

	imgRotated, err := matForRotate.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onRotate)
	}

	imgRGBRotated, _ := imgRotated.(*image.RGBA)
	if imgRGBRotated == nil {
		return nil, fmt.Errorf("wrong rotated image: %T / "+onRotate, imgRGBRotated)
	}

	return imgRGBRotated, nil
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
