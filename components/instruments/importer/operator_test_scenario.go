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
		items, err := tc.Operator.Get(tc.Source, nil)
		require.NoError(t, err)
		require.True(t, len(items) > 0)

		for _, item := range items {
			log.Printf("%#v", item)
		}
	}
}
