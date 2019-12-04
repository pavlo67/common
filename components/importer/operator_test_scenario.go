package importer

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type ImporterTestCase struct {
	Operator Operator
	Source   string
}

func TestImporterWithCases(t *testing.T, testCases []ImporterTestCase) {
	for _, tc := range testCases {
		series, err := tc.Operator.Get(tc.Source)
		require.NoError(t, err)
		require.NotNil(t, series)
		require.True(t, len(series.Items) > 0)

		for _, item := range series.Items {
			log.Printf("%#v", item)
		}
	}
}
