package convolution

import (
	"fmt"
	"image"

	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

type Mask interface {
	Prepare(onData interface{}) error
	Calculate(x, y int) pix.Value // (x, y) corresponds to the left-top (not center!) of calculation area
	Info() common.Map
	Classes() layers.ClassesCustom
}

const onLayer = "on convolution.Layer()"

func Layer(data imagelib.Described, mask Mask, scale int, addRest bool) (*layers.Layer, error) {

	rect := data.Bounds()
	xWidth0, yHeight0 := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y

	if xWidth0 <= 0 || yHeight0 <= 0 {
		return nil, fmt.Errorf("incorrect data.Rect (%#v) / "+onLayer, rect)
	} else if scale <= 0 {
		return nil, fmt.Errorf("incorrect sideTop (%d) / "+onLayer, scale)
	} else if err := mask.Prepare(data); err != nil {
		return nil, errors.Wrap(err, onLayer)
	}

	xWidth, yHeight := xWidth0/scale, yHeight0/scale
	if addRest {
		if xWidth0%scale != 0 {
			xWidth++
		}
		if yHeight0%scale != 0 {
			yHeight++
		}
	}

	settings := data.Description()
	settings.DPM = settings.DPM / float64(scale)

	lyrConvolved := layers.Layer{
		Gray: image.Gray{
			Pix:    make([]pix.Value, xWidth*yHeight),
			Stride: xWidth,

			// TODO! be careful: .Rect.Max looks oddly if rect.Min != {0,0} && convolve && sideTop > 1
			Rect: image.Rectangle{rect.Min, rect.Min.Add(image.Point{xWidth, yHeight})},
		},
		Settings: settings,
	}

	// TODO??? save previous options values with the same k
	for k, v := range mask.Info() {
		lyrConvolved.SetOptions(k, v)
	}

	minValue := mask.Calculate(rect.Min.X, rect.Min.Y)
	maxValue := minValue

	var offset, whiteCnt, blackCnt int

	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			v := mask.Calculate(rect.Min.X+x*scale, rect.Min.Y+y*scale)
			lyrConvolved.Pix[offset+x] = v
			if v == pix.ValueMax {
				whiteCnt++
				maxValue = v
			} else if v > maxValue {
				maxValue = v
			} else if v == 0 {
				blackCnt++
				minValue = v
			} else if v < minValue {
				minValue = v
			}
		}
		offset += xWidth
	}

	lyrConvolved.Min, lyrConvolved.Max = minValue, maxValue
	//if cnt := xWidth * yHeight; cnt > 0 {
	//	lyrConvolved.WhRat, lyrConvolved.BlRat = float64(whiteCnt)/float64(cnt), float64(blackCnt)/float64(cnt)
	//}
	lyrConvolved.ClassesCustom = mask.Classes()

	return &lyrConvolved, nil
}

func Metrics(data imagelib.Described, mask Mask, scale int, addRest bool) (*layers.Metrics, error) {
	rect := data.Bounds()
	xWidth0, yHeight0 := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y

	if xWidth0 <= 0 || yHeight0 <= 0 {
		return nil, fmt.Errorf("incorrect data.Rect (%#v) / "+onLayer, rect)
	} else if scale <= 0 {
		return nil, fmt.Errorf("incorrect sideTop (%d) / "+onLayer, scale)
	} else if err := mask.Prepare(data); err != nil {
		return nil, errors.Wrap(err, onLayer)
	}

	xWidth, yHeight := xWidth0/scale, yHeight0/scale
	if addRest {
		if xWidth0%scale != 0 {
			xWidth++
		}
		if yHeight0%scale != 0 {
			yHeight++
		}
	}

	minValue := mask.Calculate(rect.Min.X, rect.Min.Y)
	maxValue := minValue

	var offset, whiteCnt, blackCnt int

	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			v := mask.Calculate(rect.Min.X+x*scale, rect.Min.Y+y*scale)
			if v == pix.ValueMax {
				whiteCnt++
				maxValue = v
			} else if v > maxValue {
				maxValue = v
			} else if v == 0 {
				blackCnt++
				minValue = v
			} else if v < minValue {
				minValue = v
			}

		}
		offset += xWidth
	}

	metrics := layers.Metrics{}
	metrics.Min, metrics.Max = minValue, maxValue
	if cnt := xWidth * yHeight; cnt > 0 {
		//metrics.WhRat, metrics.BlRat = float64(whiteCnt)/float64(cnt), float64(blackCnt)/float64(cnt)
	}

	return &metrics, nil
}
