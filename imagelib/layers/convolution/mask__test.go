package convolution

import (
	"github.com/pavlo67/common/common/imagelib/pix"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

type RGBAPoint struct {
	X, Y int
	color.RGBA
}

func MaskTestScenario(t *testing.T, mask Mask, rect image.Rectangle, points []RGBAPoint, x, y int, valueExpected pix.Value) {
	require.NotNil(t, mask)

	xWidth, yHeight := rect.Dx(), rect.Dy()
	require.Truef(t, xWidth > 0, "xWidth <= 0 in %v", rect)
	require.Truef(t, yHeight > 0, "yHeight <= 0 in %v", rect)

	imgRGB := image.NewRGBA(rect)
	for _, p := range points {
		imgRGB.Set(p.X, p.Y, p.RGBA)
	}

	err := mask.Prepare(imgRGB)
	require.NoError(t, err)

	info := mask.Info()
	require.NotNil(t, info)

	value := mask.Calculate(x, y)
	require.Equal(t, valueExpected, value)
}
