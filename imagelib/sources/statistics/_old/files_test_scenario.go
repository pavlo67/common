package statistics

import (
	"fmt"
	"image"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
)

func FileSeriesTestScenario(t *testing.T, path, key string, iMin, iMax int) {

	//var descrs []sources.Description
	//err := serialization.ReadAllPartsJSON(filepath.Join(path, preparation.FramesAllDescriptionsFilename), &descrs)
	//require.NoError(t, err)
	//require.True(t, len(descrs) > 0)
	//
	//iMax = min(iMax, len(descrs))
	//
	//// checking ---------------------------------------------------------------
	//
	//for i := iMin; i < iMax; i++ {
	//	descr := descrs[i]
	//	require.Truef(t, descr.DPM > 0 && !math.IsInf(descr.DPM, 1), "wrong descr.DPMConverted: %d", descr.DPM)
	//
	//	if i > iMin {
	//		require.Equalf(t, descrs[i-1].N+1, descr.N, "%d: descr.N (%d) != descrs[i-1].N (%d) + 1", i, descr.N, descrs[i-1].N)
	//	}
	//}

	numberedFiles, err := filelib.NumberedFilesSequence(path, `^(\d{4})\.pgm$`, true)
	require.NoError(t, err)
	require.NotNil(t, numberedFiles)

	// processing -------------------------------------------------------------

	var imgGrayPrev *image.Gray

	acc := &Accumulator{
		Key:       key,
		Actor:     "file_series_reader",
		Scenario:  filelib.CurrentFile(true),
		StartedAt: time.Now(),
	}

	for i := iMin; i < iMax; i++ {
		// descr := descrs[i]

		if i%50 == 0 {
			fmt.Println(i-iMin, " / ", iMax)
		}
		imgGray, err := imagelib.ReadPGMSpecial(filepath.Join(path, fmt.Sprintf("%04d.pgm", i))) // descr.N
		require.NoError(t, err)
		require.NotNil(t, imgGray)

		acc.Add(float64(imgGray.Rect.Max.X*imgGray.Rect.Max.Y)/float64(Mega), 1, true, 0, 0)

		imgGrayPrev = imgGray
	}

	require.NotNil(t, imgGrayPrev)

	series, err := acc.Series(nil)
	require.NoError(t, err)
	require.NotNil(t, series)

	err = series.Save(path)
	require.NoError(t, err)
}
