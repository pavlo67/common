package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation2CntNearestMask{}

// type separation2CntNearestMask struct {
// 	lyr *layers.Layer
// }

// func Separation2CntNearest() convolution.Mask {
// 	return &separation2CntNearestMask{}
// }

// const onSeparationCntNearestPrepare = "on Separation2CntNearest.Prepare()"

// func (mask *separation2CntNearestMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparationCntNearestPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparationCntNearestPrepare, onData)
// 	}

// 	return nil
// }

// func (mask separation2CntNearestMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation2CntNearestMask) Info() common.Map {
// 	return common.Map{
// 		"name": fmt.Sprintf("se_fix"),
// 	}
// }

// func (mask separation2CntNearestMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	if x < lyr.Rect.Min.X || x > lyr.Rect.Max.X-1 || y < lyr.Rect.Min.Y || y > lyr.Rect.Max.Y-1 {
// 		return pix.ValueMiddle
// 	}

// 	offset := (y-lyr.Rect.Min.Y)*lyr.Stride + (x - lyr.Rect.Min.X)

// 	v := lyr.Pix[offset]

// 	if v == pix.ValueMiddle {
// 		var whNeighbors, blNeighbors bool
// 		for _, offset1 := range [4]int{offset - 1, offset + 1, offset - lyr.Stride, offset + lyr.Stride} {
// 			if offset1 < 0 || offset1 >= len(lyr.Pix) {
// 				continue
// 			} else if v := lyr.Pix[offset1]; v == 0 {
// 				blNeighbors = true
// 			} else if v == pix.ValueMax {
// 				whNeighbors = true
// 			}
// 		}
// 		if !whNeighbors && blNeighbors {
// 			return 0
// 		}
// 		//if whNeighbors && !blNeighbors {
// 		//	return imager.ValueMax
// 		//}
// 		return pix.ValueMax
// 	}

// 	return v
// }
