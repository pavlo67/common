package exporter

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/tagger"

	"github.com/pavlo67/workshop/components/data"
)

type TestCase struct {
	Operator
	CleanerOp crud.Cleaner

	Items  []data.Item   // TODO: remove the kostyl!!! use crud.Operator
	DataOp data.Operator // TODO: remove the kostyl!!! use crud.Operator
}

func TestCases(exporterOp Operator, dataOp data.Operator, cleanerOp crud.Cleaner) []TestCase {
	return []TestCase{
		{
			Operator:  exporterOp,
			DataOp:    dataOp,
			CleanerOp: cleanerOp,
			Items: []data.Item{
				{
					Key:     "etyrty",
					URL:     "ewrr",
					Title:   "ryuty",
					Summary: "wqwr",
					Tags:    []tagger.Tag{{Label: "werwe"}},
					Data: crud.Data{
						TypeKey: crud.StringTypeKey,
						Content: "wqer 3wrt5we ewrt",
					},
				},
				{
					Key:     "etyrty1",
					URL:     "ewrr1",
					Title:   "ryuty1",
					Summary: "wqwr1",
					Tags:    []tagger.Tag{{Label: "werwe1"}},
					Data: crud.Data{
						TypeKey: crud.StringTypeKey,
						Content: "wqer 3wrt5we ewrt1",
					},
				},
			},
		},
	}
}

func OperatorTestScenario(t *testing.T, testCases []TestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		// clear database -----------------------------------------------------------------------------------

		err := tc.CleanerOp.Clean(nil, nil)
		require.NoError(t, err, "what is the error on .Cleaner()?")

		// fill database ------------------------------------------------------------------------------------

		for i := 0; i < len(tc.Items); i++ {
			idI, err := tc.DataOp.Save(tc.Items[i], nil)
			require.NoError(t, err)
			require.NotEmpty(t, idI)
			tc.Items[i].ID = idI
		}

		// test export ----------------------------------------------------------------------------------------

		exportedData, items := exportDataItems(t, tc, "", tc.Items)

		// test import ----------------------------------------------------------------------------------------

		importAndExportCheck(t, tc, exportedData, items)

		// test import after cleaning -------------------------------------------------------------------------

		err = tc.CleanerOp.Clean(nil, nil)
		require.NoError(t, err)

		importAndExportCheck(t, tc, exportedData, items)

		// test export by parts -------------------------------------------------------------------------------

		exportDataItemsByParts(t, tc, tc.Items)

	}
}

func exportDataItems(t *testing.T, tc TestCase, after string, expectedItems []data.Item) (crud.Data, []data.Item) {
	exportedData, err := tc.Operator.Export(nil, after, nil)
	require.NoError(t, err)

	require.NotNil(t, exportedData)
	require.Equal(t, data.ItemsTypeKey, exportedData.TypeKey)

	var items []data.Item
	err = json.Unmarshal([]byte(exportedData.Content), &items)
	require.NoError(t, err)

	require.Equal(t, len(expectedItems), len(items))
	for i, item := range items {
		expectedItem := expectedItems[i]

		expectedItem.History = nil
		expectedItem.ID = ""

		item.History = nil
		item.ID = ""

		require.Equal(t, expectedItem, item)
		// require.True(t, reflect.DeepEqual(expectedItems[i], item))
	}

	return *exportedData, items
}

func importAndExportCheck(t *testing.T, tc TestCase, toImport crud.Data, expected []data.Item) {
	till, err := tc.Operator.Import(toImport, nil)
	require.NoError(t, err)
	require.Equal(t, string(expected[len(expected)-1].ID), till)

	_, _ = exportDataItems(t, tc, "", expected)

}

func exportDataItemsByParts(t *testing.T, tc TestCase, expectedItems []data.Item) {

	// the first part --------------------------------

	exportedData, err := tc.Operator.Export(nil, "", &crud.GetOptions{Limit: 1})
	require.NoError(t, err)

	require.NotNil(t, exportedData)
	require.Equal(t, data.ItemsTypeKey, exportedData.TypeKey)

	var items []data.Item
	err = json.Unmarshal([]byte(exportedData.Content), &items)
	require.NoError(t, err)

	require.Equal(t, 1, len(items))

	after := string(items[0].ID)

	// the rest --------------------------------------

	exportedData, err = tc.Operator.Export(nil, after, nil)
	require.NoError(t, err)

	require.NotNil(t, exportedData)
	require.Equal(t, data.ItemsTypeKey, exportedData.TypeKey)

	var itemsRest []data.Item
	err = json.Unmarshal([]byte(exportedData.Content), &itemsRest)
	require.NoError(t, err)

	require.Equal(t, len(expectedItems)-1, len(itemsRest))

	// check all -------------------------------------

	items = append(items, itemsRest...)

	require.Equal(t, len(expectedItems), len(items))
	for i, item := range items {
		expectedItem := expectedItems[i]

		expectedItem.History = nil
		expectedItem.ID = ""

		item.History = nil
		item.ID = ""

		require.Equal(t, expectedItem, item)
		// require.True(t, reflect.DeepEqual(expectedItems[i], item))
	}
}
