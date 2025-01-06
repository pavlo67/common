package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation3CntNearestMask{}

// type separation3CntNearestMask struct {
// 	lyr *layers.Layer
// }

// func Separation3CntNearest() convolution.Mask {
// 	return &separation3CntNearestMask{}
// }

// const onSeparationBoundsPrepare = "on Separation3CntNearest.Prepare()"

// func (mask *separation3CntNearestMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparationBoundsPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparationBoundsPrepare, onData)
// 	}

// 	return nil
// }

// func (mask separation3CntNearestMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation3CntNearestMask) Info() common.Map {
// 	return common.Map{
// 		"name": fmt.Sprintf("se_fix"),
// 	}
// }

// func (mask separation3CntNearestMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	if x <= lyr.Rect.Min.X || x >= lyr.Rect.Max.X-1 || y <= lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y-1 {
// 		return pix.ValueMiddle
// 	}

// 	offset := (y-lyr.Rect.Min.Y)*lyr.Stride + (x - lyr.Rect.Min.X)

// 	switch v := lyr.Pix[offset]; v {
// 	case pix.ValueMax:
// 		if lyr.Pix[offset-1] == 0 || lyr.Pix[offset+1] == 0 || lyr.Pix[offset-lyr.Stride] == 0 || lyr.Pix[offset+lyr.Stride] == 0 {
// 			return pix.ValueMiddle
// 		}
// 		return pix.ValueMax

// 	case 0:
// 		if lyr.Pix[offset-1] == pix.ValueMax || lyr.Pix[offset+1] == pix.ValueMax || lyr.Pix[offset-lyr.Stride] == pix.ValueMax || lyr.Pix[offset+lyr.Stride] == pix.ValueMax {
// 			return pix.ValueMiddle
// 		}
// 		return 0

// 	}

// 	return pix.ValueMiddle
// }
