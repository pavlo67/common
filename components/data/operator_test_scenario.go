package data

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/tagger"
)

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	ToSave   Item
	ToUpdate Item

	//DetailsToSave      Test
	//DetailsToUpdate    Test
}

func TestCases(dataOp Operator, cleanerOp crud.Cleaner) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: dataOp,
			Cleaner:  cleanerOp,
			ToSave: Item{
				ID:      "",
				URL:     "111111",
				Title:   "345456",
				Summary: "6578gj",
				Embedded: []Item{{
					Title:   "56567",
					Summary: "3333333",
					Tags:    []tagger.Tag{{Label: "1"}, {Label: "332343"}},
				}},
				Data: crud.Data{
					TypeKey: "test",
					Content: []byte(`{"AAA": "aaa", "BBB": 222}`),
				},
				Tags: []tagger.Tag{{Label: "1"}, {Label: "333"}},
				History: []crud.Action{{
					Key:    crud.CreatedAction,
					DoneAt: time.Time{},
				}},
			},

			ToUpdate: Item{
				URL:     "22222222",
				Title:   "345456rt",
				Summary: "6578eegj",
				Data: crud.Data{
					TypeKey: "test",
					Content: []byte(`{"AAA": "awraa", "BBB": 22552}`),
				},
				Tags: []tagger.Tag{{Label: "1"}, {Label: "333"}},
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

func Compare(t *testing.T, dataOp Operator, readed *Item, expectedItem Item, l logger.Operator) {
	require.NotNil(t, readed)

	l.Infof("to be saved: %#v", expectedItem)
	l.Infof("readed: %#v", readed)

	for i, action := range expectedItem.History {
		expectedItem.History[i].DoneAt = action.DoneAt.UTC()
	}

	expectedDetails := expectedItem.Data.Content
	expectedItem.Data.Content = nil

	readedDetails := readed.Data.Content
	readed.Data.Content = nil

	// TODO!!! check it carefully
	readed.History = nil
	expectedItem.History = nil

	require.Equal(t, &expectedItem, readed)
	require.Equal(t, expectedDetails, readedDetails)

}

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		var id [numRepeats1 + numRepeats2]common.ID
		var toSave [numRepeats1 + numRepeats2]Item
		// var data Tag

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean(nil, nil)
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
		//	key := field.ID
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
		//	_, err = tc.Save([]Tag{tc.ToSave}, nil)
		//	require.ErrStr(t, err, "where is an error on .Save()?")
		//	continue
		//}

		for i := 0; i < numRepeats1; i++ {
			toSave[i] = tc.ToSave
			//toSave[i].Details = &tc.DetailsToSave
			idsI, err := tc.Save([]Item{toSave[i]}, nil)
			require.NoError(t, err)
			require.True(t, len(idsI) == 1)
			id[i] = idsI[0]
		}

		var toSavePack []Item
		//tc.ToSave.Details = &tc.DetailsToSave
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
		//	 require.ErrStr(t, err)
		//	 continue
		// }

		readedSaved, err := tc.Read(id[toReadI], nil)
		require.NoError(t, err)

		tc.ToSave.ID = id[toReadI]

		Compare(t, tc, readedSaved, tc.ToSave, l)

		// test .Update & .Read -----------------------------------------------------------------------------------

		// if !tc.ExpectedUpdateOk {
		//	 err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//	 require.ErrStr(t, err, "where is an error on .Update()?")
		//	 continue
		// }

		tc.ToUpdate.ID = id[toUpdateI]
		// tc.ToUpdate.Details = &tc.DetailsToUpdate

		_, err = tc.Save([]Item{tc.ToUpdate}, nil)
		require.NoError(t, err)

		readedUpdated, err := tc.Read(id[toUpdateI], nil)
		require.NoError(t, err)

		// tc.ToUpdate.History.CreatedAt = tc.ToSave.History.CreatedAt // unchanged!!!

		Compare(t, tc, readedUpdated, tc.ToUpdate, l)

		// TODO!!!
		//	if !tc.ExcludeUpdateTest {
		//		var uniquesUpdatable []string
		//		for _, field := range description.FieldsArr {
		//			if field.Unique && (field.Updatable && !field.AutoUnique) { // || field.Additable
		//				uniquesUpdatable = append(uniquesUpdatable, field.ID)
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
		//			require.ErrStr(t, err)
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
		//		toUpdate := Tag{}
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
		//			require.ErrStr(t, err)
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
		//		require.ErrStr(t, err)
		//	}
		//

		//toUpdateResult := tc.ToUpdate
		//for _, f := range description.FieldsArr {
		//	if !f.Creatable {
		//		toUpdateResult[f.ID] = data[f.ID]
		//	}
		//}

		// test ListTags -------------------------------------------------------------------------------------

		//if !tc.ExcludeListTest {
		//	var ids []common.ID
		//	for _, idi := range id {
		//		ids = append(ids, idi)
		//	}
		//
		//	if !tc.ExpectedReadOk {
		//		// TODO: selector.InStr(keyFields[0], ids...)
		//		briefsAll, err := tc.ListTags(nil, nil)
		//
		//		require.Equal(t, 0, len(briefsAll), "why len(dataAll) is not zero after .ListTags()?")
		//		require.ErrStr(t, err)
		//		continue
		//	}
		//
		//	// TODO: selector.InStr(keyFields[0], ids...)

		briefsAll, err := tc.List(nil, &crud.GetOptions{OrderBy: []string{"id"}})
		require.NoError(t, err)
		require.True(t, len(briefsAll) == numRepeats1+numRepeats2)

		Compare(t, tc, &briefsAll[toReadI], tc.ToSave, l)
		Compare(t, tc, &briefsAll[toUpdateI], tc.ToUpdate, l)

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
		//			require.ErrStr(t, err, "where is an error on .DeleteList()?")
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
		//			require.ErrStr(t, err, "where is an error on .DeleteList()?")
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
		//		// require.ErrStr(t, err, "where is an error on .Read() after DeleteList()?")
		//
		//		require.Nil(t, nativeToRead)
		//	}

	}
}
