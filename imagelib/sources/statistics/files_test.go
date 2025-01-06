package statistics

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/imagelib/sources"
)

func TestFullImageOnSeries(t *testing.T) {
	if os.Getenv("T") != "" {
		t.Skip()
	}

	// this logFile parameter determines a file for complete logger output duplication
	// envs, l := config.PrepareTests(t, "../_env/", "")
	// require.NotNil(t, envs)

	// data ---------------------------------------------------------------------

	dataPath, err := filelib.GetDir("/home/pavlo/0/partner/_tests/at_10_15")
	require.NoError(t, err)
	require.NotEmpty(t, dataPath)

	FileSeriesTestScenario(t, dataPath, sources.TestInfoFilename)
}
