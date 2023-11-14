package imagelib

import (
	"fmt"
	"image/color"

	"github.com/pavlo67/common/common/mathlib/combinatorics"
)

type ColorRange struct {
	ColorMin, ColorMax color.RGBA
}

func (cr ColorRange) String() string {
	return fmt.Sprintf("%d_%d_%d-%d_%d_%d", cr.ColorMin.R, cr.ColorMin.G, cr.ColorMin.B, cr.ColorMax.R, cr.ColorMax.G, cr.ColorMax.B)

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
	connRange := combinatorics.Connection(cMin, cMax, cNewMin, cNewMax)
	if connRange != nil && connRange[1]-connRange[0] <= rangeMax {
		return []uint8{connRange[0], connRange[1]}
	}

	return nil
}
