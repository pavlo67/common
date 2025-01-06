package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation1CntMask{}

// type separation1CntMask struct {
// 	lyr      *layers.Layer
// 	thr      pix.Value
// 	whiteMin int
// }

// func Separation1Cnt(thr pix.Value, whiteMin int) convolution.Mask {
// 	return &separation1CntMask{
// 		thr:      thr,
// 		whiteMin: whiteMin,
// 	}
// }

// const onSeparation1CntPrepare = "on Separation1Cnt.Prepare()"

// func (mask *separation1CntMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparation1CntPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparation1CntPrepare, onData)
// 	}

// 	return nil
// }

// func (mask separation1CntMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation1CntMask) Info() common.Map {
// 	return common.Map{
// 		"name":     fmt.Sprintf("sep1_cnt_%d_%d", mask.thr, mask.whiteMin),
// 		"thr":      mask.thr,
// 		"whiteMin": mask.whiteMin,
// 	}
// }

// func (mask separation1CntMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	if x <= lyr.Rect.Min.X || x >= lyr.Rect.Max.X-1 || y <= lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y-1 {
// 		return 0
// 	}

// 	xMin, xMax, yMin, yMax := x-1, x+2, y-1, y+2

// 	offset, cnt := (yMin-lyr.Rect.Min.Y)*lyr.Stride+(-lyr.Rect.Min.X), 0
// 	for yi := yMin; yi < yMax; yi++ {
// 		for xi := xMin; xi < xMax; xi++ {
// 			if lyr.Pix[offset+xi] >= mask.thr {
// 				cnt++
// 			}
// 		}
// 		offset += lyr.Stride
// 	}

// 	if cnt >= mask.whiteMin {
// 		return pix.ValueMax
// 	}

// 	return lyr.Pix[(y-lyr.Rect.Min.Y)*lyr.Stride+(x-lyr.Rect.Min.X)]
// }
