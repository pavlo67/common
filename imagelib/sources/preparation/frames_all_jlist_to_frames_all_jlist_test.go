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

func TestFramesAllJListToFramesAllJList(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	// data ---------------------------------------------------------------------

	dataPath, err := filelib.GetDir("/home/pavlo/0/partner/_tests/at_10_15")
	require.NoError(t, err)

	jlistPath := filepath.Join(dataPath, sources.FramesAllDescriptionsFilename)
	jlistPathBak := filepath.Join(dataPath, sources.FramesAllDescriptionsFilename+".bak")

	err = filelib.CopyFile(jlistPath, jlistPathBak)
	require.NoError(t, err)

	var descrs []sources.Description
	err = serialization.ReadAllPartsJSON(jlistPath, &descrs)
	require.NoError(t, err)

	for i := range descrs {
		descrs[i].DPM = 2
	}

	err = serialization.SaveAllPartsJSON(descrs, jlistPath)
	require.NoError(t, err)
}
