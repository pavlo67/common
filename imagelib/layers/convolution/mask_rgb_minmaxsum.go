package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	frame "github.com/pavlo67/imagelib/frame"
	"github.com/pavlo67/imagelib/layers"
	"image"
	"strconv"

	"github.com/pavlo67/common/common"
)

var _ Mask = &minMaxMask{}

type minMaxMask struct {
	imgRGB *image.RGBA
	radius int
}

const onMinMax = "on RGBVariationMinMaxSumCentered(()"

func RGBVariationMinMaxSumCentered(radius int) (Mask, error) {
	if radius < 0 {
		return nil, fmt.Errorf("wrong radius: %d / "+onMinMax, radius)
	}
	return &minMaxMask{radius: radius}, nil
}

const onMinMaxPrepare = "on RGBVariationMinMaxSumCentered.Prepare()"

func (mask *minMaxMask) Prepare(onData interface{}) error {
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
			return fmt.Errorf("onData == nil (%#v) / "+onMinMaxPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onMinMaxPrepare, onData)
	}

	return nil
}

func (mask minMaxMask) Info() common.Map {
	return common.Map{
		"name": "mm_" + strconv.Itoa(mask.radius),
		"side": mask.radius,
	}
}

func (mask minMaxMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask *minMaxMask) Calculate(x, y int) pix.Value {

	imgRGB := mask.imgRGB

	xWidth, yHeight := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y
	xMin, xMax := max(x-mask.radius, 0), min(x+mask.radius+1, xWidth)
	yMin, yMax := max(y-mask.radius, 0), min(y+mask.radius+1, yHeight)

	offset := (yMin-imgRGB.Rect.Min.Y)*imgRGB.Stride + (xMin-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
	minR, minG, minB := imgRGB.Pix[offset], imgRGB.Pix[offset+1], imgRGB.Pix[offset+2]
	maxR, maxG, maxB := minR, minG, minB

	for y1 := yMin; y1 < yMax; y1++ {
		offsetX := offset
		for x1 := xMin; x1 < xMax; x1++ {
			if imgRGB.Pix[offsetX] < minR {
				minR = imgRGB.Pix[offsetX]
			} else if imgRGB.Pix[offsetX] > maxR {
				maxR = imgRGB.Pix[offsetX]
			}
			if imgRGB.Pix[offsetX+1] < minG {
				minG = imgRGB.Pix[offsetX+1]
			} else if imgRGB.Pix[offsetX+1] > maxG {
				maxG = imgRGB.Pix[offsetX+1]
			}
			if imgRGB.Pix[offsetX+2] < minB {
				minB = imgRGB.Pix[offsetX+2]
			} else if imgRGB.Pix[offsetX+2] > maxB {
				maxB = imgRGB.Pix[offsetX+2]
			}

			offsetX += coloring.NumColorsRGBA
		}
		offset += imgRGB.Stride
	}

	return pix.Value((pix.ValueSum(maxR-minR) + pix.ValueSum(maxG-minG) + pix.ValueSum(maxB-minB)) / 3)

}
