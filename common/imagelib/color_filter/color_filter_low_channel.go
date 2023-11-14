package color_filter

import (
	"fmt"
	"image/color"

	"github.com/pavlo67/common/common/imagelib/pix"

	"github.com/pavlo67/common/common"
)

var _ Operator = &lowChannelFilter{}

type lowChannelFilter struct {
	ch        int
	ch1, ch2  int
	threshold pix.ValueDelta
}

const onLowChannel = "on color_filter.LowChannel()"

func LowChannel(ch int, threshold pix.ValueDelta) (Operator, error) {
	lowChannel := lowChannelFilter{ch: ch, threshold: threshold}

	// lowChannel.ch1, lowChannel.ch2 = (ch + 1) %3, (ch + 2) %3

	switch ch {
	case 0:
		lowChannel.ch1, lowChannel.ch2 = 1, 2
	case 1:
		lowChannel.ch1, lowChannel.ch2 = 0, 2
	case 2:
		lowChannel.ch1, lowChannel.ch2 = 0, 1
	default:
		return nil, fmt.Errorf("wrong ch = (%d) / "+onLowChannel, ch)
	}

	return &lowChannel, nil
}

func (op lowChannelFilter) Test(rgba color.RGBA) bool {
	rgb := [3]uint8{rgba.R, rgba.G, rgba.B}

	return pix.ValueDelta(rgb[op.ch])-pix.ValueDelta(rgb[op.ch1]) < op.threshold ||
		pix.ValueDelta(rgb[op.ch])-pix.ValueDelta(rgb[op.ch2]) < op.threshold
}

func (op lowChannelFilter) Info() common.Map {
	return common.Map{
		"name":      "low_channel",
		"ch":        op.ch,
		"threshold": op.threshold,
	}
}
