package imagelib

import (
	"fmt"
	"image"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testSubPath = "test_images"

// TODO!!! be careful: imgs [][]image.Image is a list of image rows so it's related to the next indexing: imgs[y][x]

func TestImageCompose(t *testing.T) {
	imgs := make([][]image.Image, 2)
	imgs[0] = make([]image.Image, 2)
	imgs[1] = make([]image.Image, 2)

	var err error

	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			imgs[x][y], err = Read(fmt.Sprintf(filepath.Join(testSubPath, "image%d%d.jpg"), x, y))
			require.NoError(t, err)
			require.NotNil(t, imgs[x][y])
		}
	}

	imgComposed, err := ComposeImages(imgs)
	require.NoError(t, err)
	require.NotNil(t, imgComposed)

	imgComposedExpected, err := Read(filepath.Join(testSubPath, "img_composed.png"))
	require.NoError(t, err)
	require.NotNil(t, imgComposedExpected)

	require.Equal(t, imgComposedExpected.Bounds(), imgComposed.Bounds())

	imgComposedExpectedRGBA, _ := imgComposedExpected.(*image.RGBA)
	require.NotNil(t, imgComposedExpectedRGBA)

	require.Equal(t, imgComposedExpectedRGBA.Pix, imgComposed.Pix)

}
