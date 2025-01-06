package statistics

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
)

func TestFullImageOnSeries(t *testing.T) {
	path, err := filelib.GetDir("_test_files/")
	require.NoError(t, err)
	require.NotEmpty(t, path)

	FileSeriesTestScenario(t, path, "files", 4, 6)
}
