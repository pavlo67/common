package layers

import (
	"fmt"
	"image"
	"math"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/common/common/logger"
)

var _ logger.GetImage = &Layer{}
var _ imagelib.Bounded = &Layer{}

func (lyr *Layer) Image(opts common.Map) (image.Image, string, error) {
	gray, err := lyr.GrayWide()
	return gray, "", err
}

const onGrayWide = "on layers.Layer.GrayWide()"

func (lyr *Layer) GrayWide() (*image.Gray, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil / " + onGrayWide)
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	if xWidth <= 0 || yHeight <= 0 {
		return nil, fmt.Errorf("incorrect lyr.Rect (%#v) / "+onGrayWide, lyr.Rect)
	} else if lyr.Stride < xWidth {
		return nil, fmt.Errorf("lyr.Stride (%d) < lyr.xWidth (%#v) / "+onGrayWide, lyr.Stride, lyr.Rect)
	} else if len(lyr.Pix) < lyr.Stride*(yHeight-1)+xWidth {
		return nil, fmt.Errorf("len(lyr.values) == %d, lyr.Stride = %d, lyr.Rect == %#v / "+onGrayWide, len(lyr.Pix), lyr.Stride, lyr.Rect)
	}

	lyr.MinMax()

	vMin, vMax := lyr.Min, lyr.Max

	if vMax == vMin || (vMin == 0 && vMax == pix.ValueMax) {
		return &lyr.Gray, nil
	}

	gray := image.Gray{
		Pix:    make([]uint8, xWidth*yHeight),
		Stride: xWidth,
		Rect:   lyr.Rect,
	}

	pixValue := float64(pix.ValueMax) / float64(vMax-vMin)

	var stride, strideGray int

	// TODO!!! use original lyr.Pix optionally

	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			gray.Pix[strideGray+x] = uint8(math.Round(float64(lyr.Pix[stride+x]-vMin) * pixValue))
		}
		stride += lyr.Stride
		strideGray += xWidth
	}

	return &gray, nil

}

const onSavePNG = "on layers.Layer.SavePNG()"

func (lyr Layer) SavePNG(filename string) error {
	gray, err := lyr.GrayWide()
	if err != nil {
		return errors.Wrap(err, onSavePNG)
	} else if gray == nil {
		return fmt.Errorf("gray == nil / " + onSavePNG)
	}

	if err = imagelib.SavePNG(gray, filename); err != nil {
		return errors.Wrap(err, onSavePNG)
	}

	return nil
}
