package preparation

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/imagelib/sources"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
)

func TestImages(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	imagesPath, err := filelib.GetDir("/home/pavlo/0/partner/_tests/arcgis_0/")
	require.NoError(t, err)

	resultPath, err := filelib.GetDir("/home/pavlo/0/partner/_/odometry_data/2/")
	require.NoError(t, err)

	numberedFiles, err := filelib.NumberedFilesSequence(imagesPath, sources.RePGNStr, false)
	require.NoError(t, err)

	for i, nf := range numberedFiles {

		if i%50 == 0 {
			fmt.Println(i, " / ", len(numberedFiles))
		}

		//imgGray, err := imagelib.ReadGray(filepath.Join(imagesPath, fmt.Sprintf("%04d.pgm", descr.N)))

		filename := filepath.Join(imagesPath, fmt.Sprintf("%04d.pgm", nf.I))

		img, err := imagelib.ReadPGMSpecial(filename)
		require.NoError(t, err)
		require.NotNil(t, img)

		require.Truef(t, img.Rect.Max.X >= 998 && img.Rect.Max.Y >= 798, "%d: %+v", nf.I, img.Rect)

		imgCutted := img.SubImage(image.Rectangle{Max: image.Point{998, 798}})

		err = imagelib.SavePGM(imgCutted, filepath.Join(resultPath, fmt.Sprintf("%04d.pgm", nf.I)))
		require.NoError(t, err)

	}

}
