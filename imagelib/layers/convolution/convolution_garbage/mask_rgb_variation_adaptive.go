package convolution

// import (
// 	"fmt"
// 	frame2 "github.com/pavlo67/imagelib/frame"
// 	"github.com/pavlo67/imagelib/layers/convolution"
// 	"image"
// 	"math"
// 	"strconv"

// 	"github.com/pavlo67/imagelib/imagelib/pix"

// 	"github.com/pavlo67/common/common/imagelib"

// 	"github.com/pavlo67/common/common"
// )

// var _ convolution.Mask = &variationAdaptiveMask{}

// type variationAdaptiveMask struct {
// 	imgRGB *image.RGBA
// 	radius int

// 	//ctxPix    []layer.Value
// 	//ctxSide   int
// 	//ctxStride int

// 	inversedBy pix.Value
// 	// cnt        int64
// }

// const onVariationAdaptive = "on RGBVariationAdaptiveCentered()"

// func RGBVariationAdaptiveCentered(radius, ctxRadius int, inverseBy pix.Value) (convolution.Mask, error) {
// 	if radius < 1 {
// 		return nil, fmt.Errorf("wrong side: %d / "+onVariationAdaptive, radius)
// 	} else if ctxRadius < 1 {
// 		return nil, fmt.Errorf("wrong ctxRadius: %d / "+onVariationAdaptive, ctxRadius)
// 	}

// 	return &variationAdaptiveMask{radius: radius, inversedBy: inverseBy}, nil // ctxSide: 2*ctxRadius + 1,
// }

// func (mask variationAdaptiveMask) Side() int {
// 	return 1
// }

// const onVariationAdaptivePrepare = "on RGBVariationAdaptiveCentered.Prepare()"

// func (mask *variationAdaptiveMask) Prepare(onData interface{}) error {
// 	switch v := onData.(type) {
// 	case image.RGBA:
// 		mask.imgRGB = &v
// 	case *image.RGBA:
// 		mask.imgRGB = v
// 	case frame2.LayerRGBA:
// 		mask.imgRGB = &v.RGBA
// 	case *frame2.LayerRGBA:
// 		mask.imgRGB = &v.RGBA
// 	case frame2.Frame:
// 		mask.imgRGB = &v.RGBA
// 	case *frame2.Frame:
// 		mask.imgRGB = &v.RGBA
// 	}
// 	if mask.imgRGB == nil {
// 		if onData == nil {
// 			return fmt.Errorf("onData == nil (%#v) / "+onVariationAdaptivePrepare, onData)
// 		}
// 		return fmt.Errorf("wrong data (%T) / "+onVariationAdaptivePrepare, onData)
// 	}

// 	// xWidth, yHeight := mask.imgRGB.Rect.Max.X-mask.imgRGB.Rect.Min.X, mask.imgRGB.Rect.Max.Y-mask.imgRGB.Rect.Min.Y
// 	//xCtxMax, yCtxMax := xWidth/mask.ctxSide, yHeight/mask.ctxSide
// 	//if xWidth%mask.ctxSide > 0 {
// 	//	xCtxMax++
// 	//}
// 	//if yHeight%mask.ctxSide > 0 {
// 	//	yCtxMax++
// 	//}
// 	//
// 	//mask.ctxStride = xCtxMax
// 	//mask.ctxPix = make([]layer.Value, mask.ctxStride*yCtxMax)
// 	//
// 	//var ctxOffset int
// 	//for yCtx := 0; yCtx < yCtxMax; yCtx++ {
// 	//	for xCtx := 0; xCtx < xCtxMax; xCtx++ {
// 	//		xMin, xMax, yMin, yMax := xCtx*mask.ctxSide, (xCtx+1)*mask.ctxSide, yCtx*mask.ctxSide, (yCtx+1)*mask.ctxSide
// 	//		if xMax > xWidth {
// 	//			xMax = xWidth
// 	//		}
// 	//		if yMax > yHeight {
// 	//			yMax = yHeight
// 	//		}
// 	//
// 	//		offsetAvg := (yMin+yMax)/2*mask.imgRGB.Stride + (xMin+xMax)/2
// 	//		rAvg, gAvg, bAvg := pix.ValueSum(mask.imgRGB.Pix[offsetAvg]), pix.ValueSum(mask.imgRGB.Pix[offsetAvg+1]), pix.ValueSum(mask.imgRGB.Pix[offsetAvg+2])
// 	//
// 	//		offset := yMin * mask.imgRGB.Stride
// 	//		sum, cnt := pix.ValueSum(0), float64((xMax-xMin)*(yMax-yMin))
// 	//
// 	//		for y := yMin; y < yMax; y++ {
// 	//			for x := xMin; x < xMax; x++ {
// 	//				r, g, b := pix.ValueSum(mask.imgRGB.Pix[offset]), pix.ValueSum(mask.imgRGB.Pix[offset+1]), pix.ValueSum(mask.imgRGB.Pix[offset+2])
// 	//				sum += (r-rAvg)*(r-rAvg) + (g-gAvg)*(g-gAvg) + (b-bAvg)*(b-bAvg)
// 	//			}
// 	//			offset += mask.imgRGB.Stride
// 	//		}
// 	//		mask.ctxPix[ctxOffset+xCtx] = layer.Value(math.Round(math.Sqrt(float64(sum)/3) / cnt))
// 	//
// 	//	}
// 	//	ctxOffset += mask.ctxStride
// 	//}

// 	// mask.cnt = 0
// 	return nil
// }

// func (mask variationAdaptiveMask) Info() common.Map {
// 	return common.Map{
// 		"name":       "variationAdaptive_" + strconv.Itoa(mask.radius) + "_" + strconv.Itoa(int(mask.inversedBy)),
// 		"side":       mask.radius,
// 		"inversedBy": mask.inversedBy,
// 	}
// }

// func (mask variationAdaptiveMask) Classes() interface{} {
// 	//var whiteRatio float64
// 	//
// 	//if mask.imgRGB != nil {
// 	//	size := mask.imgRGB.Rect.Size()
// 	//
// 	//	if layerLen := float64(size.X * size.Y); layerLen > 0 {
// 	//		whiteRatio = float64(mask.cnt) / layerLen
// 	//	}
// 	//
// 	//}
// 	//
// 	//return &imager.Metrics{
// 	//	WhRat: whiteRatio,
// 	//	//BlRat:          0,
// 	//	//Avg:        0,
// 	//	//Criterion: 0,
// 	//	//Classes:            nil,
// 	//}

// 	return nil
// }

// func (mask *variationAdaptiveMask) Calculate(x, y int) pix.Value {

// 	imgRGB := mask.imgRGB

// 	//xCanon, yCanon := x-imgRGB.Rect.Min.X, y-imgRGB.Rect.Min.Y
// 	//contextOffset := (mask.ctxStride * (yCanon / mask.ctxSide)) + (xCanon / mask.ctxSide)
// 	//ctxPix := mask.ctxPix[contextOffset]
// 	//if ctxPix == 0 {
// 	//	return 0
// 	//}

// 	xWidth, yHeight := imgRGB.Rect.Max.X-imgRGB.Rect.Min.X, imgRGB.Rect.Max.Y-imgRGB.Rect.Min.Y
// 	xMin, xMax, yMin, yMax := x-mask.radius, x+mask.radius+1, y-mask.radius, y+mask.radius+1
// 	if xMin < 0 {
// 		xMin = 0
// 	}
// 	if xMax > xWidth {
// 		xMax = xWidth
// 	}
// 	if yMin < 0 {
// 		yMin = 0
// 	}
// 	if yMax > yHeight {
// 		yMax = yHeight
// 	}

// 	offsetCenter := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
// 	clr := imgRGB.Pix[offsetCenter : offsetCenter+3]

// 	offset := (yMin-imgRGB.Rect.Min.Y)*imgRGB.Stride + (xMin-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA

// 	sum, cnt := pix.ValueSum(0), float64((xMax-xMin)*(yMax-yMin))

// 	for y1 := yMin; y1 < yMax; y1++ {
// 		offsetX := offset
// 		for x1 := xMin; x1 < xMax; x1++ {
// 			sum += (pix.ValueSum(imgRGB.Pix[offsetX])-pix.ValueSum(clr[0]))*(pix.ValueSum(imgRGB.Pix[offsetX])-pix.ValueSum(clr[0])) +
// 				(pix.ValueSum(imgRGB.Pix[offsetX+1])-pix.ValueSum(clr[1]))*(pix.ValueSum(imgRGB.Pix[offsetX+1])-pix.ValueSum(clr[1])) +
// 				(pix.ValueSum(imgRGB.Pix[offsetX+2])-pix.ValueSum(clr[2]))*(pix.ValueSum(imgRGB.Pix[offsetX+2])-pix.ValueSum(clr[2]))
// 			offsetX += coloring.NumColorsRGBA
// 		}
// 		offset += imgRGB.Stride
// 	}

// 	v := pix.Value(math.Sqrt(float64(sum)/3) / cnt)

// 	if mask.inversedBy == 0 {
// 		//if v == pix.ValueMax {
// 		//	mask.cnt++
// 		//}
// 		return v
// 	} else {
// 		if v < mask.inversedBy {
// 			//mask.cnt++
// 			return pix.ValueMax
// 		}
// 		return 0
// 	}
// }
