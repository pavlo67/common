package convolution

//import (
//	"fmt"
//	"image"
//
//	"github.com/pavlo67/common/common/imagelib"
//
//	"github.com/pavlo67/common/common"
//
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//var _ convolution.Mask = &neighborsSimilarMask{}
//
//type neighborsSimilarMask struct {
//	imgRGB      *image.RGBA
//	pixDeltaMax pix.ValueDelta
//	thrClose   int
//}
//
//const onNeighborsSimilar = "on NeighborsSimilar()"
//
//func NeighborsSimilar(colorDelta layer.ValueDelta, thrClose int) (Mask, error) {
//	return &neighborsSimilarMask{
//		pixDeltaMax: colorDelta,
//		thrClose:   thrClose,
//	}, nil
//}
//
//func (mask *neighborsSimilarMask) Side() int {
//	return 1
//}
//
//const onNeighborsSimilarPrepare = "on NeighborsSimilar.GetNext()"
//
//func (mask *neighborsSimilarMask) GetNext(onData interface{}) error {
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
//			return fmt.Errorf("onData == nil (%#v) / "+onNeighborsSimilarPrepare, onData)
//		}
//		return fmt.Errorf("wrong data (%T) / "+onNeighborsSimilarPrepare, onData)
//	}
//
//	return nil
//}
//
//func (mask neighborsSimilarMask) Info() common.Mappings {
//	info := common.Mappings{
//		"name":        "neighbors_similar",
//		"pixDeltaMax": mask.pixDeltaMax,
//	}
//
//	return info
//}
//
//func (mask *neighborsSimilarMask) Calculate(x, y int) layer.Value {
//
//	imgRGB := mask.imgRGB
//
//	offsetCenter := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
//	clr := imgRGB.Pix[offsetCenter : offsetCenter+3]
//
//	offsets, cnt := make([]int, 0, 4), 0
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
//		//delta := (pix.ValueDelta(imgRGB.Pix[offset])-pix.ValueDelta(clr[0]))*(pix.ValueDelta(imgRGB.Pix[offset])-pix.ValueDelta(clr[0])) +
//		//	(pix.ValueDelta(imgRGB.Pix[offset+1])-pix.ValueDelta(clr[1]))*(pix.ValueDelta(imgRGB.Pix[offset+1])-pix.ValueDelta(clr[1])) +
//		//	(pix.ValueDelta(imgRGB.Pix[offset+2])-pix.ValueDelta(clr[2]))*(pix.ValueDelta(imgRGB.Pix[offset+2])-pix.ValueDelta(clr[2]))
//
//		var delta pix.ValueDelta
//		deltaR := (pix.ValueDelta(imgRGB.Pix[offset]) - pix.ValueDelta(clr[0]))
//		if deltaR >= 0 {
//			delta += deltaR
//		} else {
//			delta -= deltaR
//		}
//		deltaG := (pix.ValueDelta(imgRGB.Pix[offset+1]) - pix.ValueDelta(clr[1]))
//		if deltaG >= 0 {
//			delta += deltaG
//		} else {
//			delta -= deltaG
//		}
//		deltaB := (pix.ValueDelta(imgRGB.Pix[offset+2]) - pix.ValueDelta(clr[2]))
//		if deltaB >= 0 {
//			delta += deltaB
//		} else {
//			delta -= deltaB
//		}
//
//		//fmt.Println(deltaR, deltaB, deltaG, delta, mask.pixDeltaMax)
//
//		if delta <= mask.pixDeltaMax {
//			cnt++
//		}
//	}
//
//	if mask.thrClose <= 0 {
//		return layer.Value(cnt) * layer.Value(0x30)
//	} else if cnt >= mask.thrClose {
//		return layer.ValueMax
//	}
//
//	return 0
//}
//
//const onNeighborsSimilarInitStat = "on NeighborsSimilar.InitStat()"
//
//func (mask *neighborsSimilarMask) InitStat(thrClose float64) error {
//	//mask.pixDeltaMax = int32(math.Round(thrClose))
//	//if mask.pixDeltaMax <= 0 {
//	//	return fmt.Errorf("wrong thrClose: %f / "+onNeighborsSimilarInitStat, thrClose)
//	//}
//	//
//	//mask.colorCounts = map[int32]int64{}
//	return nil
//}
//
//func (mask *neighborsSimilarMask) Classes() interface{} {
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
