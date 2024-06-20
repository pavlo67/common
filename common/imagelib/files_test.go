package imagelib

import (
	"github.com/pavlo67/common/common/filelib"
	"github.com/stretchr/testify/require"
	"image"
	"os"
	"reflect"
	"testing"
)

func TestSavePGM(t *testing.T) {
	xWidth, yHeigth := 100, 100

	imgGray := image.NewGray(image.Rect(0, 0, xWidth, yHeigth))
	imgGray.Pix[100] = 200

	filename := filelib.CurrentPath() + "test.pgm"

	err := SavePGM(imgGray, filename)
	require.NoError(t, err)

	fi, err := os.Stat(filename)
	require.NoError(t, err)
	require.NotNil(t, fi)

	headerSize := fi.Size() - int64(xWidth*yHeigth)

	require.Equal(t, int64(15), headerSize)

	imgGrayReaded, err := ReadGray(filename)
	require.NoError(t, err)
	require.Equal(t, imgGray.Rect, imgGrayReaded.Rect)
	require.Equal(t, imgGray.Stride, imgGrayReaded.Stride)
	require.True(t, reflect.DeepEqual(imgGray.Pix, imgGrayReaded.Pix))
	require.Equal(t, imgGray.Pix[100], imgGrayReaded.Pix[100])

}
