package convolution_rgb

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/frame"
	"image"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

type Mask interface {
	Prepare(onData interface{}) error
	Calculate(x, y int) frame.ValueRGBA
	Info() common.Map
}

const onLayer = "on convolution_rgb.Layer()"

func Layer(data imagelib.Bounded, dpm float64, mask Mask, step int, addRest bool) (*frame.LayerRGBA, error) {

	if step < 0 {
		return nil, fmt.Errorf("incorrect step (%d) / "+onLayer, step)
	} else if err := mask.Prepare(data); err != nil {
		return nil, errors.Wrap(err, onLayer)
	}

	rect := data.Bounds()
	xWidth0, yHeight0 := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y

	if xWidth0 <= 0 || yHeight0 <= 0 {
		return nil, fmt.Errorf("incorrect data.Rect (%#v) / "+onLayer, rect)
	} else if step <= 0 {
		return nil, fmt.Errorf("step (%d) <= 0 / "+onLayer, step)
	}

	xWidth, yHeight := xWidth0/step, yHeight0/step

	if addRest {
		if xWidth0%step != 0 {
			xWidth++
		}
		if yHeight0%step != 0 {
			yHeight++
		}
	}

	lyrConvolved := frame.LayerRGBA{
		RGBA: image.RGBA{
			Pix:    make([]pix.Value, xWidth*yHeight*coloring.NumColorsRGBA),
			Stride: xWidth * coloring.NumColorsRGBA,
			// TODO! be careful: .Rect.Max looks oddly if rect.Min != {0,0} and step > 1
			Rect: image.Rectangle{rect.Min, image.Point{rect.Min.X + xWidth, rect.Min.Y + yHeight}},
		},
		Settings: imagelib.Settings{
			DPM: dpm / float64(step),
		},
	}

	for k, v := range mask.Info() {
		lyrConvolved.SetOptions(k, v)
	}

	var offset int
	for y := 0; y < yHeight; y++ {
		offsetX := offset
		for x := 0; x < xWidth; x++ {
			v := mask.Calculate(rect.Min.X+x*step, rect.Min.Y+y*step)
			lyrConvolved.Pix[offsetX] = v[0]
			lyrConvolved.Pix[offsetX+1] = v[1]
			lyrConvolved.Pix[offsetX+2] = v[2]
			lyrConvolved.Pix[offsetX+3] = v[3]
			offsetX += coloring.NumColorsRGBA
		}
		offset += lyrConvolved.Stride
	}

	return &lyrConvolved, nil
}
