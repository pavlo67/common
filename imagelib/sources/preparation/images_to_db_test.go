package preparation

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
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

func TestImagesToTestDB(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	imagesPath := "/home/pavlo/0/partner/_/"
	dpmRequired := 2.

	processingPath, err := filelib.Dir(filepath.Join(imagesPath, time.Now().Format(time.RFC3339)[:19]))

	numberedFiles, err := filelib.NumberedFilesSequence(imagesPath, sources.RePGNPNGStr, false)
	require.NoError(t, err)

	var descrs []sources.Description
	err = serialization.ReadAllPartsJSON(filepath.Join(imagesPath, sources.FramesDescriptionsFilename), &descrs)
	require.NoError(t, err)

	if len(numberedFiles) < 1 {
		t.Fatalf("no files found in %s (%s)", imagesPath, sources.RePGNPNGStr)
	}

	if len(descrs) < 1 {
		t.Fatalf("no descriptions found in %s (%s)", imagesPath, sources.FramesDescriptionsFilename)
	}

	require.Equal(t, numberedFiles[0].I, descrs[0].N)
	require.Equal(t, numberedFiles[len(numberedFiles)-1].I, descrs[len(descrs)-1].N)

	descrsAll := []sources.Description{descrs[0]}

	require.NotNilf(t, descrs[0].GeoPoint, "N = %d, descr.GeoPoint == nil", descrs[0].N)
	for i, descr := range descrs[:len(descrs)-1] {
		descrNext := descrs[i+1]
		require.Truef(t, descrNext.N > descr.N, "non-sequental numeration: descrs[%d].N = %d, descrs[%d].N = %d", i, descr.N, i+1, descrNext.N)
		require.NotNilf(t, descrNext.GeoPoint, "N = %d, descr.GeoPoint == nil", descrNext.N)
		descrsAll = append(descrsAll, sources.InterpolatedDescriptions(descr, descrNext)...)
		descrsAll = append(descrsAll, descrs[i+1])
	}

	require.Equal(t, len(numberedFiles), len(descrsAll))
	require.Equal(t, numberedFiles[0].I, descrsAll[0].N)
	require.Truef(t, descrs[0].DPM > 0 && !math.IsInf(descrs[0].DPM, 1), "N = %d, wrong descr.DPM: %f", descrs[0].N, descrs[0].DPM)

	for i, descr := range descrsAll[:len(descrsAll)-1] {
		descrNext := descrsAll[i+1]
		require.Truef(t, descrNext.N == descr.N+1, "non-sequental numeration: descrsAll[%d].N = %d, descrsAll[%d].N = %d", i, descr.N, i+1, descrNext.N)
		require.Truef(t, descr.DPM > 0 && !math.IsInf(descr.DPM, 1), "N = %d, wrong descr.DPM: %f", descr.N, descr.DPM)

	}

	series := sources.Series{
		Path:        processingPath,
		DPMRequired: &dpmRequired,
	}
	err = serialization.Read(filepath.Join(imagesPath, sources.VideoInfoFilename), serialization.MarshalerJSON, &series.Info)
	require.NoError(t, err)

	for i, descr := range descrsAll {
		if i%10 == 0 {
			fmt.Println(i, " / ", len(descrsAll))
		}

		filename := filepath.Join(imagesPath, fmt.Sprintf("%04d.pgm", descr.N))
		if ok, _ := filelib.FileExists(filename, false); !ok {
			filename = filepath.Join(imagesPath, fmt.Sprintf("%04d.png", descr.N))
		}

		img, err := imagelib.Read(filename)
		require.NoError(t, err)
		require.NotNil(t, img)

		err = series.Add(img, descr)
		require.NoError(t, err)
	}

	err = series.Save()
	require.NoError(t, err)
}
