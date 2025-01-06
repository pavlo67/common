package statistics

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"image"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/sources"
)

func FileSeriesTestScenario(t *testing.T, path, testInfoFilename string) {
	//var videoInfo video.Info
	//err := serialization.Read(filepath.Join(path, preparation.VideoInfoFilename), serialization.MarshalerJSON, &videoInfo)
	//require.NoError(t, err)

	testInfoPath := filepath.Join(path, testInfoFilename)
	var testInfo sources.TestInfo
	err := serialization.Read(testInfoPath, serialization.MarshalerJSON, &testInfo)
	require.NoError(t, err)

	var descrs []sources.Description
	err = serialization.ReadAllPartsJSON(filepath.Join(path, sources.FramesAllDescriptionsFilename), &descrs)
	require.NoError(t, err)
	require.True(t, len(descrs) > 0)

	iMin := testInfo.NFrom
	iMax := testInfo.NTo

	iMax = min(iMax, len(descrs))

	// checking ---------------------------------------------------------------

	for i := iMin; i < iMax; i++ {
		descr := descrs[i]
		require.Truef(t, descr.DPM > 0 && !math.IsInf(descr.DPM, 1), "wrong descr.DPMConverted: %d", descr.DPM)

		if i > iMin {
			require.Equalf(t, descrs[i-1].N+1, descr.N, "%d: descr.N (%d) != descrs[i-1].N (%d) + 1", i, descr.N, descrs[i-1].N)
		}
	}

	// processing -------------------------------------------------------------

	var imgGrayPrev *image.Gray

	acc := &Accumulator{
		Path:      testInfoPath,
		Actor:     "file_series_reader",
		Scenario:  filelib.CurrentFile(true),
		StartedAt: time.Now(),
	}

	for i := iMin; i < iMax; i++ {
		descr := descrs[i]

		if i%50 == 0 {
			fmt.Println(i-iMin, " / ", iMax)
		}
		imgGray, err := imagelib.ReadPGMSpecial(filepath.Join(path, fmt.Sprintf("%04d.pgm", descr.N)))
		require.NoError(t, err)
		require.NotNil(t, imgGray)

		acc.Add(float64(imgGray.Rect.Max.X*imgGray.Rect.Max.Y)/float64(Mega), 1, true, 0, 0)

		imgGrayPrev = imgGray
	}

	require.NotNil(t, imgGrayPrev)

	series, err := acc.Series()
	require.NoError(t, err)
	require.NotNil(t, series)

	var timesStr string
	for _, tValues := range series.Times {
		for _, tValue := range tValues {
			timesStr += fmt.Sprintf("%f", tValue) + "\t"
		}
		timesStr += "\n"
	}
	err = os.WriteFile(filepath.Join(path, "times_"+time.Now().Format(time.RFC3339)[:19]+".xls"), []byte(timesStr), 0644)
	require.NoError(t, err)

	series.Times = nil

	err = serialization.SavePart(*series, serialization.MarshalerJSON, filepath.Join(path, StatisticsFilename))
	require.NoError(t, err)
}
