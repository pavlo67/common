package crud

import (
	"os"
	"testing"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/stretchr/testify/require"
)

type OperatorTestCase struct {
	Operator
	Cleaner

	DetailsToRead interface{}

	ToSave          Item
	ExpectedSaveErr error
	ExpectedReadErr error

	ExpectedListErr error
	ExcludeListTest bool

	ToUpdate          Item
	ExpectedUpdateErr error
	ExcludeUpdateTest bool

	ExpectedRemoveErr error
	ExcludeRemoveTest bool
}

// TODO: тест чистки бази
// TODO: test created_at, updated_at
// TODO: test GetOptions

const numRepeats = 3
const toReadI = 0   // must be < numRepeats
const toUpdateI = 1 // must be < numRepeats
const toDeleteI = 2 // must be < numRepeats

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		var id [numRepeats]common.ID
		var toSave [numRepeats]Item
		// var data Item

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean()
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

		if tc.ExpectedSaveErr != nil {
			_, err = tc.Save(tc.ToSave, nil)
			require.Error(t, err, "where is an error on .Save()?")
			continue
		}

		for i := 0; i < numRepeats; i++ {
			toSave[i] = tc.ToSave

			idI, err := tc.Save(toSave[i], nil)
			require.NoError(t, err, "what is the error on .Create()?")
			require.NotEmpty(t, idI)

			id[i] = *idI
		}

		// test Read ----------------------------------------------------------------------------------------

		if tc.ExpectedReadErr != nil {
			_, err = tc.Read(id[toReadI], nil)
			require.Error(t, err)
			continue
		}

		l.Infof("saved: %#v", tc.ToSave)

		doc, err := tc.Read(id[toReadI], nil)
		require.NoError(t, err)

		l.Infof("readed: %#v", doc)

		err = tc.Details(doc, tc.DetailsToRead)
		require.NoError(t, err)

		l.Infof("readed details: %#v", tc.DetailsToRead)

		// TODO!!!
		// testData(t, nil, []string{string(id[toReadI])}, toSave[toReadI], data, true, "on .Read()")

		//toUpdateResult := tc.ToUpdate
		//for _, f := range description.FieldsArr {
		//	if !f.Creatable {
		//		toUpdateResult[f.Key] = data[f.Key]
		//	}
		//}

		// test List -------------------------------------------------------------------------------------

		if !tc.ExcludeListTest {
			var ids []common.ID
			for _, idi := range id {
				ids = append(ids, idi)
			}

			if tc.ExpectedReadErr != nil {
				// TODO: selectors.InStr(keyFields[0], ids...)
				briefsAll, err := tc.List(nil, nil)

				require.Equal(t, 0, len(briefsAll), "why len(dataAll) is not zero after .List()?")
				require.Error(t, err)
				continue
			}

			// TODO: selectors.InStr(keyFields[0], ids...)
			briefsAll, err := tc.List(nil, nil)
			require.NoError(t, err, "what is the error on .ReadList()?")
			require.True(t, len(briefsAll) >= numRepeats, "must be len(dataAll) (%d) >= numRepeats (%d)", len(briefsAll), numRepeats)

			// TODO!!!
			//for i, native := range nativeAll {
			//	testData(t, keyFields, []string{id[i]}, toSave[i], data, true, description, "on .ReadList()")
			//}
		}

		//	// test Update --------------------------------------------------------------------------------------
		//
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
		//		if tc.ExpectedUpdateErr != nil {
		//			err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//			require.Error(t, err, "where is an error on .Update()?")
		//			continue
		//		}
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
		//	// test DeleteList --------------------------------------------------------------------------------------
		//
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
