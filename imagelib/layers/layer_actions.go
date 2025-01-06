package layers

import (
	"fmt"
	"image"
	"math"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/common/common/mathlib/plane"
)

func (lyr *Layer) SubLayer(r image.Rectangle) *Layer {
	if lyr == nil {
		return nil // &img{}
	}
	r = r.Intersect(lyr.Rect)
	if r.Empty() {
		return nil // &img{}
	}
	i := lyr.Offset(r.Min.X, r.Min.Y)
	subLayer := Layer{
		Gray: image.Gray{
			Pix:    lyr.Pix[i:],
			Stride: lyr.Stride,
			Rect:   r,
		},
		Settings: lyr.Settings.Copy(),
	}

	subLayer.MinMax()

	return &subLayer
}

func (lyr *Layer) ThresholdValue(threshold float64) pix.Value {
	if lyr == nil {
		return 0
	}

	return lyr.Min + pix.Value(math.Round(float64(lyr.Max-lyr.Min)*threshold))
}

const onThresholded = "on layer.Thresholded()"

func (lyr *Layer) Thresholded(thr pix.Value, inverse bool) (*Layer, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onThresholded)
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	if xWidth <= 0 || yHeight <= 0 {
		return nil, fmt.Errorf("incorrect lyr.Rect (%#v) / "+onThresholded, lyr.Rect)
	} else if lyr.Stride < xWidth {
		return nil, fmt.Errorf("lyr,Stride (%d) < lyr.xWidth (%#v) / "+onThresholded, lyr.Stride, lyr.Rect)
	} else if len(lyr.Pix) < lyr.Stride*(yHeight-1)+xWidth {
		return nil, fmt.Errorf("len(lyr.values) == %d, lyr.Stride = %d, lyr.Rect == %#v / "+onThresholded, len(lyr.Pix), lyr.Stride, lyr.Rect)
	}

	lenValues := xWidth * yHeight

	lyrThresholded := Layer{
		Gray: image.Gray{
			Pix:    make([]pix.Value, lenValues),
			Stride: xWidth,
			Rect:   lyr.Rect,
		},
		Settings: lyr.Settings.Copy(),
	}
	lyrThresholded.Settings.SetOptions("thr", thr)
	lyrThresholded.Settings.SetOptions("inversed", inverse)

	var whiteCount int

	var stride, strideThresholded int
	if inverse {
		for y := 0; y < yHeight; y++ {
			for x := 0; x < xWidth; x++ {
				if lyr.Pix[stride+x] <= thr {
					lyrThresholded.Pix[strideThresholded+x] = pix.ValueMax
					whiteCount++
				} else {
					lyrThresholded.Pix[strideThresholded+x] = 0
				}
			}
			stride += lyr.Stride
			strideThresholded += xWidth
		}
	} else {
		for y := 0; y < yHeight; y++ {
			for x := 0; x < xWidth; x++ {
				if lyr.Pix[stride+x] >= thr {
					lyrThresholded.Pix[strideThresholded+x] = pix.ValueMax
					whiteCount++
				} else {
					lyrThresholded.Pix[strideThresholded+x] = 0
				}
			}
			stride += lyr.Stride
			strideThresholded += xWidth
		}
	}

	if whiteCount == lenValues {
		lyrThresholded.Min, lyrThresholded.Max = pix.ValueMax, pix.ValueMax
	} else if whiteCount == 0 {
		lyrThresholded.Min, lyrThresholded.Max = 0, 0
	} else {
		lyrThresholded.Min, lyrThresholded.Max = 0, pix.ValueMax
	}

	//lyrThresholded.WhRat = float64(whiteCount) / float64(lenValues)
	//lyrThresholded.BlRat = 1 - lyrThresholded.Metrics.WhRat

	return &lyrThresholded, nil
}

const onCenter = "on Layer.Center()"

func (lyr *Layer) Center(thr pix.Value) (*plane.Point2, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onCenter)
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	if xWidth <= 0 || yHeight <= 0 {
		return nil, fmt.Errorf("incorrect lyr.Rect (%#v) / "+onCenter, lyr.Rect)
	} else if lyr.Stride < xWidth {
		return nil, fmt.Errorf("lyr,Stride (%d) < lyr.xWidth (%#v) / "+onCenter, lyr.Stride, lyr.Rect)
	} else if len(lyr.Pix) < lyr.Stride*(yHeight-1)+xWidth {
		return nil, fmt.Errorf("len(lyr.values) == %d, lyr.Stride = %d, lyr.Rect == %#v / "+onCenter, len(lyr.Pix), lyr.Stride, lyr.Rect)
	}

	var offset int

	// TODO!!! be careful with a (practical impossible!) overflow
	var xSum, ySum, cnt int64

	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			if lyr.Pix[offset+x] >= thr {
				xSum += int64(x)
				ySum += int64(y)
				cnt++
			}
		}
		offset += lyr.Stride
	}

	// log.Printf("lyr.Rect: %v, xSum: %d, ySum: %d, cnt: %d", lyr.Rect, xSum, ySum, cnt)

	if cnt > 0 {
		return &plane.Point2{
			float64(lyr.Rect.Min.X) + float64(xSum)/float64(cnt),
			float64(lyr.Rect.Min.Y) + float64(ySum)/float64(cnt),
		}, nil
	}

	return nil, nil
}

const onRectThr = "on Layer.RectThr()"

func (lyr *Layer) RectThr(thr pix.Value) (*image.Rectangle, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onRectThr)
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	if xWidth <= 0 || yHeight <= 0 {
		return nil, fmt.Errorf("incorrect lyr.Rect (%#v) / "+onRectThr, lyr.Rect)
	} else if lyr.Stride < xWidth {
		return nil, fmt.Errorf("lyr,Stride (%d) < lyr.xWidth (%#v) / "+onRectThr, lyr.Stride, lyr.Rect)
	} else if len(lyr.Pix) < lyr.Stride*(yHeight-1)+xWidth {
		return nil, fmt.Errorf("len(lyr.values) == %d, lyr.Stride = %d, lyr.Rect == %#v / "+onRectThr, len(lyr.Pix), lyr.Stride, lyr.Rect)
	}

	var offset, xMin, xMax, yMin, yMax int

	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			if lyr.Pix[offset+x] >= thr {
				if xMin == xMax {
					xMin, xMax, yMin, yMax = x, x+1, y, y+1
				} else {
					if x >= xMax {
						xMax = x + 1
					} else if x < xMin {
						xMin = x
					}
					if y >= yMax {
						yMax = y + 1
					} else if y < yMin {
						yMin = y
					}
				}
			}
		}
		offset += lyr.Stride
	}

	if xMin == xMax {
		return nil, nil
	}

	return &image.Rectangle{image.Point{xMin, yMin}.Add(lyr.Rect.Min), image.Point{xMax, yMax}.Add(lyr.Rect.Min)}, nil
}

const onInversed = "on layer.Inversed()"

func (lyr *Layer) Inversed() (*Layer, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onInversed)
	}

	lyrInversed := Layer{
		Gray: image.Gray{
			Pix:    make([]pix.Value, len(lyr.Pix)),
			Stride: lyr.Stride,
			Rect:   lyr.Gray.Rect,
		},
		Settings: lyr.Settings.Copy(),
	}
	lyrInversed.Settings.SetOptions("inversed", true)

	for i, v := range lyr.Pix {
		lyrInversed.Pix[i] = pix.ValueMax - v
	}
	//lyrInversed.WhRat, lyrInversed.BlRat = lyr.Metrics.BlRat, lyr.Metrics.WhRat
	lyrInversed.Min, lyrInversed.Max = pix.ValueMax-lyrInversed.Max, pix.ValueMax-lyrInversed.Min

	return &lyrInversed, nil
}

const onTransposed = "on lyr.BrClassesTr()"

func (lyr *Layer) Transposed() (*Layer, error) {

	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onTransposed)
	}

	xWidth0, yHeight0 := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	if lyr.Stride < xWidth0 {
		return nil, fmt.Errorf("lyr,Stride (%d) < lyr.xWidth0 (%#v) / "+onTransposed, lyr.Stride, lyr.Rect)
	} else if len(lyr.Pix) < lyr.Stride*(yHeight0-1)+xWidth0 {
		return nil, fmt.Errorf("len(lyr.values) == %d, lyr.Stride = %d, lyr.Rect == %#v / "+onTransposed, len(lyr.Pix), lyr.Stride, lyr.Rect)
	}

	xWidth, yHeight := yHeight0, xWidth0

	layerTransposed := Layer{
		Gray: image.Gray{
			Pix:    make([]pix.Value, xWidth*yHeight),
			Stride: xWidth,
			Rect:   image.Rectangle{Min: image.Point{lyr.Rect.Min.Y, lyr.Rect.Min.X}, Max: image.Point{lyr.Rect.Max.Y, lyr.Rect.Max.X}},
		},
		Settings: lyr.Settings.Copy(),
		Metrics:  lyr.Metrics,
	}

	var offset int
	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			layerTransposed.Pix[offset+x] = lyr.Pix[xWidth0*x+y]
		}
		offset += xWidth
	}

	return &layerTransposed, nil
}

func (lyr *Layer) Summa(xMin, yMin, xMax, yMax int) (pix.ValueSum, int) {
	if lyr == nil {
		return 0, 0
	}

	// log.Print(xMin, xMax, yMin, yMax)
	// log.Printf("%#v", lyr.Rect)

	if xMin < lyr.Rect.Min.X {
		xMin = lyr.Rect.Min.X
	}
	if xMax > lyr.Rect.Max.X {
		xMax = lyr.Rect.Max.X
	}
	if yMin < lyr.Rect.Min.Y {
		yMin = lyr.Rect.Min.Y
	}
	if yMax > lyr.Rect.Max.Y {
		yMax = lyr.Rect.Max.Y
	}

	// log.Print(xMin, xMax, yMin, yMax)

	sum, cnt, offset := pix.ValueSum(0), 0, (yMin-lyr.Rect.Min.Y)*lyr.Stride-lyr.Rect.Min.X
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			sum += pix.ValueSum(lyr.Pix[offset+x])
		}
		cnt += xMax - xMin
		offset += lyr.Stride
	}

	return sum, cnt
}

const onZones = "on lyr.Zones()"

func (lyr Layer) Zones(halfSide float64) ([]Layer, error) {

	if lyr.Settings.DPM == 0 {
		return nil, fmt.Errorf("ls.DPM == 0 / " + onZones)
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	halfSidePix := int(math.Round(halfSide * lyr.Settings.DPM))

	var zones []Layer
	var stride int

	for yInt := 0; yInt < yHeight; yInt++ {
		xInt := 0
		for xInt < xWidth {
			if lyr.Pix[stride+xInt] > 0 {
				inZones := In(zones, image.Point{lyr.Rect.Min.X + xInt, lyr.Rect.Min.Y + yInt})
				if inZones > 0 {
					xInt += inZones
				} else {
					r := image.Rect(xInt-halfSidePix, yInt-halfSidePix, xInt+halfSidePix+1, yInt+halfSidePix+1)
					if zone := lyr.SubLayer(r); zone != nil {
						zones = append(zones, *zone)
						xInt += halfSidePix + 1
					} else {
						return nil, fmt.Errorf("lyr.SubLayer(%v) == nil / "+onZones, r)
					}

				}
				continue
			} else {
				xInt++
			}
		}
		stride += lyr.Stride
	}

	return zones, nil
}

func In(zones []Layer, p image.Point) int {
	for _, zone := range zones {
		if p.In(zone.Rect) {
			return zone.Rect.Max.X - p.X
		}
	}

	return 0
}
