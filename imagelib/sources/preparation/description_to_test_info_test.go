package preparation

import (
	"github.com/pavlo67/imagelib/sources"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/serialization"
)

func TestDescriptionToTestInfo(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	// data ---------------------------------------------------------------------

	dataPath, err := filelib.GetDir("/home/pavlo/0/partner/_tests/at_10_15")
	require.NoError(t, err)

	testInfoPath := filepath.Join(dataPath, sources.TestInfoFilename)

	testInfo := sources.TestInfo{
		NFrom: 50,
		NTo:   900,
	}

	err = serialization.Save(testInfo, serialization.MarshalerJSON, testInfoPath)
	require.NoError(t, err)
}
