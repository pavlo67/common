package convolution

// import (
// 	"fmt"
// 	"github.com/pavlo67/imagelib/imagelib/pix"
// 	"github.com/pavlo67/imagelib/layers"
// 	"github.com/pavlo67/imagelib/layers/convolution"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &distancesMask{}

// type distancesMask struct {
// 	lyr      *layers.Layer
// 	thrClose pix.Value
// 	thrFar   pix.Value
// }

// func Distances(thrClose, thrFar pix.Value) convolution.Mask {
// 	return &distancesMask{
// 		thrClose: thrClose,
// 		thrFar:   thrFar,
// 	}
// }

// const onDistancesPrepare = "on Distances.Prepare()"

// func (mask *distancesMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case layers.Layer:
// 		mask.lyr = &v
// 	case *layers.Layer:
// 		mask.lyr = v
// 	}
// 	if mask.lyr == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onDistancesPrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onDistancesPrepare, onData)
// 	}

// 	return nil
// }

// func (mask distancesMask) Classes() interface{} {
// 	return nil
// }

// func (mask distancesMask) Info() common.Map {
// 	return common.Map{
// 		"name":     fmt.Sprintf("dist_%d_%d", mask.thrClose, mask.thrFar),
// 		"thrClose": mask.thrClose,
// 		"thrFar":   mask.thrFar,
// 	}
// }

// func (mask distancesMask) Calculate(x, y int) pix.Value {
// 	lyr := mask.lyr
// 	if lyr == nil {
// 		return 0
// 	}

// 	if x <= lyr.Rect.Min.X || x >= lyr.Rect.Max.X-1 || y <= lyr.Rect.Min.Y || y >= lyr.Rect.Max.Y-1 {
// 		return pix.ValueMiddle
// 	}

// 	offset := (y-lyr.Rect.Min.Y)*lyr.Stride + (x - lyr.Rect.Min.X)
// 	v := lyr.Pix[offset]
// 	var dist, distMax pix.Value
// 	if v1 := lyr.Pix[offset-1]; v1 >= v {
// 		distMax = v1 - v
// 	} else {
// 		distMax = v - v1
// 	}

// 	for _, offset1 := range [3]int{offset + 1, offset - lyr.Stride, offset + lyr.Stride} {
// 		if v1 := lyr.Pix[offset1]; v1 >= v {
// 			dist = v1 - v
// 		} else {
// 			dist = v - v1
// 		}
// 		if dist > distMax {
// 			distMax = dist
// 		}
// 	}

// 	if distMax >= mask.thrFar {
// 		return pix.ValueMax
// 	} else if distMax <= mask.thrClose {
// 		return 0
// 	}

// 	return pix.ValueMiddle
// }
