package dataimporter

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
		require.True(t, len(series.Data) > 0)

		for _, item := range series.Data {
			log.Printf("%#v", item)
		}
	}
}
