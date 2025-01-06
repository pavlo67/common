package convolution

import (
	"github.com/pavlo67/imagelib/layers"
	"math"
)

func SeparationRatio(lyr *layers.Layer, blackMax, whiteMin int) (blackRatio, whiteRatio float64) {
	if lyr == nil {
		return math.NaN(), math.NaN()
	}

	xWidth1, yHeight1 := lyr.Rect.Max.X-lyr.Rect.Min.X-1, lyr.Rect.Max.Y-lyr.Rect.Min.Y-1

	var blackCnt, whiteCnt int64

	offset := lyr.Stride
	for y := 1; y < yHeight1; y++ {
		for x := 1; x < xWidth1; x++ {
			cnt := 0
			// xyIsWhite := false
			if lyr.Pix[offset+x-lyr.Stride-1] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x-lyr.Stride] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x-lyr.Stride+1] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x-1] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x] != 0 {
				// xyIsWhite = true
				cnt++
			}
			if lyr.Pix[offset+x+1] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x+lyr.Stride-1] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x+lyr.Stride] != 0 {
				cnt++
			}
			if lyr.Pix[offset+x+lyr.Stride+1] != 0 {
				cnt++
			}

			if cnt >= whiteMin { // && xyIsWhite
				whiteCnt++
			} else if cnt <= blackMax { // && !xyIsWhite
				blackCnt++
			}
		}

		offset += lyr.Stride
	}

	total := float64((xWidth1 + 2) * (yHeight1 + 2))

	return float64(blackCnt) / total, float64(whiteCnt) / total
}
