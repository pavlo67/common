package selectors

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/basis"
)

//func TestMain(m *testing.M) {
//	if err := os.Setenv("ENV", "test"); err != nil {
//		log.Fatalln("No test environment!!!")
//	}
//	os.Exit(m.Run())
//}
//
func TestSelectors(t *testing.T) {

	mySelector := And(
		Or(
			FieldEqual("b", 4),
			FieldEqual("c", 1, 2),
		),
		FieldEqual("a", 1, 2, 3),
		FieldEqual("z", nil),
	)
	//fmt.Println(mySelector)

	condition, values, err := Mysql(nil, mySelector)
	fmt.Println(condition, "\n", values, "\n", err)
	require.Equal(t, nil, err, "wrong test selectors")
	require.Equal(t, []interface{}([]interface{}{4, 1, 2, 1, 2, 3}), values, "wrong test selectors")
	require.Equal(t, "((`b` = ?) OR (`c` in (?,?))) AND (`a` in (?,?,?)) AND (`z` UserIS NULL)", condition, "wrong test selectors")

	myBoltSelector := And(
		Or(
			FieldEqual("b", 4),
			FieldEqual("c", 1, 2),
		),
		FieldEqual("a", 1, 2, 3),
	)
	data := map[string]string{"b": "0", "c": "2", "a": "2"}
	res, err := CheckMap(myBoltSelector, data)
	require.Equal(t, nil, err, "wrong test bolt selectors")
	log.Println("res1:", res)
	require.Equal(t, true, res, "wrong test bolt selectors")

	data = map[string]string{"b": "4", "c": "2", "a": "0"}
	res, err = CheckMap(myBoltSelector, data)
	require.Equal(t, nil, err, "wrong test bolt selectors")
	log.Println("res2:", res)
	require.Equal(t, false, res, "wrong test bolt selectors")

	mySelector = FieldEqual("id", "11", "22", "33")
	keys, newSelector, err := BoltPrepare(mySelector)
	require.Equal(t, nil, err, "wrong test bolt prepare selectors")
	require.Equal(t, [][]byte{basis.Uint64ToByte(11), basis.Uint64ToByte(22), basis.Uint64ToByte(33)}, keys, "wrong test prepare selectors")
	require.Equal(t, newSelector, nil, "wrong test bolt prepare selectors")

	mySelector = And(
		Or(
			FieldEqual("b", 4),
			FieldEqual("c", 1, 2),
		),
		FieldEqual("id", "1", "2", "3"),
		FieldEqual("z", nil),
	)
	filteredSelector := And(
		Or(
			FieldEqual("b", 4),
			FieldEqual("c", 1, 2),
		),
		FieldEqual("z", nil),
	)
	keys, newSelector, err = BoltPrepare(mySelector)
	require.Equal(t, nil, err, "wrong test bolt prepare selectors")
	log.Println(mySelector)
	log.Println(newSelector)
	require.Equal(t, [][]byte{basis.Uint64ToByte(1), basis.Uint64ToByte(2), basis.Uint64ToByte(3)}, keys, "wrong test prepare selectors")
	require.Equal(t, filteredSelector, newSelector, "wrong test bolt prepare selectors")

}
