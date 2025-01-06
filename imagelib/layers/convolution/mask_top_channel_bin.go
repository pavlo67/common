package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	frame2 "github.com/pavlo67/imagelib/frame"
	"github.com/pavlo67/imagelib/layers"
	"image"
	"strconv"

	"github.com/pavlo67/common/common"
)

var _ Mask = &topChannelBinaryMask{}

type topChannelBinaryMask struct {
	imgRGB   *image.RGBA
	ch       int
	ch1, ch2 int
	thr      pix.ValueDelta
	inverse  bool
}

const onTopChannelBinary = "on convolution.TopChannelBinary()"

func TopChannelBinary(ch int, thr pix.ValueDelta, inverse bool) Mask {
	return &topChannelBinaryMask{
		ch:      ch,
		ch1:     (ch + 1) % 3,
		ch2:     (ch + 2) % 3,
		thr:     thr,
		inverse: inverse,
	}
}

func (mask *topChannelBinaryMask) Side() int {
	return 1
}

const onTopChannelBinaryPrepare = "on topChannelBinaryMask.Prepare()"

func (mask *topChannelBinaryMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case image.RGBA:
		mask.imgRGB = &v
	case *image.RGBA:
		mask.imgRGB = v
	case frame2.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case *frame2.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case frame2.Frame:
		mask.imgRGB = &v.RGBA
	case *frame2.Frame:
		mask.imgRGB = &v.RGBA
	}
	if mask.imgRGB == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v)", onData)
		}
		return fmt.Errorf("wrong data (%T)", onData)
	}

	return nil
}

func (mask topChannelBinaryMask) Info() common.Map {
	return common.Map{
		"name": "top_ch_" + strconv.Itoa(int(mask.thr)),
		"thr":  mask.thr,
	}
}

func (mask topChannelBinaryMask) Classes() layers.ClassesCustom {
	return nil
}

func (mask topChannelBinaryMask) Calculate(x, y int) pix.Value {
	imgRGB := mask.imgRGB

	offset := (y-imgRGB.Rect.Min.Y)*imgRGB.Stride + (x-imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
	if offset < 0 || offset >= len(imgRGB.Pix)-coloring.NumColorsRGB {
		return 0
	}

	rgb := imgRGB.Pix[offset : offset+coloring.NumColorsRGB]

	var clrOtherMax uint8
	if rgb[mask.ch1] > rgb[mask.ch2] {
		clrOtherMax = rgb[mask.ch1]
	} else {
		clrOtherMax = rgb[mask.ch2]
	}

	vDelta := pix.ValueDelta(rgb[mask.ch]) - pix.ValueDelta(clrOtherMax) - mask.thr
	if (mask.inverse && vDelta <= 0) || (!mask.inverse && vDelta >= 0) {
		return pix.ValueMax
	}

	return 0
}
