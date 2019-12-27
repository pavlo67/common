package filesloader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
)

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	PathToStore string
	URLToLoad   string
}

func TestCases(flOp Operator, cleanerOp crud.Cleaner, pathToStore string) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator:    flOp,
			Cleaner:     cleanerOp,
			PathToStore: pathToStore,
			URLToLoad:   "http://abc.ru",
		},
	}
}

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		//// ClearDatabase ------------------------------------------------------------------------------------
		//
		//err := tc.Cleaner.Clean(nil, nil)
		//require.NoError(t, err, "what is the error on .Cleaner()?")

		// test .Load -----------------------------------------------------------------------------------------

		item, err := tc.Load(tc.URLToLoad, "", nil)
		require.NoError(t, err)
		require.NotNil(t, item)

		l.Infof("%#v", item)

		require.Equal(t, tc.PathToStore, item.Path[:len(tc.PathToStore)])

		isDir, err := item.IsDir()
		require.NoError(t, err)
		require.True(t, isDir)

		files, err := item.FilesList()
		require.NoError(t, err)
		require.True(t, len(files) > 0)
		require.Equal(t, "0.html", files[0].Name())
		require.False(t, files[0].IsDir())
		require.True(t, files[0].Size() > 0)

	}
}
