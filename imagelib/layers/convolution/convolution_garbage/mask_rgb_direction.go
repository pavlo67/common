package convolution

//import (
//	"fmt"
//	"image"
//	"sort"
//
//	"github.com/pavlo67/common/common"
//	"github.com/pavlo67/common/common/imagelib"
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//// this algo doesn't work effectively :-(((
//
//var _ convolution.Mask = &directionMask{}
//
//type directionMask struct {
//	imgRGB            *image.RGBA
//	side            []int
//	differenceMax     pix.ValueDelta
//	gapMaxPix         int
//	garbageRatioMax   float64
//	extensionRatioMin float64
//	// inversedBy        layer.Value
//}
//
//const onDirection = "on Bearing()"
//
//func Direction(side []int, differenceMax pix.ValueDelta, gapMaxPix int, garbageRatioMax, extensionRatioMin float64) (Mask, error) {
//	if len(side) < 1 {
//		return nil, fmt.Errorf("empty side: %d / "+onDirection, side)
//	} else if side[0] < 1 {
//		return nil, fmt.Errorf("wrong side[0]: %d / "+onDirection, side[0])
//	}
//	for i := 1; i < len(side); i++ {
//		if side[i] <= side[i-1] {
//			return nil, fmt.Errorf("side[%d]: %d  <= side[%d]: %d / "+onDirection, i, side[i], i-1, side[i-1])
//		}
//	}
//
//	return &directionMask{
//		side:            side,
//		differenceMax:     differenceMax,
//		gapMaxPix:         gapMaxPix,
//		garbageRatioMax:   garbageRatioMax,
//		extensionRatioMin: extensionRatioMin,
//		// inversedBy:        inversedBy,
//	}, nil
//}
//
//func (mask *directionMask) Side() int {
//	return 1
//}
//
//const onDirectionPrepare = "on Bearing.GetNext()"
//
//func (mask *directionMask) GetNext(onData interface{}) error {
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
//			return fmt.Errorf("onData == nil (%#v) / "+onDirectionPrepare, onData)
//		}
//		return fmt.Errorf("wrong data (%T) / "+onDirectionPrepare, onData)
//	}
//
//	return nil
//}
//
//func (mask directionMask) Info() common.Mappings {
//	return common.Mappings{
//		"name":              "direction",
//		"side":            mask.side,
//		"gapMaxPix":         mask.gapMaxPix,
//		"differenceMax":     mask.differenceMax,
//		"garbageRatioMax":   mask.garbageRatioMax,
//		"extensionRatioMin": mask.extensionRatioMin,
//	}
//}
//
//
//func (mask directionMask) Classes() interface{} {
//sizes := mask.imgRGB.Rect.Size()
//if pixLen := sizes.X * sizes.Y; pixLen > 0 {
//return &imager.Metrics{
//WhRat: float64(mask.cnt) / float64(pixLen),
//}
//}
//
//	return nil
//}
//
//func (mask directionMask) Calculate(x, y int) layer.Value {
//
//	imgRGB, rect := mask.imgRGB, mask.imgRGB.Rect
//
//	offsetCenter := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
//	clr := imgRGB.Pix[offsetCenter : offsetCenter+3]
//
//	for _, r := range mask.side {
//		xMin, yMin, xMax, yMax := x-r, y-r, x+r+1, y+r+1
//
//		if xMin < rect.Min.X {
//			xMin = rect.Min.X
//		}
//
//		if xMax >= rect.Max.X {
//			xMax = rect.Max.X
//		}
//
//		if yMin < rect.Min.Y {
//			yMin = rect.Min.Y
//		}
//
//		if yMax >= rect.Max.Y {
//			yMax = rect.Max.Y
//		}
//
//		xWidth, yHeight := xMax-xMin-2, yMax-yMin-2
//		lenPoints := 2 * (xWidth + yHeight)
//		points := make([]image.Contour, lenPoints)
//
//		for x := 0; x < xWidth; x++ {
//			points[x] = image.Contour{xMin + 1 + x, yMin}
//			points[2*xWidth+yHeight-x-1] = image.Contour{xMax - 2 - x, yMax - 1}
//		}
//
//		//log.Print(yMin, yMax, xWidth, yHeight, points)
//
//		for y := 0; y < yHeight; y++ {
//			points[xWidth+y] = image.Contour{xMax - 1, yMin + 1 + y}
//			points[lenPoints-1-y] = image.Contour{xMin, yMax - 2 - y}
//		}
//
//		//log.Fatal(xMin, xMax, xWidth, yHeight, points)
//
//		var groups [][2]int
//		groupFrom, groupTo := -1, -1
//
//		for i, p := range points {
//			offset := (p.Y-rect.Min.Y)*imgRGB.Stride + (p.X-rect.Min.X)*coloring.NumColorsRGBA
//			clr1 := imgRGB.Pix[offset : offset+3]
//
//			var delta pix.ValueDelta
//			for i, v := range clr1 {
//				if v > clr[i] {
//					delta += pix.ValueDelta(v) - pix.ValueDelta(clr[i])
//				} else {
//					delta += pix.ValueDelta(clr[i]) - pix.ValueDelta(v)
//				}
//			}
//			if delta <= mask.differenceMax {
//				if groupFrom < 0 {
//					groupFrom = i
//				}
//				groupTo = -1
//			} else if groupFrom >= 0 {
//				if groupTo < 0 {
//					groupTo = i
//				}
//				if i-groupTo >= mask.gapMaxPix {
//					groups = append(groups, [2]int{groupFrom, groupTo})
//					groupFrom, groupTo = -1, -1
//				}
//			}
//		}
//		if groupFrom >= 0 {
//			if groupTo >= 0 {
//				groups = append(groups, [2]int{groupFrom, groupTo})
//			} else {
//				groups = append(groups, [2]int{groupFrom, lenPoints})
//			}
//		}
//
//		// log.Print(groups)
//
//		if len(groups) > 1 && groups[0][0]+lenPoints-(groups[len(groups)-1][1]-1) <= mask.gapMaxPix {
//			groups[0][0] = groups[len(groups)-1][0]
//			groups = groups[:len(groups)-1]
//		}
//
//		if len(groups) > 2 {
//			sort.Slice(groups, func(i, j int) bool { return groups[i][1]-groups[i][0] > groups[j][1]-groups[j][0] })
//			if float64(groups[2][1]-groups[2][0])/float64(groups[1][1]-groups[1][0]) <= mask.garbageRatioMax {
//				// ok!!!
//				return layer.Value(r)
//			}
//
//			// too many groups of similar points: it's not a point with "linear directed neighborhood"
//			break
//		}
//
//		if len(groups) == 2 {
//			// ok!!!
//			return layer.Value(r)
//		}
//
//		if len(groups) != 1 || float64((groups[0][1]+lenPoints-groups[0][0])%lenPoints)/float64(lenPoints) < mask.extensionRatioMin {
//			break
//		}
//	}
//
//	//if mask.inversedBy == 0 {
//	//	return v
//	//} else if v < mask.inversedBy {
//	//	return layer.ValueMax
//	//}
//
//	return 0
//}
