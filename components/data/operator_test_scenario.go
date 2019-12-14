package data

import (
	"os"
	"testing"

	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/tagger"
	"github.com/stretchr/testify/require"
)

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	ToSave   Item
	ToUpdate Item

	DetailsToSave      Test
	DetailsToReadSaved Test

	DetailsToUpdate      Test
	DetailsToReadUpdated Test
}

type Test struct {
	AAA string
	BBB int
}

const TypeKeyTest TypeKey = "test"

var TypeTest = Type{
	Key:      TypeKeyTest,
	Exemplar: Test{},
}

var TypeString = Type{
	Key:      TypeKeyString,
	Exemplar: "",
}

func TestCases(dataOp Operator, cleanerOp crud.Cleaner) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: dataOp,
			Cleaner:  cleanerOp,
			ToSave: Item{
				ID:       "",
				TypeKey:  "test",
				ExportID: "rtuy",
				URL:      "111111",
				Title:    "345456",
				Summary:  "6578gj",
				Embedded: []Item{{
					ExportID: "wq3r",
					Title:    "56567",
					Summary:  "3333333",
					Tags:     []tagger.Tag{"1", "332343"},
				}},
				Tags: []tagger.Tag{"1", "333"},
				Status: crud.Status{
					CreatedAt: time.Now(),
				},
				Origin: flow.Origin{},
			},
			DetailsToSave: Test{
				AAA: "aaa",
				BBB: 222,
			},

			ToUpdate: Item{
				URL:     "22222222",
				Title:   "345456rt",
				Summary: "6578eegj",
				Tags:    []tagger.Tag{"1", "333"},
				Status: crud.Status{
					CreatedAt: time.Now().Add(time.Minute),
				},
			},
			DetailsToUpdate: Test{
				AAA: "awraa",
				BBB: 22552,
			},
		},
	}
}

// TODO: тест чистки бази
// TODO: test created_at, updated_at
// TODO: test GetOptions

const numRepeats1 = 2
const numRepeats2 = 3
const toReadI = 0   // must be < numRepeats1 + numRepeats2
const toUpdateI = 1 // must be < numRepeats1 + numRepeats2
const toDeleteI = 2 // must be < numRepeats1 + numRepeats2

func Compare(t *testing.T, dataOp Operator, readed *Item, expectedItem Item, expectedDetails, detailsToRead Test, l logger.Operator) {
	require.NotNil(t, readed)

	err := dataOp.SetDetails(readed)
	require.NoError(t, err)

	l.Infof("to be saved: %#v", expectedItem)
	l.Infof("readed: %#v", readed)
	l.Infof("readed details: %#v", detailsToRead)

	expectedItem.Status.CreatedAt = expectedItem.Status.CreatedAt.UTC()
	expectedItem.Details = nil
	expectedItem.DetailsRaw = nil

	readed.Details = nil
	readed.DetailsRaw = nil

	// kostyl!!!
	require.Equal(t, expectedItem.Status.CreatedAt.Format(time.RFC3339), readed.Status.CreatedAt.Format(time.RFC3339))
	readed.Status.CreatedAt = expectedItem.Status.CreatedAt

	require.Equal(t, &expectedItem, readed)
	require.Equal(t, expectedDetails, detailsToRead)

}

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		var id [numRepeats1 + numRepeats2]common.ID
		var toSave [numRepeats1 + numRepeats2]Item
		// var data Item

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean(nil)
		require.NoError(t, err, "what is the error on .Cleaner()?")

		// test Describe ------------------------------------------------------------------------------------

		//description := tc.Description()
		//
		//keyFields := description.PrimaryKeys()
		//
		//if len(keyFields) > 1 {
		//	require.FailNow(t, "too many key fields", keyFields)
		//} else if len(keyFields) < 1 {
		//	keyFields = append(keyFields, "id")
		//}
		//
		////for _, fieldKey := range tc.DescribedFields {
		////	require.NotEmpty(t, description.FieldsArr[fieldKey], "on .Describe(): "+fieldKey+"???")
		////}

		// test Create --------------------------------------------------------------------------------------

		//var uniques, autoUniques []string
		//
		//for _, field := range description.FieldsArr {
		//	key := field.Key
		//	if field.Unique {
		//		if field.AutoUnique {
		//			autoUniques = append(autoUniques, key)
		//		} else {
		//			uniques = append(uniques, key)
		//		}
		//	}
		//}
		//
		//nativeToCreate, err := tc.ItemToNative(tc.ToSave)
		//require.NoError(t, err)

		//if !tc.ExpectedSaveOk {
		//	_, err = tc.Save([]Item{tc.ToSave}, nil)
		//	require.Error(t, err, "where is an error on .Save()?")
		//	continue
		//}

		for i := 0; i < numRepeats1; i++ {
			toSave[i] = tc.ToSave
			toSave[i].Details = &tc.DetailsToSave
			idsI, err := tc.Save([]Item{toSave[i]}, nil)
			require.NoError(t, err)
			require.True(t, len(idsI) == 1)
			id[i] = idsI[0]
		}

		var toSavePack []Item
		tc.ToSave.Details = &tc.DetailsToSave
		for j := 0; j < numRepeats2; j++ {
			toSavePack = append(toSavePack, tc.ToSave)
		}
		idsI, err := tc.Save(toSavePack, nil)
		require.NoError(t, err)
		require.True(t, len(idsI) == numRepeats2)
		for j := 0; j < numRepeats2; j++ {
			id[numRepeats1+i] = idsI[i]
		}

		// test .Read ----------------------------------------------------------------------------------------

		// if !tc.ExpectedReadOk {
		// 	 _, err = tc.Read(id[toReadI], nil)
		//	 require.Error(t, err)
		//	 continue
		// }

		readedSaved, err := tc.Read(id[toReadI], nil)
		require.NoError(t, err)

		tc.ToSave.ID = id[toReadI]

		Compare(t, tc, readedSaved, tc.ToSave, tc.DetailsToSave, tc.DetailsToReadSaved, l)

		// test .Update & .Read -----------------------------------------------------------------------------------

		// if !tc.ExpectedUpdateOk {
		//	 err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//	 require.Error(t, err, "where is an error on .Update()?")
		//	 continue
		// }

		tc.ToUpdate.ID = id[toUpdateI]
		tc.ToUpdate.Details = &tc.DetailsToUpdate

		_, err = tc.Save([]Item{tc.ToUpdate}, nil)
		require.NoError(t, err)

		readedUpdated, err := tc.Read(id[toUpdateI], nil)
		require.NoError(t, err)

		tc.ToUpdate.ExportID = tc.ToSave.ExportID                 // unchanged!!!
		tc.ToUpdate.Origin = tc.ToSave.Origin                     // unchanged!!!
		tc.ToUpdate.Status.CreatedAt = tc.ToSave.Status.CreatedAt // unchanged!!!

		Compare(t, tc, readedUpdated, tc.ToUpdate, tc.DetailsToUpdate, tc.DetailsToReadUpdated, l)

		// TODO!!!
		//	if !tc.ExcludeUpdateTest {
		//		var uniquesUpdatable []string
		//		for _, field := range description.FieldsArr {
		//			if field.Unique && (field.Updatable && !field.AutoUnique) { // || field.Additable
		//				uniquesUpdatable = append(uniquesUpdatable, field.Key)
		//			}
		//		}
		//
		//		//tc.ToUpdate[keyFields[0]] = id[toUpdateI]
		//
		//		nativeToUpdate, err := tc.ItemToNative(tc.ToUpdate)
		//		require.NoError(t, err)
		//
		//
		//		if tc.ISToUpdateBad != nil {
		//			err = tc.Update(*tc.ISToUpdateBad, id[toUpdateI], nativeToUpdate)
		//			require.Error(t, err)
		//		}
		//
		//		// update 1: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is an error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		// update 2: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is an error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToUpdate, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		toUpdate := Item{}
		//		for k, v := range toUpdateResult {
		//			toUpdate[k] = v
		//		}
		//
		//		// can't duplicate uniques fields
		//		for _, key := range uniquesUpdatable {
		//			toUpdate[key] = toSave[0][key]
		//			nativeToUpdate, err := tc.ItemToNative(toUpdate)
		//			require.NoError(t, err)
		//			err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//			require.Error(t, err)
		//			toUpdate[key] = toUpdateResult[key]
		//		}
		//
		//		// update 3: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is the error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		// can't update absent record
		//		toUpdate[keyFields[0]] += "123"
		//		nativeToUpdate, err = tc.ItemToNative(toUpdate)
		//		require.NoError(t, err)
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.Error(t, err)
		//	}
		//

		//toUpdateResult := tc.ToUpdate
		//for _, f := range description.FieldsArr {
		//	if !f.Creatable {
		//		toUpdateResult[f.Key] = data[f.Key]
		//	}
		//}

		// test List -------------------------------------------------------------------------------------

		//if !tc.ExcludeListTest {
		//	var ids []common.ID
		//	for _, idi := range id {
		//		ids = append(ids, idi)
		//	}
		//
		//	if !tc.ExpectedReadOk {
		//		// TODO: selector.InStr(keyFields[0], ids...)
		//		briefsAll, err := tc.List(nil, nil)
		//
		//		require.Equal(t, 0, len(briefsAll), "why len(dataAll) is not zero after .List()?")
		//		require.Error(t, err)
		//		continue
		//	}
		//
		//	// TODO: selector.InStr(keyFields[0], ids...)

		briefsAll, err := tc.List(nil, nil)
		require.NoError(t, err)
		require.True(t, len(briefsAll) == numRepeats1+numRepeats2)

		Compare(t, tc, &briefsAll[toReadI], tc.ToSave, tc.DetailsToSave, tc.DetailsToReadSaved, l)
		Compare(t, tc, &briefsAll[toUpdateI], tc.ToUpdate, tc.DetailsToUpdate, tc.DetailsToReadUpdated, l)

		// test .Delete --------------------------------------------------------------------------------------

		err = tc.Remove(id[toDeleteI], nil)
		require.NoError(t, err)

		readDeleted, err := tc.Read(id[toDeleteI], nil)
		require.Error(t, err)
		require.Nil(t, readDeleted)

		briefsAll, err = tc.List(nil, nil)
		require.NoError(t, err)
		require.True(t, len(briefsAll) == numRepeats1+numRepeats2-1)

		//	if !tc.ExcludeRemoveTest {
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//
		//		if tc.ExpectedRemoveErr != nil {
		//			err = tc.Delete(tc.ISToDelete, id[toDeleteI])
		//			require.Error(t, err, "where is an error on .DeleteList()?")
		//			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//			require.NoError(t, err, "what is the error on .Read() after Update()?")
		//			data, err = tc.NativeToItem(nativeToRead)
		//			require.NoError(t, err)
		//			require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//			continue
		//		}
		//
		//		if tc.ISToDeleteBad != nil {
		//			err = tc.Delete(*tc.ISToDeleteBad, id[toDeleteI])
		//			require.Error(t, err, "where is an error on .DeleteList()?")
		//			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//			require.NoError(t, err, "what is the error on .Read() after Update()?")
		//			data, err = tc.NativeToItem(nativeToRead)
		//			require.NoError(t, err)
		//			require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//		}
		//
		//		err = tc.Delete(tc.ISToDelete, id[toDeleteI])
		//		require.NoError(t, err, "what is the error on .DeleteList()?")
		//
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//
		//		// it depends on implementation
		//		// require.Error(t, err, "where is an error on .Read() after DeleteList()?")
		//
		//		require.Nil(t, nativeToRead)
		//	}

	}
}

//func testData(t *testing.T, keyFields, expectedID []string, expectedData, data Item, onCreate bool, description Description, on string) {
//	if expectedData == nil {
//		require.Nil(t, data)
//		return
//	}
//	require.NotNil(t, data)
//
//	require.Equal(t, len(keyFields), len(expectedID))
//	for i, f := range keyFields {
//		require.Equal(t, expectedID[i], data[f], on+": incorrect key value in field '%s'???", f)
//	}
//
//	for _, field := range description.FieldsArr {
//		key := field.Key
//
//		// TODO: check key field
//
//		if (onCreate && field.Creatable) || (!onCreate && field.Updatable) {
//			if expectedData[key] == "" && field.NotEmpty {
//				require.NotEmpty(t, data[key], on+": "+key+"???")
//			} else {
//				require.Equal(t, expectedData[key], data[key], on+": "+key+"???")
//			}
//
//			//} else if !onCreate && field.Additable {
//			//	if expectedData[key] == "" {
//			//		require.Equal(t, expectedData[key], data[key], on+": "+key+"???")
//			//	} else {
//			//		require.True(t, len(data[key]) > len(expectedData[key]), on+": "+key+"???")
//			//	}
//
//		} else if field.NotEmpty {
//			require.NotEmpty(t, data[key], on+": "+key+"???")
//		}
//	}
//}
