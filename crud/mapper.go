package crud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/basis/encrlib"
)

type Mapper interface {
	Fields() []Field
	Write(StringMap) (interface{}, error)
	Read(interface{}) (StringMap, error)
}

func MapperTest(t *testing.T, testCases []Mapper) {

	// TODO: NotEmpty-fields
	// TODO: check all fields with reflect

	for i, tc := range testCases {
		fmt.Println(i)

		fields := tc.Fields()

		data0 := StringMap{}

		for _, f := range fields {
			if f.Creatable {
				maxLength := f.MaxLength
				if maxLength == 0 {
					maxLength = 1
				}
				data0[f.Key] = encrlib.RandomString(maxLength)
			}
		}
		mapped, err := tc.Write(data0)
		require.NoError(t, err)

		data1, err := tc.Read(mapped)

		for k, v := range data0 {
			require.Equal(t, v, data1[k], fmt.Sprintf("??? for key: '%s'", k))
		}

	}
}
