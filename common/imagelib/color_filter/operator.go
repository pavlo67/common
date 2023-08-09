package color_filter

import (
	"image/color"

	"github.com/pavlo67/common/common"
)

type Operator interface {
	Test(rgb color.RGBA) bool
	Info() common.Map
}
