package opencvlib

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/mathlib/plane"
	"image"
	"math"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"
)

const onImageToMat = "on opencvlib.onImageToMat()"

func ImageToMat(img image.Image) (_ *gocv.Mat, isColor bool, _ error) {
	if img == nil {
		return nil, false, fmt.Errorf(onImageToMat + ": nil img")
	} else if imgGray, ok := img.(*image.Gray); imgGray != nil {
		mat, err := gocv.ImageGrayToMatGray(imgGray)
		if err != nil {
			return nil, false, errors.Wrap(err, onImageToMat)
		}
		return &mat, false, nil
	} else if ok {
		return nil, false, fmt.Errorf(onImageToMat + ": nil img.(*Gray)")

	}
	mat, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		return nil, false, errors.Wrap(err, onImageToMat)
	}
	return &mat, true, nil
}

func ImgSize(mat *gocv.Mat) image.Point {
	var size image.Point

	if mat == nil {
		return size
	}

	matSize := mat.Size()
	if len(matSize) > 0 {
		size.Y = matSize[0]
		if len(matSize) > 1 {
			size.X = matSize[1]
		}
	}

	return size
}

const onResize = "on Resize()"

func Resize(imgRGB image.RGBA, scale float64) (*image.RGBA, float64, error) {
	if scale == 1 || scale == 0 {
		return &imgRGB, 1, nil
	} else if scale < 0 || math.IsNaN(scale) || math.IsInf(scale, 0) {
		return nil, 0, fmt.Errorf("wrong resize scale (%f) / "+onResize, scale)
	}

	if imgRGB.Rect.Dx() != imgRGB.Stride {
		// TODO!!! do it faster
		imgRGBCopied := imagelib.ImageToRGBACopied(&imgRGB)
		imgRGB = *imgRGBCopied
	}

	mat, err := gocv.ImageToMatRGB(&imgRGB)
	if err != nil {
		return nil, 0, errors.Wrap(err, onResize)
	}
	defer mat.Close()

	matForResize := gocv.NewMat()
	defer matForResize.Close()

	// gocv.Resize(mat, &matForResize, image.Point{}, scale, scale, gocv.InterpolationDefault)

	pt := image.Point{int(math.Round(float64(imgRGB.Rect.Dx()) * scale)), int(math.Round(float64(imgRGB.Rect.Dy()) * scale))}

	gocv.Resize(mat, &matForResize, pt, 0, 0, gocv.InterpolationDefault)

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

const onResizeGray = "on opencvlib.ResizeGray()"

func ResizeGray(imgGray *image.Gray, ratio float64) (*image.Gray, float64, error) {
	if imgGray == nil {
		return nil, 0, errors.New("imgGray == nil / " + onResizeGray)
	} else if ratio == 1 || ratio == 0 {
		return imgGray, 1, nil
	} else if ratio < 0 || math.IsNaN(ratio) || math.IsInf(ratio, 0) {
		return nil, 0, fmt.Errorf("wrong resize ratio (%f) / "+onResizeGray, ratio)
	}

	mat, err := gocv.ImageGrayToMatGray(imgGray)
	if err != nil {
		return nil, 0, errors.Wrap(err, onResizeGray)
	}
	defer mat.Close()

	matForResize := gocv.NewMat()
	defer matForResize.Close()

	gocv.Resize(mat, &matForResize, image.Point{}, ratio, ratio, gocv.InterpolationDefault)

	imgResized, err := matForResize.ToImage()
	if err != nil {
		return nil, 0, errors.Wrap(err, onResizeGray)
	}

	grayResized, ok := imgResized.(*image.Gray)
	if !ok {
		return nil, 0, fmt.Errorf("resized image has wrong type: %T / "+onResizeGray, grayResized)
	}

	return grayResized, ratio, nil
}

const onRotate = "on Rotate()"

func Rotate(imgRGB *image.RGBA, angle float64) (*image.RGBA, error) {

	if imgRGB == nil {
		return nil, errors.New("imgRGB == nil / " + onRotate)
	} else if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return nil, fmt.Errorf("wrong rotation angle (%f) / "+onRotate, angle)
	}

	dx, dy := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y

	mat, err := gocv.ImageToMatRGB(imgRGB)
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

const onRotateGray = "on RotateGray()"

func RotateGray(imgGray *image.Gray, angle float64) (*image.Gray, error) {

	if imgGray == nil {
		return nil, errors.New("imgGray == nil / " + onRotateGray)
	} else if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return nil, fmt.Errorf("wrong rotation angle (%f) / "+onRotateGray, angle)
	}

	dx, dy := imgGray.Rect.Max.X-imgGray.Rect.Min.X, imgGray.Rect.Max.Y-imgGray.Rect.Min.Y

	mat, err := gocv.ImageGrayToMatGray(imgGray)
	if err != nil {
		return nil, errors.Wrap(err, onRotateGray)
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
		return nil, errors.Wrap(err, onRotateGray)
	}

	imgGrayRotated, _ := imgRotated.(*image.Gray)
	if imgGrayRotated == nil {
		return nil, fmt.Errorf("imgGrayRotated == nil (imgRotated: %T) / "+onRotateGray, imgRotated)
	}

	return imgGrayRotated, nil
}

const onPositionImage = "on PositionImage()"

func PositionImage(imgRGBA0 image.RGBA, scale float64, rotation plane.XToYAngle, imgSideX, imgSideY int, l logger.Operator) (_ *image.RGBA, leftTop image.Point, _ error) {

	// TODO!!! don't convert gocv.Mat to imgRGBA

	imgRGBA, _, err := Resize(imgRGBA0, scale)
	if err != nil {
		return nil, leftTop, errors.Wrap(err, onPositionImage)
	} else if imgRGBA == nil {
		return nil, leftTop, errors.New("imgRGBA == nil / " + onPositionImage)
	}

	//imagelib.Save(imgRGBA, "test.png")

	var imgRGBRotated *image.RGBA

	if rotation == 0 {
		imgRGBRotated = imgRGBA

	} else {

		// TODO!!! why it doesn't work properly???
		// imgRGBRotated, err = Rotate(*imgRGBCentered, -rotation)

		imgRGBRotated, err = Rotate(imagelib.ImageToRGBACopied(imgRGBA), float64(-rotation))
		if err != nil {
			return nil, leftTop, errors.Wrap(err, onPositionImage)
		} else if imgRGBRotated == nil {
			return nil, leftTop, errors.New("imgRGBRotated == nil / " + onPositionImage)
		}

	}

	//if l != nil {
	//	l.Image("original.png", imagelib.GetImage1(&imgRGBA0), nil)
	//	l.Image("resized.png", imagelib.GetImage1(imgRGBA), nil)
	//	l.Image("rotated.png", imagelib.GetImage1(imgRGBRotated), nil)
	//}

	imgRect := image.Rectangle{Max: imgRGBRotated.Rect.Canon().Size()}
	imgRGBRotated.Rect = imgRect

	var imgRGBFinal *image.RGBA
	sub := false

	imgRectSub := imgRect

	if xDelta := imgRect.Max.X - imgSideX; xDelta > 0 {
		imgRectSub.Min.X += xDelta / 2
		imgRectSub.Max.X, sub = imgRectSub.Min.X+imgSideX, true
	}
	if yDelta := imgRect.Max.Y - imgSideY; yDelta > 0 {
		imgRectSub.Min.Y += yDelta / 2
		imgRectSub.Max.Y, sub = imgRectSub.Min.Y+imgSideY, true
	}

	if sub {
		leftTop = imgRectSub.Min
		imgRGBFinal, _ = imgRGBRotated.SubImage(imgRectSub).(*image.RGBA)

	} else {
		imgRGBFinal = imgRGBRotated

	}

	return imgRGBFinal, leftTop, nil
}

const onTranspose = "on Transpose()"

func Transpose(imgRGB *image.RGBA) (*image.RGBA, error) {
	if imgRGB == nil {
		return nil, errors.New("imgRGB == nil / " + onTranspose)
	}
	mat, err := gocv.ImageToMatRGB(imgRGB)
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

const onMorphExGray = "on MorphExGray()"

func MorphExGray(imgGray *image.Gray, morphType gocv.MorphType, size int) (*image.Gray, error) {

	if imgGray == nil {
		return nil, errors.New("imgGray == nil / " + onMorphExGray)
	}

	mat, err := gocv.ImageGrayToMatGray(imgGray)
	if err != nil {
		return nil, errors.Wrap(err, onMorphExGray)
	}
	defer mat.Close()

	matForTransform := gocv.NewMat()
	defer matForTransform.Close()

	morphEl := gocv.GetStructuringElement(gocv.MorphRect, image.Point{size, size})
	gocv.MorphologyEx(mat, &matForTransform, morphType, morphEl)

	imgTransformed, err := matForTransform.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onMorphExGray)
	}

	imgGrayTransformed, ok := imgTransformed.(*image.Gray)
	if !ok {
		return nil, fmt.Errorf("transposed image has wrong type: %T / "+onMorphExGray, imgGray)
	}

	return imgGrayTransformed, nil
}

const onPrepare = "on Prepare()"

func Prepare(mat gocv.Mat, colorConversionCode gocv.ColorConversionCode, scale float64) (image.Image, error) {
	if !(scale > 0) || math.IsInf(scale, 1) {
		return nil, fmt.Errorf("wrong scale: %f / "+onPrepare, scale)
	}

	dims := mat.Size()
	if len(dims) != 2 {
		return nil, fmt.Errorf("wrong mat size: %+v / "+onPrepare, dims)
	}

	var matConverted *gocv.Mat

	if colorConversionCode < 0 {
		matConverted = &mat
	} else {
		matColorConverted := gocv.NewMat()
		defer matColorConverted.Close()

		gocv.CvtColor(mat, &matColorConverted, colorConversionCode)
		matConverted = &matColorConverted
	}

	if scale != 1 {
		matForResize := gocv.NewMat()
		defer matForResize.Close()

		gocv.Resize(*matConverted, &matForResize, image.Point{}, scale, scale, gocv.InterpolationDefault)
		matConverted = &matForResize
	}

	img, err := matConverted.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, onPrepare)
	}

	return img, nil
}

//const onPositionImage = "on PositionImage()"
//
//func PositionImage(imgRGBA0 *image.RGBA, scale0 float64, rotation plane.XToYAngle, imgSideX, imgSideY int, l logger.Operator) (*image.RGBA, error) {
//
//	// TODO!!! don't convert gocv.Mat to imgRGBA
//
//	imgRGBA, _, err := Resize(imgRGBA0, scale0)
//	if err != nil {
//		return nil, errors.Wrap(err, onPositionImage)
//	} else if imgRGBA == nil {
//		return nil, errors.New("imgRGBA == nil / " + onPositionImage)
//	}
//
//	var imgRGBRotated *image.RGBA
//
//	if rotation == 0 {
//		imgRGBRotated = imgRGBA
//
//	} else {
//
//		// TODO!!! why it doesn't work properly???
//		// imgRGBRotated, err = Rotate(*imgRGBCentered, -rotation)
//
//		imgRGBRotated, err = Rotate(imagelib.ImageToRGBACopied(imgRGBA), float64(-rotation))
//		if err != nil {
//			return nil, errors.Wrap(err, onPositionImage)
//		} else if imgRGBRotated == nil {
//			return nil, errors.New("imgRGBRotated == nil / " + onPositionImage)
//		}
//
//	}
//
//	if l != nil {
//		l.Image("original.png", imagelib.GetImage1(imgRGBA0), nil)
//		l.Image("resized.png", imagelib.GetImage1(imgRGBA), nil)
//		l.Image("rotated.png", imagelib.GetImage1(imgRGBRotated), nil)
//	}
//
//	imgRect := image.Rectangle{Max: imgRGBRotated.Rect.Canon().Size()}
//	imgRGBRotated.Rect = imgRect
//
//	var imgRGBFinal *image.RGBA
//	sub := false
//
//	imgRectSub := imgRect
//
//	if xDelta := imgRect.Max.X - imgSideX; xDelta > 0 {
//		imgRectSub.Min.X += xDelta / 2
//		imgRectSub.Max.X, sub = imgRectSub.Min.X+imgSideX, true
//	}
//	if yDelta := imgRect.Max.Y - imgSideY; yDelta > 0 {
//		imgRectSub.Min.Y += yDelta / 2
//		imgRectSub.Max.Y, sub = imgRectSub.Min.Y+imgSideY, true
//	}
//
//	if sub {
//		imgRGBFinal, _ = imgRGBRotated.SubImage(imgRectSub).(*image.RGBA)
//
//	} else {
//		imgRGBFinal = imgRGBRotated
//
//	}
//
//	return imgRGBFinal, nil
//}
//
//const onRotateResized = "on RotateResized()"
//
//// DEPRECATED!!!
//func RotateResized(imgRGB image.RGBA, angle plane.XToYAngle, targetSide int) (*image.RGBA, float64, error) {
//
//	if math.IsNaN(float64(angle)) || math.IsInf(float64(angle), 0) {
//		return nil, 0, fmt.Errorf("wrong rotation angle (%f) / "+onRotateResized, angle)
//	}
//
//	dx, dy := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y
//
//	sideMin := dx
//	if dy < sideMin {
//		sideMin = dy
//	}
//
//	if sideMin <= 0 {
//		return nil, 0, fmt.Errorf("wrong image rectangle: %v / "+onRotateResized, imgRGB.Rect)
//	}
//
//	if targetSide <= 0 {
//		return nil, 0, fmt.Errorf("wrong target side: %d / "+onRotateResized, targetSide)
//	}
//
//	scale1 := float64(targetSide) / float64(sideMin)
//
//	mat, err := gocv.ImageToMatRGB(&imgRGB)
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onRotateResized)
//	}
//	defer mat.Close()
//
//	var matForResize gocv.Mat
//	var center image.Point
//
//	if scale1 == 1 {
//		matForResize = mat
//		center = image.Point{dx / 2, dy / 2}
//
//	} else {
//		matForResize = gocv.NewMat()
//		defer matForResize.Close()
//		gocv.Resize(mat, &matForResize, image.Point{}, scale1, scale1, gocv.InterpolationDefault)
//		center = image.Point{int(float64(dx)*scale1) / 2, int(float64(dy)*scale1) / 2}
//	}
//
//	diag := scale1 * math.Sqrt(float64(dx*dx+dy*dy))
//
//	// log.Fatal(targetSideMin, sideMin, scale1, float64(dx)*scale1, float64(dy)*scale1, diag, center)
//
//	scale2 := 1.
//	if diag != float64(targetSide) {
//		scale2 = float64(targetSide) / diag
//	}
//
//	matForRotate := gocv.NewMat()
//	defer matForRotate.Close()
//
//	angleDegrees := float64(angle * 180 / math.Pi)
//
//	m := gocv.GetRotationMatrix2D(center, angleDegrees, scale2)
//
//	sideX, sideY := int(math.Round(float64(dx)*scale1)), int(math.Round(float64(dy)*scale1))
//
//	gocv.WarpAffine(matForResize, &matForRotate, m, image.Point{sideX, sideY})
//
//	imgRotated, err := matForRotate.ToImage()
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onRotateResized)
//	}
//
//	imgRGBRotated, _ := imgRotated.(*image.RGBA)
//	if imgRGBRotated == nil {
//		return nil, 0, fmt.Errorf("wrong resized image: %T / "+onRotateResized, imgRGBRotated)
//	}
//
//	delta2 := (sideX - sideY) / 2
//	var imgRGBFinal *image.RGBA
//	if sideX > sideY {
//		imgRGBFinal, _ = imgRGBRotated.SubImage(image.Rectangle{image.Point{delta2, 0}, image.Point{delta2 + targetSide, sideY}}).(*image.RGBA)
//	} else {
//		imgRGBFinal, _ = imgRGBRotated.SubImage(image.Rectangle{image.Point{0, -delta2}, image.Point{sideX, -delta2 + targetSide}}).(*image.RGBA)
//	}
//
//	imgRGBFinal.Rect = imagelib.Normalize(imgRGBFinal.Rect)
//
//	return imgRGBFinal, scale1 * scale2, nil
//
//}

//const onResizeToRange = "on ResizeToRange()"
//
//func ResizeToRange(imgRGB image.RGBA, dpm float64, dpmRange [2]float64) (*image.RGBA, float64, error) {
//	if !(dpm > 0 && !math.IsInf(dpm, 1)) {
//		return nil, 0, fmt.Errorf("wrong dpm: %f / "+onResizeToRange, dpm)
//	}
//	if dpm >= dpmRange[0] && dpm <= dpmRange[1] {
//		return &imgRGB, dpm, nil
//	}
//
//	imgRGBResized, resizeRatio, err := Resize(imgRGB, 0.5*(dpmRange[0]+dpmRange[1])/dpm)
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onResizeToRange)
//	} else if imgRGBResized == nil {
//		return nil, 0, errors.New("resized img == nil / " + onResizeToRange)
//	}
//
//	return imgRGBResized, dpm * resizeRatio, nil
//}

//const onResize = "on Resize()"
//
//func Resize(imgRGB *image.RGBA, ratio float64) (*image.RGBA, float64, error) {
//	if imgRGB == nil {
//		return nil, 0, errors.New("imgRGB == nil / " + onResize)
//	} else if ratio == 1 || ratio == 0 {
//		return imgRGB, 1, nil
//	} else if ratio < 0 || math.IsNaN(ratio) || math.IsInf(ratio, 0) {
//		return nil, 0, fmt.Errorf("wrong resize ratio (%f) / "+onResize, ratio)
//	}
//
//	mat, err := gocv.ImageToMatRGB(imgRGB)
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onResize)
//	}
//	defer mat.Close()
//
//	matForResize := gocv.NewMat()
//	defer matForResize.Close()
//
//	gocv.Resize(mat, &matForResize, image.Point{}, ratio, ratio, gocv.InterpolationDefault)
//
//	imgResized, err := matForResize.ToImage()
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onResize)
//	}
//
//	rgbaResized, ok := imgResized.(*image.RGBA)
//	if !ok {
//		return nil, 0, fmt.Errorf("resized image has wrong type: %T / "+onResize, rgbaResized)
//	}
//
//	return rgbaResized, ratio, nil
//}
