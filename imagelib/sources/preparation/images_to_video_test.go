package preparation

import (
	"github.com/pavlo67/common/common/serialization"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/imagelib/opencvlib"
	"github.com/pavlo67/imagelib/video"
)

func TestImagesToVideo(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	// data ---------------------------------------------------------------------

	dataPath, err := filelib.GetDir("/home/pavlo/0/partner/_/odometry_data/2/")
	require.NoError(t, err)

	//jlistAllPath := filepath.Join(dataPath, shifter.FramesDescriptionFile)
	//var descrs []sources.Description
	//err = serialization.ReadAllPartsJSON(jlistAllPath, &descrs)
	//require.NoError(t, err)
	//require.Truef(t, len(descrs) > 0, "descrs = []")

	var info video.Info
	err = serialization.Read(filepath.Join(dataPath, video.InfoFilename), serialization.MarshalerJSON, &info)
	require.NoError(t, err)

	divider := 1.
	if info.FPSDivider != nil {
		divider = float64(max(*info.FPSDivider, 1))
	}

	isColor := false

	var re *regexp.Regexp
	if isColor {
		re, err = regexp.Compile(`^\d{4}\.png$`)
	} else {
		re, err = regexp.Compile(`^\d{4}\.pgm$`)
	}

	require.NoError(t, err)
	require.NotNil(t, re)

	err = opencvlib.WriteMP4(filepath.Join(dataPath, filepath.Base(dataPath)+".mp4"), dataPath, *re, false, info.FPS/divider, info.Rectangle.Max.X, info.Rectangle.Max.Y, isColor)
	require.NoError(t, err)

}
