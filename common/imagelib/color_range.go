package imagelib

import (
	"fmt"
	"image/color"
)

type ColorRange struct {
	ColorMin, ColorMax color.RGBA
}

func (cr ColorRange) String() string {
	return fmt.Sprintf("%d_%d_%d-%d_%d_%d", cr.ColorMin.R, cr.ColorMin.G, cr.ColorMin.B, cr.ColorMax.R, cr.ColorMax.G, cr.ColorMax.B)

}

func CorrectColorsRanges(colorRanges []ColorRange, rangeMax uint8) []ColorRange {
	var colorRangesCorrected []ColorRange

COLOR_RANGE:
	for _, cr := range colorRanges {
		for i, colorRange := range colorRangesCorrected {
			if rangeR := CheckRange(colorRange.ColorMin.R, colorRange.ColorMax.R, cr.ColorMin.R, cr.ColorMax.R, rangeMax); rangeR != nil {
				colorRangesCorrected[i].ColorMin.R, colorRangesCorrected[i].ColorMax.R = rangeR[0], rangeR[1]
				continue COLOR_RANGE
			}
			colorRangesCorrected = append(colorRangesCorrected, cr)
		}
	}

	return colorRangesCorrected
}

func CheckRange(cMin, cMax, cNewMin, cNewMax uint8, rangeMax uint8) []uint8 {
	if cNewMin < cMin {
		cMin = cNewMin

	} else if cNewMax <= cMax {
		return []uint8{cMin, cMax} // returning original {cMin, cMax}

	}

	if cNewMax > cMax {
		cMax = cNewMax
	}

	if cMax-cMin <= rangeMax {
		return []uint8{cMin, cMax}
	}

	return nil
}
