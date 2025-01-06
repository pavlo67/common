package images

import (
	"github.com/pavlo67/common/common/geolib"
	"github.com/pavlo67/imagelib/sources"
	"image"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/db"
)

const testKey = "0001"

func OperatorTestScenario(t *testing.T, imagesOp Operator, imagesCleanerOp db.Cleaner) {
	require.Equal(t, "test", os.Getenv("ENV"))

	require.NotNil(t, imagesOp)
	require.NotNil(t, imagesCleanerOp)

	err := imagesCleanerOp.Clean()
	require.NoError(t, err)

	testImage := image.NewRGBA(image.Rect(0, 0, 10, 10))
	testImage.Set(5, 5, colornames.White)
	testImage.Set(7, 7, colornames.Orange)

	testDescr := sources.Description{
		N: 1,
		//ImageRef: sources.ImageRef{
		//	ImagePath: "11",
		//	SourceKey: "33",
		//},
		GeoPoint: &geolib.Point{1, 2},
		Bearing:  3,
		DPM:      2,
	}

	testImageChecked, err := imagesOp.Check(testKey)
	require.NoError(t, err)
	require.False(t, testImageChecked)

	imgPath, err := imagesOp.Save(testImage, testDescr, testKey)
	require.NoError(t, err)
	require.NotEmpty(t, imgPath)
	testDescr.ImagePath = imgPath

	testImageChecked, err = imagesOp.Check(testKey)
	require.NoError(t, err)
	require.True(t, testImageChecked)

	testImageReaded, testDescrReaded, err := imagesOp.Get(testKey)
	require.NoError(t, err)
	require.NotNil(t, testImageReaded)
	require.NotNil(t, testDescrReaded)

	require.Truef(t, reflect.DeepEqual(testDescr, *testDescrReaded), "testDescr: %+v, testDescrReaded: %+v", testDescr, testDescrReaded)

	rect := testImage.Bounds()
	require.Equal(t, rect, testImageReaded.Bounds())

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			rExpected, gExpected, bExpected, aExpected := testImage.At(x, y).RGBA()
			r, g, b, a := testImageReaded.At(x, y).RGBA()

			require.Equal(t, rExpected, r)
			require.Equal(t, gExpected, g)
			require.Equal(t, bExpected, b)
			require.Equal(t, aExpected, a)
		}
	}

	imgPaths, err := imagesOp.ListPaths(testKey)
	require.NoError(t, err)
	require.True(t, len(imgPaths) == 1)

}
