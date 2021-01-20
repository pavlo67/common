package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func ScenarioContentImportExportItself(t *testing.T, contentOriginal Content) {

	// contentOriginal --> List
	listOriginal, descrOriginal, err := contentOriginal.Export()
	require.NoError(t, err)

	// contentOriginal (duplicated) <-- List
	contentDuplicated := contentOriginal.NewEmpty()
	err = contentDuplicated.Import(listOriginal, descrOriginal)
	require.NoError(t, err)
	listDuplicated, descrDuplicated, err := contentDuplicated.Export()
	require.NoError(t, err)

	// contentOriginal must be equal to contentOriginal (duplicated)
	err = descrOriginal.IsEqualTo(descrDuplicated)
	require.NoError(t, err)
	err = descrOriginal.ValuesAreEqual(listOriginal, listDuplicated)
	require.NoError(t, err)
}

func ScenarioContentImportExport(t *testing.T, content0, contentOriginalImported Content) {

	// content0 --> contentOriginalImported
	list0, descr0, err := content0.Export()
	require.NoError(t, err)
	contentOriginalImported = contentOriginalImported.NewEmpty()
	err = contentOriginalImported.Import(list0, descr0)
	require.NoError(t, err)
	listOriginal, descrOriginal, err := contentOriginalImported.Export()
	require.NoError(t, err)

	// contentOriginalImported --> content0
	content0 = content0.NewEmpty()
	err = content0.Import(listOriginal, descrOriginal)
	require.NoError(t, err)
	list1, descr1, err := content0.Export()
	require.NoError(t, err)

	// content0 --> contentOriginalImported (duplicated)
	contentOriginalImported = contentOriginalImported.NewEmpty()
	err = contentOriginalImported.Import(list1, descr1)
	require.NoError(t, err)
	listDuplicated, descrDuplicated, err := contentOriginalImported.Export()
	require.NoError(t, err)

	// contentOriginalImported must be equal to contentOriginalImported (duplicated)
	err = descrOriginal.IsEqualTo(descrDuplicated)
	require.NoError(t, err)
	err = descrOriginal.ValuesAreEqual(listOriginal, listDuplicated)
	require.NoError(t, err)

}
