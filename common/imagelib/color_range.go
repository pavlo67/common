package imagelib

import (
	"fmt"
	"image"
	"image/color"

	"github.com/pavlo67/common/common/mathlib/sets"

	"github.com/pavlo67/common/common/imagelib/pix"
)

type ColorRange struct {
	ColorMin, ColorMax color.RGBA
}

func (cr ColorRange) String() string {
	return fmt.Sprintf("%d_%d_%d-%d_%d_%d", cr.ColorMin.R, cr.ColorMin.G, cr.ColorMin.B, cr.ColorMax.R, cr.ColorMax.G, cr.ColorMax.B)

}

func ColorRangesIntersect(cr0, cr1 ColorRange) bool {
	return sets.Intersect(cr0.ColorMin.R, cr0.ColorMax.R, cr1.ColorMin.R, cr1.ColorMax.R) &&
		sets.Intersect(cr0.ColorMin.G, cr0.ColorMax.G, cr1.ColorMin.G, cr1.ColorMax.G) &&
		sets.Intersect(cr0.ColorMin.B, cr0.ColorMax.B, cr1.ColorMin.B, cr1.ColorMax.B)

}

func GetColorRange(imgRGB *image.RGBA, halfRange pix.Value, points ...image.Point) *ColorRange {
	var cr *ColorRange

	var rSum, gSum, bSum, cnt pix.ValueSum

	for _, p := range points {
		if p.X < imgRGB.Rect.Min.X || p.X >= imgRGB.Rect.Max.X || p.Y < imgRGB.Rect.Min.Y || p.Y >= imgRGB.Rect.Max.Y {
			continue
		}
		offset := (p.Y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (p.X-imgRGB.Rect.Min.X)*NumColorsRGBA
		clr := imgRGB.Pix[offset : offset+3]

		rSum += pix.ValueSum(clr[0])
		gSum += pix.ValueSum(clr[1])
		bSum += pix.ValueSum(clr[2])
		cnt++

		//if cr == nil {
		//	cr = &ColorRange{
		//		ColorMin: color.RGBA{R: clr[0], G: clr[1], B: clr[2], A: 255},
		//		ColorMax: color.RGBA{R: clr[0], G: clr[1], B: clr[2], A: 255},
		//	}
		//
		//} else {
		//	if clr[0] < cr.ColorMin.R {
		//		cr.ColorMin.R = clr[0]
		//	} else if clr[0] > cr.ColorMax.R {
		//		cr.ColorMax.R = clr[0]
		//	}
		//	if clr[1] < cr.ColorMin.G {
		//		cr.ColorMin.G = clr[1]
		//	} else if clr[1] > cr.ColorMax.G {
		//		cr.ColorMax.G = clr[1]
		//	}
		//	if clr[2] < cr.ColorMin.B {
		//		cr.ColorMin.B = clr[2]
		//	} else if clr[2] > cr.ColorMax.B {
		//		cr.ColorMax.B = clr[2]
		//	}
		//}

		offset += NumColorsRGBA
	}

	if cnt > 0 {
		rAvg, gAvg, bAvg := rSum/cnt, gSum/cnt, bSum/cnt
		return &ColorRange{
			ColorMin: color.RGBA{
				R: pix.Value(max(0, rAvg-pix.ValueSum(halfRange))),
				G: pix.Value(max(0, gAvg-pix.ValueSum(halfRange))),
				B: pix.Value(max(0, bAvg-pix.ValueSum(halfRange))),
				A: 255,
			},
			ColorMax: color.RGBA{
				R: pix.Value(min(255, rAvg+pix.ValueSum(halfRange))),
				G: pix.Value(min(255, gAvg+pix.ValueSum(halfRange))),
				B: pix.Value(min(255, bAvg+pix.ValueSum(halfRange))),
				A: 0,
			},
		}
	}

	return cr
}

func CorrectColorsRanges(colorRanges []ColorRange, rangeMax uint8) []ColorRange {

	if len(colorRanges) < 1 {
		return nil
	}
	colorRangesCorrected := []ColorRange{colorRanges[0]}

COLOR_RANGE:
	for _, cr := range colorRanges[1:] {
		for i, colorRange := range colorRangesCorrected {
			var rangeR, rangeG, rangeB []uint8
			if rangeR = CheckRange(colorRange.ColorMin.R, colorRange.ColorMax.R, cr.ColorMin.R, cr.ColorMax.R, rangeMax); rangeR == nil {
				break
			}
			if rangeG = CheckRange(colorRange.ColorMin.G, colorRange.ColorMax.G, cr.ColorMin.G, cr.ColorMax.G, rangeMax); rangeG == nil {
				break
			}
			if rangeB = CheckRange(colorRange.ColorMin.B, colorRange.ColorMax.B, cr.ColorMin.B, cr.ColorMax.B, rangeMax); rangeB == nil {
				break
			}

			colorRangesCorrected[i].ColorMin.R, colorRangesCorrected[i].ColorMax.R = rangeR[0], rangeR[1]
			colorRangesCorrected[i].ColorMin.G, colorRangesCorrected[i].ColorMax.G = rangeG[0], rangeG[1]
			colorRangesCorrected[i].ColorMin.B, colorRangesCorrected[i].ColorMax.B = rangeB[0], rangeB[1]
			continue COLOR_RANGE
		}

		colorRangesCorrected = append(colorRangesCorrected, cr)
	}

	return colorRangesCorrected
}

func CheckRange(cMin, cMax, cNewMin, cNewMax uint8, rangeMax uint8) []uint8 {
	connRange := sets.Connection(cMin, cMax, cNewMin, cNewMax)
	if connRange != nil && connRange[1]-connRange[0] <= rangeMax {
		return []uint8{connRange[0], connRange[1]}
	}

	return nil
}
