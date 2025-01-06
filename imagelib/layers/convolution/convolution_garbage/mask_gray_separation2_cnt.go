package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation2CntMask{}

// type separation2CntMask struct {
// 	lyr      *layers.Layer
// 	thr      pix.Value
// 	whiteMin int
// }

// func Separation2Cnt(thr pix.Value, whiteMin int) convolution.Mask {
// 	return &separation2CntMask{
// 		thr:      thr,
// 		whiteMin: whiteMin,
// 	}
// }

// const onSeparation2CntPrepare = "on Separation2Cnt.Prepare()"

// func (mask *separation2CntMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparation2CntPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparation2CntPrepare, onData)
// 	}

// 	return nil
// }

// func (mask separation2CntMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation2CntMask) Info() common.Map {
// 	return common.Map{
// 		"name":     fmt.Sprintf("sep_cnt2_%d_%d", mask.thr, mask.whiteMin),
// 		"thr":      mask.thr,
// 		"whiteMin": mask.whiteMin,
// 	}
// }

// func (mask separation2CntMask) Calculate(x, y int) pix.Value {
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

// 	return 0
// }
