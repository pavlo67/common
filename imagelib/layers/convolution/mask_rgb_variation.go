package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	frame2 "github.com/pavlo67/imagelib/frame"
	"github.com/pavlo67/imagelib/layers"
	"image"
	"math"
	"strconv"

	"github.com/pavlo67/common/common"
)

var _ Mask = &variationMask{}

type variationMask struct {
	imgRGB *image.RGBA
	radius int
}

func RGBVariationCentered(radius uint) Mask {
	return &variationMask{radius: int(radius)}
}

const onVariationPrepare = "on RGBVariationCentered.Prepare()"

func (mask *variationMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case image.RGBA:
		mask.imgRGB = &v
	case *image.RGBA:
		mask.imgRGB = v
	case frame2.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case *frame2.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case frame2.Frame:
		mask.imgRGB = &v.RGBA
	case *frame2.Frame:
		mask.imgRGB = &v.RGBA
	}
	if mask.imgRGB == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onVariationPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onVariationPrepare, onData)
	}

	return nil
}

func (mask variationMask) Info() common.Map {
	return common.Map{
		"name": "vr_" + strconv.Itoa(mask.radius),
		"side": mask.radius,
	}
}

func (mask variationMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask *variationMask) Calculate(x, y int) pix.Value {

	imgRGB := mask.imgRGB

	xWidth, yHeight := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y
	xMin, xMax, yMin, yMax := x-mask.radius, x+mask.radius+1, y-mask.radius, y+mask.radius+1
	if xMin < 0 {
		xMin = 0
	}
	if xMax > xWidth {
		xMax = xWidth
	}
	if yMin < 0 {
		yMin = 0
	}
	if yMax > yHeight {
		yMax = yHeight
	}

	offsetCenter := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
	clr := imgRGB.Pix[offsetCenter : offsetCenter+3]

	offset := (yMin-imgRGB.Rect.Min.Y)*imgRGB.Stride + (xMin-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA

	sum, cnt := pix.ValueSum(0), float64((xMax-xMin)*(yMax-yMin))

	for y1 := yMin; y1 < yMax; y1++ {
		offsetX := offset
		for x1 := xMin; x1 < xMax; x1++ {
			sum += (pix.ValueSum(imgRGB.Pix[offsetX])-pix.ValueSum(clr[0]))*(pix.ValueSum(imgRGB.Pix[offsetX])-pix.ValueSum(clr[0])) +
				(pix.ValueSum(imgRGB.Pix[offsetX+1])-pix.ValueSum(clr[1]))*(pix.ValueSum(imgRGB.Pix[offsetX+1])-pix.ValueSum(clr[1])) +
				(pix.ValueSum(imgRGB.Pix[offsetX+2])-pix.ValueSum(clr[2]))*(pix.ValueSum(imgRGB.Pix[offsetX+2])-pix.ValueSum(clr[2]))

			offsetX += coloring.NumColorsRGBA
		}
		offset += imgRGB.Stride
	}

	return pix.Value(math.Sqrt(float64(sum)/3) / cnt)
}
