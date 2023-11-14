package color_filter

import (
	"fmt"
	"image/color"

	"github.com/pavlo67/common/common/imagelib/pix"

	"github.com/pavlo67/common/common"
)

var _ Operator = &topChannelFilter{}

type topChannelFilter struct {
	ch        int
	ch1, ch2  int
	threshold pix.ValueDelta
}

const onTopChannel = "on color_filter.TopChannel()"

func TopChannel(ch int, threshold pix.ValueDelta) (Operator, error) {
	topChannel := topChannelFilter{ch: ch, threshold: threshold}

	// topChannel.ch1, topChannel.ch2 = (ch + 1) %3, (ch + 2) %3

	switch ch {
	case 0:
		topChannel.ch1, topChannel.ch2 = 1, 2
	case 1:
		topChannel.ch1, topChannel.ch2 = 0, 2
	case 2:
		topChannel.ch1, topChannel.ch2 = 0, 1
	default:
		return nil, fmt.Errorf("wrong ch = (%d) / "+onTopChannel, ch)
	}

	return &topChannel, nil
}

func (op topChannelFilter) Test(rgba color.RGBA) bool {
	rgb := [3]uint8{rgba.R, rgba.G, rgba.B}

	return pix.ValueDelta(rgb[op.ch])-pix.ValueDelta(rgb[op.ch1]) > op.threshold &&
		pix.ValueDelta(rgb[op.ch])-pix.ValueDelta(rgb[op.ch2]) > op.threshold
}

func (op topChannelFilter) Info() common.Map {
	return common.Map{
		"name":      "top_channel",
		"ch":        op.ch,
		"threshold": op.threshold,
	}
}
