package convolution

import (
	"fmt"
	"image"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"

	"github.com/pavlo67/imagelib/frame"
	"github.com/pavlo67/imagelib/layers"
)

var _ Mask = &colorRangeMask{}

type colorRangeMask struct {
	imgRGB *image.RGBA
	coloring.ColorRange
	// cnt int32
}

const onColorRange = "on ColorRange()"

func ColorRange(colorRange coloring.ColorRange) Mask {
	return &colorRangeMask{ColorRange: colorRange}
}

func (mask *colorRangeMask) Side() int {
	return 1
}

const onColorRangePrepare = "on ColorRange.GetNext()"

func (mask *colorRangeMask) Prepare(onData interface{}) error {
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
			return fmt.Errorf("onData == nil (%#v) / "+onColorRangePrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onColorRangePrepare, onData)
	}

	// mask.cnt = 0
	return nil
}

func (mask colorRangeMask) Info() common.Map {
	return common.Map{
		"name":     "cr_" + mask.ColorRange.String(),
		"colorMin": mask.ColorRange.ColorMin,
		"colorMax": mask.ColorRange.ColorMax,
	}
}

func (mask colorRangeMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask *colorRangeMask) Calculate(x, y int) pix.Value {
	imgRGB := mask.imgRGB
	offset := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
	r, g, b := imgRGB.Pix[offset], imgRGB.Pix[offset+1], imgRGB.Pix[offset+2]

	if r >= mask.ColorRange.ColorMin.R && r <= mask.ColorRange.ColorMax.R &&
		g >= mask.ColorRange.ColorMin.G && g <= mask.ColorRange.ColorMax.G &&
		b >= mask.ColorRange.ColorMin.B && b <= mask.ColorRange.ColorMax.B {
		return pix.ValueMax
	}

	return 0
}
