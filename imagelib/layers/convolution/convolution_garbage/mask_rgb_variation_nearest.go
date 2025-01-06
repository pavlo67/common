package convolution

//import (
//	"fmt"
//	"image"
//	"image/color"
//	"math"
//	"sort"
//
//	"github.com/pavlo67/imagelib/imagelib/color_filter"
//
//	"github.com/pavlo67/common/common/imagelib"
//
//	"github.com/pavlo67/common/common"
//
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//var _ convolution.Mask = &variationNearestMask{}
//
//type variationNearestMask struct {
//	imgRGB      *image.RGBA
//	colorFilter color_filter.Operator
//	shadyValue  layer.Value
//
//	colorDelta  int32
//	colorCounts map[int32]int64
//	inversedBy  layer.Value
//}
//
//const onVariationNearest = "on VariationNearest()"
//
//func VariationNearest(colorFilter color_filter.Operator, shadyValue, inversedBy layer.Value) convolution.Mask {
//	return &variationNearestMask{
//		colorFilter: colorFilter,
//		shadyValue:  shadyValue,
//		colorCounts: map[int32]int64{},
//		inversedBy:  inversedBy,
//	}
//}
//
//func (mask *variationNearestMask) Side() int {
//	return 1
//}
//
//const onVariationNearestPrepare = "on VariationNearest.GetNext()"
//
//func (mask *variationNearestMask) GetNext(onData interface{}) error {
//	switch v := onData.(type) {
//	case image.RGBA:
//		mask.imgRGB = &v
//	case *image.RGBA:
//		mask.imgRGB = v
//case imager.LayerRGBA:
//mask.imgRGB = &v.RGBA
//case *imager.LayerRGBA:
//mask.imgRGB = &v.RGBA
//case imager.Frame:
//mask.imgRGB = &v.RGBA
//case *imager.Frame:
//mask.imgRGB = &v.RGBA
//	}
//	if mask.imgRGB == nil {
//		if onData == nil {
//			return fmt.Errorf("onData == nil (%#v) / "+onVariationNearestPrepare, onData)
//		}
//		return fmt.Errorf("wrong data (%T) / "+onVariationNearestPrepare, onData)
//	}
//
//	return nil
//}
//
//func (mask variationNearestMask) Info() common.Mappings {
//	info := common.Mappings{
//		"name":       "variation_nearest",
//		"inversedBy": mask.inversedBy,
//	}
//
//	if mask.colorFilter != nil {
//		for k, v := range mask.colorFilter.Info() {
//			info["color_filter_"+k] = v
//		}
//	}
//
//	return info
//}
//
//func (mask *variationNearestMask) Calculate(x, y int) layer.Value {
//
//	imgRGB := mask.imgRGB
//
//	offsetCenter := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
//	clr := imgRGB.Pix[offsetCenter : offsetCenter+3]
//
//	if mask.colorFilter != nil && !mask.colorFilter.Test(color.RGBA{clr[0], clr[1], clr[2], 0}) {
//		return mask.shadyValue
//	}
//
//	sum, offsets := pix.ValueSum(0), make([]int, 0, 4)
//
//	if x > imgRGB.Rect.Min.X {
//		offsets = append(offsets, offsetCenter-coloring.NumColorsRGBA)
//	}
//	if x < imgRGB.Rect.Max.X-1 {
//		offsets = append(offsets, offsetCenter+coloring.NumColorsRGBA)
//	}
//	if y > imgRGB.Rect.Min.Y {
//		offsets = append(offsets, offsetCenter-imgRGB.Stride)
//	}
//	if y < imgRGB.Rect.Max.Y-1 {
//		offsets = append(offsets, offsetCenter+imgRGB.Stride)
//	}
//
//	for _, offset := range offsets {
//		sum += pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset])-pix.ValueDelta(clr[0]))*pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset])-pix.ValueDelta(clr[0])) +
//			pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset+1])-pix.ValueDelta(clr[1]))*pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset+1])-pix.ValueDelta(clr[1])) +
//			pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset+2])-pix.ValueDelta(clr[2]))*pix.ValueSum(pix.ValueDelta(imgRGB.Pix[offset+2])-pix.ValueDelta(clr[2]))
//	}
//
//	var v layer.Value
//	if len(offsets) > 0 {
//		v = layer.Value(math.Sqrt(float64(sum) / (3 * float64(len(offsets)))))
//	}
//
//	if mask.inversedBy == 0 {
//		return v
//
//	} else if v < mask.inversedBy {
//		//// TODO: uncomment it if mask.Classes() is needed
//		//if mask.pixDeltaMax > 0 {
//		//	mask.colorCounts[((int32(clr[0])/mask.pixDeltaMax)<<16+(int32(clr[1])/mask.pixDeltaMax)<<8+(int32(clr[2])/mask.pixDeltaMax))]++
//		//}
//		return layer.ValueMax
//
//	}
//
//	return 0
//}
//
//const onVariationNearestInitStat = "on VariationNearest.InitStat()"
//
//func (mask *variationNearestMask) InitStat(thrClose float64) error {
//
//	mask.colorDelta = int32(math.Round(thrClose))
//	if mask.colorDelta <= 0 {
//		return fmt.Errorf("wrong thrClose: %f / "+onVariationNearestInitStat, thrClose)
//	}
//
//	mask.colorCounts = map[int32]int64{}
//	return nil
//}
//
//func (mask *variationNearestMask) Classes() interface{} {
//	if mask == nil {
//		return nil
//	}
//
//	var colorStats []imagelib.ColorStat
//	for colorKey, cnt := range mask.colorCounts {
//		colorStats = append(colorStats, imagelib.ColorStat{colorKey, cnt})
//	}
//	sort.Slice(colorStats, func(i, j int) bool { return colorStats[i].Count > colorStats[j].Count })
//
//	return colorStats
//}
