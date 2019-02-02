package crud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/partes/selectors"
	"github.com/pavlo67/punctum/auth"
)

type OperatorTestCase struct {
	Operator
	Cleaner

	KeyField string

	ISToCreate        auth.ID
	ISToCreateBad     *auth.ID
	ToCreate          StringMap
	ExpectedCreateErr error

	ISToRead        auth.ID
	ISToReadBad     *auth.ID
	ExpectedReadErr error

	ISToUpdate        auth.ID
	ISToUpdateBad     *auth.ID
	ToUpdate          StringMap
	ExpectedUpdateErr error

	ISToDelete        auth.ID
	ISToDeleteBad     *auth.ID
	ExpectedDeleteErr error

	ExcludeReadListTest bool
	ExcludeUpdateTest   bool
	ExcludeDeleteTest   bool
}

// TODO: тест чистки бази
// TODO: test r_view, r_owner, managers change
// TODO: test created_at, updated_at
// TODO: test ReadOptions

const numRepeats = 3
const toReadI = 0   // must be < numRepeats
const toUpdateI = 1 // must be < numRepeats
const toDeleteI = 2 // must be < numRepeats

func OperatorTest(t *testing.T, testCases []OperatorTestCase) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println(i)

		var id [numRepeats]string
		var toCreate [numRepeats]StringMap
		var data StringMap

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner()
		require.NoError(t, err, "what is the error on .Cleaner()?")

		// test Describe ------------------------------------------------------------------------------------

		if tc.KeyField == "" {
			tc.KeyField = "id"
		}

		description, err := tc.Describe()
		require.NoError(t, err, "what is the error on .Describe()?")

		//for _, fieldKey := range tc.DescribedFields {
		//	require.NotEmpty(t, description.Fields[fieldKey], "on .Describe(): "+fieldKey+"???")
		//}

		// test Create --------------------------------------------------------------------------------------

		var uniques, autoUniques []string

		for _, field := range description.Fields {
			key := field.Key
			if field.Unique {
				if field.AutoUnique {
					autoUniques = append(autoUniques, key)
				} else {
					uniques = append(uniques, key)
				}
			}
		}

		nativeToCreate, err := tc.StringMapToNative(tc.ToCreate)
		require.NoError(t, err)

		if tc.ExpectedCreateErr != nil {
			_, err = tc.Create(tc.ISToCreate, nativeToCreate)
			require.Error(t, err, "where is an error on .Create()?")
			continue
		}

		if tc.ISToCreateBad != nil {
			_, err = tc.Create(*tc.ISToCreateBad, nativeToCreate)
			require.Error(t, err, "where is an error on .Create()?")
		}

		for i := 0; i < numRepeats; i++ {
			toCreate[i] = StringMap{}
			for k, v := range tc.ToCreate {
				toCreate[i][k] = v
			}

			if i > 0 {
				for _, key := range autoUniques {
					toCreate[i][key] = toCreate[0][key] + strings.Repeat("n", i)
				}

				for _, key := range uniques {
					toCreate[i][key] = toCreate[0][key] + strings.Repeat("n", i)

					buf := toCreate[i][key]
					toCreate[i][key] = toCreate[0][key]
					// testing error!!!

					nativeToCreateI, err := tc.StringMapToNative(toCreate[i])
					require.NoError(t, err)

					_, err = tc.Create(tc.ISToCreate, nativeToCreateI)
					require.Errorf(t, err, "%#v", nativeToCreateI)

					toCreate[i][key] = buf
				}

			}

			nativeI, err := tc.StringMapToNative(toCreate[i])
			require.NoError(t, err)

			id[i], err = tc.Create(tc.ISToCreate, nativeI)
			require.NoError(t, err, "what is the error on .Create()?")
			require.NotEmpty(t, id[i])
		}

		// test Read ----------------------------------------------------------------------------------------

		if tc.ExpectedReadErr != nil {
			_, err = tc.Read(tc.ISToRead, id[toReadI])
			require.Error(t, err, "where is an error on .Read()?")
			continue
		}

		if tc.ISToReadBad != nil {
			_, err = tc.Read(*tc.ISToReadBad, id[toReadI])
			require.Error(t, err, "where is an error on .Read()?")
		}

		nativeToRead, err := tc.Read(tc.ISToRead, id[toReadI])
		require.NoError(t, err, "what is the error on .Read()?")

		data, err = tc.NativeToStringMap(nativeToRead)
		require.NoError(t, err)
		testData(t, tc.KeyField, id[toReadI], toCreate[toReadI], data, true, description, "on .Read()")

		toUpdateResult := tc.ToUpdate
		for _, f := range description.Fields {
			if !f.Creatable {
				toUpdateResult[f.Key] = data[f.Key]
			}
		}

		// test ReadList -------------------------------------------------------------------------------------

		if !tc.ExcludeReadListTest {
			var ids []interface{}
			for _, idi := range id {
				ids = append(ids, idi)
			}

			options := &ReadOptions{
				Selector: selectors.FieldEqual(tc.KeyField, ids...),
				SortBy:   []string{tc.KeyField},
			}
			if tc.ExpectedReadErr != nil {
				nativeAll, _, err := tc.ReadList(tc.ISToRead, options)
				require.Equal(t, 0, len(nativeAll), "why len(dataAll) is not zero after .ReadList()?")
				require.Error(t, err)
				continue
			}

			nativeAll, _, err := tc.ReadList(tc.ISToRead, options)
			require.NoError(t, err, "what is the error on .ReadList()?")
			require.True(t, len(nativeAll) >= numRepeats, "must be len(dataAll) (%d) >= numRepeats (%d)", len(nativeAll), numRepeats)

			// require.True(t, numAll == uint64(len(nativeAll)), "must be numAll (%d) == len(dataAll) (%d)", numAll, len(nativeAll))

			// require.Equal(t, numAll, uint64(len(dataAll)))

			for i, native := range nativeAll {
				data, err := tc.NativeToStringMap(native)
				require.NoError(t, err)
				testData(t, tc.KeyField, id[i], toCreate[i], data, true, description, "on .ReadList()")
			}
		}

		// test Update --------------------------------------------------------------------------------------

		if !tc.ExcludeUpdateTest {
			var uniquesUpdatable []string
			for _, field := range description.Fields {
				if field.Unique && (field.Updatable || field.Additable) && !field.AutoUnique {
					uniquesUpdatable = append(uniquesUpdatable, field.Key)
				}
			}

			tc.ToUpdate[tc.KeyField] = id[toUpdateI]

			nativeToUpdate, err := tc.StringMapToNative(tc.ToUpdate)
			require.NoError(t, err)

			if tc.ExpectedUpdateErr != nil {
				err = tc.Update(tc.ISToUpdate, nativeToUpdate)
				require.Error(t, err, "where is an error on .Update()?")
				continue
			}

			if tc.ISToUpdateBad != nil {
				err = tc.Update(*tc.ISToUpdateBad, nativeToUpdate)
				require.Error(t, err)
			}

			// update 1: ok
			err = tc.Update(tc.ISToUpdate, nativeToUpdate)
			require.NoError(t, err, "what is an error on .Update()?")
			nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
			require.NoError(t, err, "what is the error on .Read() after Update()?")
			data, err = tc.NativeToStringMap(nativeToRead)
			require.NoError(t, err)
			testData(t, tc.KeyField, id[toUpdateI], toUpdateResult, data, false, description, "on .Read() after Update()")

			// update 2: ok
			err = tc.Update(tc.ISToUpdate, nativeToUpdate)
			require.NoError(t, err, "what is an error on .Update()?")
			nativeToRead, err = tc.Read(tc.ISToUpdate, id[toUpdateI])
			require.NoError(t, err, "what is the error on .Read() after Update()?")
			data, err = tc.NativeToStringMap(nativeToRead)
			require.NoError(t, err)
			testData(t, tc.KeyField, id[toUpdateI], toUpdateResult, data, false, description, "on .Read() after Update()")

			toUpdate := StringMap{}
			for k, v := range toUpdateResult {
				toUpdate[k] = v
			}

			// can't duplicate uniques fields
			for _, key := range uniquesUpdatable {
				toUpdate[key] = toCreate[0][key]
				nativeToUpdate, err := tc.StringMapToNative(toUpdate)
				require.NoError(t, err)
				err = tc.Update(tc.ISToUpdate, nativeToUpdate)
				require.Error(t, err)
				toUpdate[key] = toUpdateResult[key]
			}

			// update 3: ok
			err = tc.Update(tc.ISToUpdate, nativeToUpdate)
			require.NoError(t, err, "what is the error on .Update()?")
			nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
			require.NoError(t, err, "what is the error on .Read() after Update()?")
			data, err = tc.NativeToStringMap(nativeToRead)
			require.NoError(t, err)
			testData(t, tc.KeyField, id[toUpdateI], toUpdateResult, data, false, description, "on .Read() after Update()")

			// can't update absent record
			toUpdate[tc.KeyField] += "123"
			nativeToUpdate, err = tc.StringMapToNative(toUpdate)
			require.NoError(t, err)
			err = tc.Update(tc.ISToUpdate, nativeToUpdate)
			require.Error(t, err)
		}

		// test DeleteList --------------------------------------------------------------------------------------

		if !tc.ExcludeDeleteTest {
			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
			require.NoError(t, err, "what is the error on .Read() after Update()?")
			data, err = tc.NativeToStringMap(nativeToRead)
			require.NoError(t, err)
			require.Equal(t, id[toDeleteI], data[tc.KeyField])

			if tc.ExpectedDeleteErr != nil {
				err = tc.Delete(tc.ISToDelete, id[toDeleteI])
				require.Error(t, err, "where is an error on .DeleteList()?")
				nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
				require.NoError(t, err, "what is the error on .Read() after Update()?")
				data, err = tc.NativeToStringMap(nativeToRead)
				require.NoError(t, err)
				require.Equal(t, id[toDeleteI], data[tc.KeyField])
				continue
			}

			if tc.ISToDeleteBad != nil {
				err = tc.Delete(*tc.ISToDeleteBad, id[toDeleteI])
				require.Error(t, err, "where is an error on .DeleteList()?")
				nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
				require.NoError(t, err, "what is the error on .Read() after Update()?")
				data, err = tc.NativeToStringMap(nativeToRead)
				require.NoError(t, err)
				require.Equal(t, id[toDeleteI], data[tc.KeyField])
			}

			err = tc.Delete(tc.ISToDelete, id[toDeleteI])
			require.NoError(t, err, "what is the error on .DeleteList()?")

			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])

			// it depends on implementation
			// require.Error(t, err, "where is an error on .Read() after DeleteList()?")

			require.Nil(t, nativeToRead)
		}
	}
}

func testData(t *testing.T, keyField, expectedID string, expectedData, data StringMap, onCreate bool, description Description, on string) {
	if expectedData == nil {
		require.Nil(t, data)
		return
	}
	require.NotNil(t, data)
	require.Equalf(t, expectedID, data[keyField], on+": incorrect key value in field '%s'???", keyField)
	for _, field := range description.Fields {
		key := field.Key

		// TODO: check key field

		if (onCreate && field.Creatable) || (!onCreate && field.Updatable) {
			if expectedData[key] == "" && field.NotEmpty {
				require.NotEmpty(t, data[key], on+": "+key+"???")
			} else {
				require.Equal(t, expectedData[key], data[key], on+": "+key+"???")
			}

		} else if !onCreate && field.Additable {
			if expectedData[key] == "" {
				require.Equal(t, expectedData[key], data[key], on+": "+key+"???")
			} else {
				require.True(t, len(data[key]) > len(expectedData[key]), on+": "+key+"???")
			}

		} else if field.NotEmpty {
			require.NotEmpty(t, data[key], on+": "+key+"???")
		}
	}
}
