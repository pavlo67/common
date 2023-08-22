package imagelib

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageCompose(t *testing.T) {
	imgs := make([][]image.Image, 2)
	imgs[0] = make([]image.Image, 2)
	imgs[1] = make([]image.Image, 2)

	var err error

	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			imgs[x][y], err = ReadImage(fmt.Sprintf("%d%d.jpg", x, y))
			require.NoError(t, err)
			require.NotNil(t, imgs[x][y])
		}
	}

	imgComposed, err := ComposeImages(imgs)
	require.NoError(t, err)
	require.NotNil(t, imgComposed)

	err = SavePNG(imgComposed, "img_composed.png")
	require.NoError(t, err)
}
