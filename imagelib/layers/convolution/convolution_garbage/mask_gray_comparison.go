package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"
// 	"strconv"

// 	"github.com/pavlo67/imagelib/imagelib/pix"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &comparisonMask{}

// type comparisonMask struct {
// 	lyr      *layers.Layer
// 	lyrBase  layers.Layer
// 	valueMax pix.Value
// 	scale    pix.Value
// }

// func Comparison(lyrBase layers.Layer, valueMax pix.Value) convolution.Mask {
// 	var scale pix.Value
// 	if valueMax > 0 {
// 		scale = pix.Value(3 * pix.ValueSum(pix.ValueMax) / pix.ValueSum(4*valueMax))
// 	}

// 	return &comparisonMask{
// 		lyrBase:  lyrBase,
// 		valueMax: valueMax,
// 		scale:    scale,
// 	}
// }

// const onComparisonPrepare = "on Comparison.Prepare()"

// func (mask *comparisonMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onComparisonPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onComparisonPrepare, onData)
// 	}

// 	if mask.lyr.Rect.Min != mask.lyrBase.Rect.Min {
// 		return fmt.Errorf("mask.lyr.Rect.Min (%v) != mask.lyrBase.Rect.Min (%v) / "+onComparisonPrepare, mask.lyr.Rect.Min, mask.lyrBase.Rect.Min)
// 	} else if mask.lyr.Rect.Max != mask.lyrBase.Rect.Max {
// 		return fmt.Errorf("mask.lyr.Rect.Max (%v) != mask.lyrBase.Rect.Max (%v) / "+onComparisonPrepare, mask.lyr.Rect.Max, mask.lyrBase.Rect.Max)
// 	}

// 	return nil
// }

// func (mask comparisonMask) Classes() interface{} {
// 	return nil
// }

// func (mask comparisonMask) Info() common.Map {
// 	return common.Map{
// 		"name":     "cm_" + strconv.Itoa(int(mask.valueMax)) + "_" + mask.lyrBase.Options.StringDefault("name", ""),
// 		"valueMax": mask.valueMax,
// 	}
// }

// func (mask comparisonMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	offset, offsetBase := (y-lyr.Rect.Min.Y)*lyr.Stride+(x-lyr.Rect.Min.X), (y-mask.lyrBase.Rect.Min.Y)*mask.lyrBase.Stride+(x-mask.lyrBase.Rect.Min.X)

// 	if mask.lyrBase.Pix[offsetBase] == 0 {
// 		return pix.ValueMax
// 	}

// 	v := lyr.Pix[offset] / mask.lyrBase.Pix[offsetBase]
// 	if v > mask.valueMax {
// 		return pix.ValueMax
// 	}

// 	return v * mask.scale
// }
