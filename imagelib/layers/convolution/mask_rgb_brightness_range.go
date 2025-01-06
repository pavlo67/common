package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"
	"image"
	"strconv"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/imagelib/frame"
)

var _ Mask = &brightnessRangeMask{}

type brightnessRangeMask struct {
	imgRGB     *image.RGBA
	classRange uint8
	classMin   uint8
	classMax   uint8
	cnt        int64
}

const onBrightnessRange = "on RGBBrightnessRange()"

func RGBBrightnessRange(classRange, classMin, classMax uint8) (Mask, error) {
	if classRange == 0 {
		return nil, fmt.Errorf("classRange == 0 / " + onBrightnessRange)
	} else if classMax < classMin || classMax > 0xFF/classRange {
		return nil, fmt.Errorf("classRange == %d, classMin = %d, 0xFF/classRange = %d, wrong classMax = %d / "+onBrightnessRange, classRange, classMin, 0xFF/classRange, classMax)
	}

	return &brightnessRangeMask{
		classRange: classRange,
		classMin:   classMin,
		classMax:   classMax,
	}, nil
}

func (mask *brightnessRangeMask) Side() int {
	return 1
}

const onBrightnessRangePrepare = "on RGBBrightnessRange.GetNext()"

func (mask *brightnessRangeMask) Prepare(onData interface{}) error {
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
			return fmt.Errorf("onData == nil (%#v) / "+onBrightnessRangePrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onBrightnessRangePrepare, onData)
	}
	mask.cnt = 0

	return nil
}

func (mask brightnessRangeMask) Info() common.Map {
	return common.Map{
		"name":       "br_" + strconv.Itoa(int(mask.classRange)) + "_" + strconv.Itoa(int(mask.classMin)) + "_" + strconv.Itoa(int(mask.classMax)),
		"classRange": mask.classRange,
		"classMin":   mask.classMin,
		"classMax":   mask.classMax,
	}
}

func (mask brightnessRangeMask) Classes() layers.ClassesCustom {

	return nil
}

func (mask *brightnessRangeMask) Calculate(x, y int) pix.Value {
	offset := (y-mask.imgRGB.Rect.Min.Y)*mask.imgRGB.Stride + (x-mask.imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
	clr := mask.imgRGB.Pix[offset : offset+3]

	brightnessRange := int(clr[0]) + int(clr[1]) + int(clr[2])
	n := pix.Value(brightnessRange / (3 * int(mask.classRange)))
	if n > mask.classMax {
		return mask.classRange * (n - 1)
	} else if n < mask.classMin {
		return mask.classRange * n
	}

	mask.cnt++
	return pix.ValueMax
}
