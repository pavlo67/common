package preparation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/sources"
)

func TestFramesJListToFramesAllJList(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	// this logFile parameter determines a file for complete logger output duplication
	//envs, _ := config.PrepareTests(t, "../_env/", "")
	//require.NotNil(t, envs)

	// data ---------------------------------------------------------------------

	dataPath, err := filelib.GetDir("_data/odometry/at_10_15")
	require.NoError(t, err)

	jlistPath := filepath.Join(dataPath, sources.FramesDescriptionsFilename)

	var descrs []sources.Description
	err = serialization.ReadAllPartsJSON(jlistPath, &descrs)
	require.NoError(t, err)

	require.Truef(t, len(descrs) > 0, "descrs = []")

	descrsAll := []sources.Description{descrs[0]}

	// t.Fatalf("%+v", descrs[0])

	for i, descr := range descrs[:len(descrs)-1] {
		descrNext := descrs[i+1]
		require.Truef(t, descrNext.N > descr.N, "non-sequental numeration: descr.N = %d, descrNext.N = %d", descr.N, descrNext.N)
		require.NotNilf(t, descrNext.GeoPoint, "%d: descr.GeoPoint == nil", descrNext.N)

		descrsAll = append(descrsAll, sources.InterpolatedDescriptions(descr, descrNext)...)
		descrsAll = append(descrsAll, descrs[i+1])
	}

	jlistAllPath := filepath.Join(dataPath, sources.FramesAllDescriptionsFilename)

	err = serialization.SaveAllPartsJSON(descrsAll, jlistAllPath)
	require.NoError(t, err)
}
