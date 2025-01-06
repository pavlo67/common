package convolution_rgb

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	"image"
	"math"
	"strconv"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/imagelib/frame"
)

var _ Mask = &averagingRGBMask{}

type averagingRGBMask struct {
	imgRGB *image.RGBA
	radius int
}

func RGBAveraging(radiusUint uint) Mask {
	return &averagingRGBMask{
		radius: int(radiusUint),
	}
}

const onAveragingRGBPrepare = "on RGBAveraging.GetNext()"

func (mask *averagingRGBMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case image.RGBA:
		mask.imgRGB = &v
	case *image.RGBA:
		mask.imgRGB = v
	case frame.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case *frame.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case frame.Frame:
		mask.imgRGB = &v.RGBA
	case *frame.Frame:
		mask.imgRGB = &v.RGBA
	}
	if mask.imgRGB == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onAveragingRGBPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onAveragingRGBPrepare, onData)
	}

	return nil
}

//func (mask averagingRGBMask) Classes() *imager.Metrics {
//	return nil
//}

func (mask averagingRGBMask) Info() common.Map {
	return common.Map{
		"name": "av_rgb_" + strconv.Itoa(mask.radius),
		"side": mask.radius,
	}
}

func (mask averagingRGBMask) Calculate(x, y int) frame.ValueRGBA {
	imgRGB := mask.imgRGB
	if imgRGB == nil {
		return frame.ValueRGBA{}
	}

	xMin, xMax, yMin, yMax := x-mask.radius, x+mask.radius+1, y-mask.radius, y+mask.radius+1
	if xMin < 0 {
		xMin = 0
	}
	if xMax > imgRGB.Rect.Max.X {
		xMax = imgRGB.Rect.Max.X
	}
	if yMin < 0 {
		yMin = 0
	}
	if yMax > imgRGB.Rect.Max.Y {
		yMax = imgRGB.Rect.Max.Y
	}

	offset, sumR, sumG, sumB, cnt := (yMin-imgRGB.Rect.Min.Y)*imgRGB.Stride+(xMin-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA, pix.ValueSum(0), pix.ValueSum(0), pix.ValueSum(0), 0
	for yi := yMin; yi < yMax; yi++ {
		offsetX := offset

		for xi := xMin; xi < xMax; xi++ {
			sumR += pix.ValueSum(imgRGB.Pix[offsetX])
			sumG += pix.ValueSum(imgRGB.Pix[offsetX+1])
			sumB += pix.ValueSum(imgRGB.Pix[offsetX+2])
			cnt++

			offsetX += coloring.NumColorsRGBA
		}
		offset += imgRGB.Stride
	}

	if cnt <= 0 {
		return frame.ValueRGBA{}
	}

	return frame.ValueRGBA{
		pix.Value(math.Round(float64(sumR) / float64(cnt))),
		pix.Value(math.Round(float64(sumG) / float64(cnt))),
		pix.Value(math.Round(float64(sumB) / float64(cnt))),
		pix.ValueMax,
	}
}
