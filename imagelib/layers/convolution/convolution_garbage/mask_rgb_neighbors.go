package convolution

//import (
//	"fmt"
//
//	"github.com/pavlo67/maps/mapping/imager"
//
//	"github.com/pavlo67/common/common"
//
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//var _ convolution.Mask = &neighborsMask{}
//
//type neighborsMask struct {
//	imgRGB       *layer.Layer
//	thrClose layer.Value
//}
//
//const onNeighbors = "on Neighbors()"
//
//func Neighbors(thrClose layer.Value) (Mask, error) {
//	return &neighborsMask{thrClose: thrClose}, nil
//}
//
//func (mask *neighborsMask) Side() int {
//	return 1
//}
//
//const onNeighborsPrepare = "on Neighbors.GetNext()"
//
//func (mask *neighborsMask) GetNext(onData interface{}) error {
//	switch v := onData.(type) {
//	case layer.Layer:
//		mask.imgRGB = &v
//	case *layer.Layer:
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
//			return fmt.Errorf("onData == nil (%#v) / "+onNeighborsPrepare, onData)
//		}
//		return fmt.Errorf("wrong data (%T) / "+onNeighborsPrepare, onData)
//	}
//
//	return nil
//}
//
//func (mask neighborsMask) Info() common.Mappings {
//	info := common.Mappings{
//		"name":      "neighbors",
//		"thrClose": mask.thrClose,
//	}
//
//	return info
//}
//
//func (mask *neighborsMask) Calculate(x, y int) layer.Value {
//
//	imgRGB := mask.imgRGB
//	offset := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x - imgRGB.Rect.Min.X)
//	cnt := 0
//
//	if x > imgRGB.Rect.Min.X && imgRGB.Pix[offset-1] >= mask.thrClose {
//		cnt++
//	}
//	if x < imgRGB.Rect.Max.X-1 && imgRGB.Pix[offset+1] >= mask.thrClose {
//		cnt++
//	}
//	if y > imgRGB.Rect.Min.Y && imgRGB.Pix[offset-imgRGB.Stride] >= mask.thrClose {
//		cnt++
//	}
//	if y < imgRGB.Rect.Max.Y-1 && imgRGB.Pix[offset+imgRGB.Stride] >= mask.thrClose {
//		cnt++
//	}
//
//	if cnt > 2 {
//		return layer.ValueMax
//	} else if cnt < 2 {
//		return 0
//	}
//
//	return imgRGB.Pix[offset]
//}
//
//const onNeighborsInitStat = "on Neighbors.InitStat()"
//
//func (mask *neighborsMask) InitStat(thrClose float64) error {
//	//mask.pixDeltaMax = int32(math.Round(thrClose))
//	//if mask.pixDeltaMax <= 0 {
//	//	return fmt.Errorf("wrong thrClose: %f / "+onNeighborsSimilarInitStat, thrClose)
//	//}
//	//
//	//mask.colorCounts = map[int32]int64{}
//	return nil
//}
//
//func (mask *neighborsMask) Classes() *imager.PreparationSettings {
//	if mask == nil {
//		return nil
//	}
//
//	//var colorStats []imagelib.ColorStat
//	//for colorKey, cnt := range mask.colorCounts {
//	//	colorStats = append(colorStats, imagelib.ColorStat{colorKey, cnt})
//	//}
//	//sort.Slice(colorStats, func(i, j int) bool { return colorStats[i].Count > colorStats[j].Count })
//	//
//	//return colorStats
//
//	return nil
//}
