package founts_leveldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/processor/founts"
)

func TestCRUD(t *testing.T) {
	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment set!!!")
	}

	fountsOp, err := New("test")

	operatorCRUD := founts.OperatorCRUD{fountsOp}
	testCases, err := operatorCRUD.TestCases(func() error { return fountsOp.clean() })
	require.NoError(t, err)

	crud.OperatorTest(t, testCases)

	err = fountsOp.Close()
	require.NoError(t, err)
}
