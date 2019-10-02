package manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadManifest(t *testing.T) {
	manifest, err := ReadManifest("test_data")
	require.NoError(t, err)
	require.NotNil(t, manifest)

	fmt.Printf("%#v", manifest)

}
