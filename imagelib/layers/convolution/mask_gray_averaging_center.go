package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"
	"math"
	"strconv"

	"github.com/pavlo67/common/common"
)

var _ Mask = &averagingCenterMask{}

type averagingCenterMask struct {
	lyr               *layers.Layer
	side, left, right int
}

func AveragingCenter(side int) Mask {
	if side < 1 {
		side = 1
	}
	left := side / 2

	return &averagingCenterMask{
		side:  side,
		left:  left,
		right: side - left,
	}
}

const onAveragingCenterPrepare = "on averagingCenter.Prepare()"

func (mask *averagingCenterMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case layers.Layer:
		mask.lyr = &v
	case *layers.Layer:
		mask.lyr = v
	}
	if mask.lyr == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onAveragingCenterPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onAveragingCenterPrepare, onData)
	}

	return nil
}

func (mask averagingCenterMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask averagingCenterMask) Info() common.Map {
	return common.Map{
		"name": "avg_center_" + strconv.Itoa(mask.side),
		"side": mask.side,
	}
}

func (mask averagingCenterMask) Calculate(x, y int) pix.Value {

	lyr := mask.lyr
	if lyr == nil || x < lyr.Rect.Min.X || x >= lyr.Rect.Max.X || y < lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y {
		return pix.ValueMiddle
	}

	xMin, yMin := max(x-mask.left, lyr.Rect.Min.X), max(y-mask.left, lyr.Rect.Min.Y)
	xMax, yMax := min(x+mask.right, lyr.Rect.Max.X), min(y+mask.right, lyr.Rect.Max.Y)
	offset, vSum := (yMin-lyr.Rect.Min.Y)*lyr.Stride+(-lyr.Rect.Min.X), pix.ValueSum(0)
	for yi := yMin; yi < yMax; yi++ {
		for xi := xMin; xi < xMax; xi++ {
			vSum += pix.ValueSum(lyr.Pix[offset+xi])
		}
		offset += lyr.Stride
	}

	return pix.Value(math.Round(float64(vSum) / float64((yMax-yMin)*(xMax-xMin))))
}
