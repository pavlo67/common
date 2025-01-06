package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"

	"github.com/pavlo67/common/common"
)

var _ Mask = &bitwiceAndShiftedMask{}

type bitwiceAndShiftedMask struct {
	lyr    *layers.Layer
	lyr1   layers.Layer
	dx, dy int
}

func BitwiceAndShifted(lyr1 layers.Layer, dx, dy int) Mask {
	return &bitwiceAndShiftedMask{lyr1: lyr1, dx: dx, dy: dy}
}

const onBitwiceAndShiftedPrepare = "on BitwiceAndShifted.Prepare()"

func (mask *bitwiceAndShiftedMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case layers.Layer:
		mask.lyr = &v
	case *layers.Layer:
		mask.lyr = v
	}
	if mask.lyr == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onBitwiceAndShiftedPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onBitwiceAndShiftedPrepare, onData)
	}
	if mask.lyr.Rect.Min.X < mask.lyr1.Rect.Min.X || mask.lyr.Rect.Min.Y < mask.lyr1.Rect.Min.Y ||
		mask.lyr.Rect.Max.X > mask.lyr1.Rect.Max.X || mask.lyr.Rect.Max.Y > mask.lyr1.Rect.Max.Y {
		return fmt.Errorf("wrong rectangles: lyr1 (%v) vs lyr (%v) / "+onBitwiceAndShiftedPrepare, mask.lyr1.Rect, mask.lyr.Rect)
	}

	return nil
}

func (mask bitwiceAndShiftedMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask bitwiceAndShiftedMask) Info() common.Map {
	return common.Map{
		"name": "and",
	}
}

func (mask bitwiceAndShiftedMask) Calculate(x, y int) pix.Value {
	if x < mask.lyr.Rect.Min.X || y < mask.lyr.Rect.Min.Y || x >= mask.lyr.Rect.Max.X || y >= mask.lyr.Rect.Max.Y {
		return 0
	}

	x1, y1 := x-mask.dx, y-mask.dy
	if x1 < mask.lyr1.Rect.Min.X || y1 < mask.lyr1.Rect.Min.Y || x1 >= mask.lyr1.Rect.Max.X || y1 >= mask.lyr1.Rect.Max.Y {
		return 0
	}

	v := mask.lyr.Pix[(y-mask.lyr.Rect.Min.Y)*mask.lyr.Stride+(x-mask.lyr.Rect.Min.X)]
	v1 := mask.lyr1.Pix[(y1-mask.lyr1.Rect.Min.Y)*mask.lyr1.Stride+(x1-mask.lyr1.Rect.Min.X)]

	if v > 0 && v1 > 0 {
		return pix.ValueMax
	}

	return 0
}
