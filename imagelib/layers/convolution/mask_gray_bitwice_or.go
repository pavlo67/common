package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"

	"github.com/pavlo67/common/common"
)

var _ Mask = &bitwiceOrMask{}

type bitwiceOrMask struct {
	lyr  *layers.Layer
	lyr1 layers.Layer
}

func BitwiceOr(lyr1 layers.Layer) Mask {
	return &bitwiceOrMask{
		lyr1: lyr1,
	}
}

const onBitwiceOrPrepare = "on BitwiceOr.Prepare()"

func (mask *bitwiceOrMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case layers.Layer:
		mask.lyr = &v
	case *layers.Layer:
		mask.lyr = v
	}
	if mask.lyr == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onBitwiceOrPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onBitwiceOrPrepare, onData)
	}
	if mask.lyr.Rect.Min.X < mask.lyr1.Rect.Min.X || mask.lyr.Rect.Min.Y < mask.lyr1.Rect.Min.Y ||
		mask.lyr.Rect.Max.X > mask.lyr1.Rect.Max.X || mask.lyr.Rect.Max.Y > mask.lyr1.Rect.Max.Y {
		return fmt.Errorf("wrong rectangles: lyr1 (%v) vs lyr (%v) / "+onBitwiceOrPrepare, mask.lyr1.Rect, mask.lyr.Rect)
	}

	return nil
}

func (mask bitwiceOrMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask bitwiceOrMask) Info() common.Map {
	return common.Map{
		"name": "and",
	}
}

func (mask bitwiceOrMask) Calculate(x, y int) pix.Value {
	if x < mask.lyr.Rect.Min.X || y < mask.lyr.Rect.Min.Y || x >= mask.lyr.Rect.Max.X || y >= mask.lyr.Rect.Max.Y {
		return 0
	}

	v := mask.lyr.Pix[(y-mask.lyr.Rect.Min.Y)*mask.lyr.Stride+(x-mask.lyr.Rect.Min.X)]
	v1 := mask.lyr1.Pix[(y-mask.lyr1.Rect.Min.Y)*mask.lyr1.Stride+(x-mask.lyr1.Rect.Min.X)]

	if v > 0 || v1 > 0 {
		return pix.ValueMax
	}

	return 0
}
