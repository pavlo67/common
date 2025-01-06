package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation3CntMask{}

// type separation3CntMask struct {
// 	lyr       *layers.Layer
// 	threshold pix.Value
// 	blackMax  int
// 	whiteMin  int
// }

// func Separation3Cnt(threshold pix.Value, blackMax, whiteMin int) convolution.Mask {
// 	return &separation3CntMask{
// 		threshold: threshold,
// 		blackMax:  blackMax,
// 		whiteMin:  whiteMin,
// 	}
// }

// const onSeparationByCntPrepare = "on Separation3Cnt.Prepare()"

// func (mask *separation3CntMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparationByCntPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparationByCntPrepare, onData)
// 	}

// 	return nil
// }

// func (mask separation3CntMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation3CntMask) Info() common.Map {
// 	return common.Map{
// 		"name":     fmt.Sprintf("sep_%d_%d_%d", mask.threshold, mask.blackMax, mask.whiteMin),
// 		"blackMax": mask.blackMax,
// 		"whiteMin": mask.whiteMin,
// 		"thrClose": mask.threshold,
// 	}
// }

// func (mask separation3CntMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	if x <= lyr.Rect.Min.X || x >= lyr.Rect.Max.X-1 || y <= lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y-1 {
// 		return pix.ValueMiddle
// 	}

// 	xMin, xMax, yMin, yMax := x-1, x+2, y-1, y+2

// 	offset, cnt := (yMin-lyr.Rect.Min.Y)*lyr.Stride+(-lyr.Rect.Min.X), 0
// 	for yi := yMin; yi < yMax; yi++ {
// 		for xi := xMin; xi < xMax; xi++ {
// 			if lyr.Pix[offset+xi] >= mask.threshold {
// 				cnt++
// 			}
// 		}
// 		offset += lyr.Stride
// 	}

// 	if cnt <= mask.blackMax {
// 		return 0
// 	} else if cnt >= mask.whiteMin {
// 		return pix.ValueMax
// 	}

// 	return pix.ValueMiddle
// }
