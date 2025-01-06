package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &separation3SumLeftTopMask{}

// type separation3SumLeftTopMask struct {
// 	lyr      *layers.Layer
// 	side     int
// 	thrBlack float64
// 	thrWhite float64
// 	lyrTop   *layers.Layer
// 	sideTop  int
// }

// const onSeparationBySum = "on convolution.Separation3SumLeftTop()"

// func Separation3SumLeftTop(side int, thrBlack, thrWhite float64, lyrTop *layers.Layer, sideTop int) (convolution.Mask, error) {
// 	if side < 1 {
// 		return nil, fmt.Errorf("side (%d) < 1 / "+onSeparationBySum, side)
// 	} else if thrBlack >= thrWhite {
// 		return nil, fmt.Errorf("thrBlack (%f) >= thrWhite (%f) / "+onSeparationBySum, thrBlack, thrWhite)
// 	} else if lyrTop != nil && sideTop < 1 {
// 		return nil, fmt.Errorf("sideTop (%d) < 1 / "+onSeparationBySum, sideTop)
// 	}

// 	// log.Fatal(thrBlack, thrWhite)

// 	return &separation3SumLeftTopMask{
// 		side:     side,
// 		thrBlack: thrBlack,
// 		thrWhite: thrWhite,
// 		lyrTop:   lyrTop,
// 		sideTop:  sideTop,
// 	}, nil
// }

// const onSeparationBySumPrepare = "on Separation3SumLeftTop.GetNext()"

// func (mask *separation3SumLeftTopMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}

// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onSeparationBySumPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onSeparationBySumPrepare, onData)

// 	} else if mask.lyrTop != nil {
// 		size, sizeTop := mask.lyr.Rect.Size(), mask.lyrTop.Rect.Size()
// 		if sizeTop.X*mask.sideTop < size.X {
// 			return fmt.Errorf("lyrTop.X (%d) * sideTop (%d) < lyr.X (%d) / "+onSeparationBySumPrepare, sizeTop.X, mask.sideTop, size.X)
// 		}
// 		if sizeTop.Y*mask.sideTop < size.Y {
// 			return fmt.Errorf("lyrTop.Y (%d) * sideTop (%d) < lyr.Y (%d) / "+onSeparationBySumPrepare, sizeTop.Y, mask.sideTop, size.Y)
// 		}

// 	}

// 	return nil
// }

// func (mask separation3SumLeftTopMask) Classes() interface{} {
// 	return nil
// }

// func (mask separation3SumLeftTopMask) Info() common.Map {
// 	return common.Map{
// 		"name":     fmt.Sprintf("sep_%d_%d_%.2f_%.2f", mask.side, mask.sideTop, mask.thrBlack, mask.thrWhite),
// 		"side":     mask.side,
// 		"sideTop":  mask.sideTop,
// 		"thrBlack": mask.thrBlack,
// 		"thrWhite": mask.thrWhite,
// 	}
// }

// func (mask separation3SumLeftTopMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil || x < lyr.Rect.Min.X || x >= lyr.Rect.Max.X || y < lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y {
// 		return pix.ValueMiddle
// 	}

// 	xAbs, yAbs := x-lyr.Rect.Min.X, y-lyr.Rect.Min.Y

// 	var vTop pix.Value
// 	if mask.lyrTop != nil {

// 		// fmt.Print(yAbs, mask.lyrTop.Stride, xAbs)

// 		vTop = mask.lyrTop.Pix[(yAbs/mask.sideTop)*mask.lyrTop.Stride+(xAbs/mask.sideTop)]
// 		if vTop == pix.ValueMax || vTop == 0 {
// 			return vTop
// 		}
// 	}

// 	xMax, yMax := min(x+mask.side, lyr.Rect.Max.X), min(y+mask.side, lyr.Rect.Max.Y)
// 	offset, vSum := (y-lyr.Rect.Min.Y)*lyr.Stride+(-lyr.Rect.Min.X), pix.ValueSum(0)
// 	for yi := y; yi < yMax; yi++ {
// 		for xi := x; xi < xMax; xi++ {
// 			vSum += pix.ValueSum(lyr.Pix[offset+xi])
// 		}
// 		offset += lyr.Stride
// 	}

// 	vRatio := float64(vSum) / (float64(pix.ValueMax) * float64((yMax-y)*(xMax-x)))

// 	if vRatio <= mask.thrBlack {
// 		//if mask.lyrTop != nil {
// 		//	log.Print(x, y, 0)
// 		//}

// 		return 0
// 	} else if vRatio >= mask.thrWhite {
// 		//if mask.lyrTop != nil {
// 		//	log.Print(x, y, imager.ValueMax)
// 		//}

// 		return pix.ValueMax
// 	}

// 	return pix.ValueMiddle
// }
