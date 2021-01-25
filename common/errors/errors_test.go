package errors

import (
	"testing"

	"github.com/pavlo67/common/common"

	"github.com/stretchr/testify/require"
)

func TestWrappedError(t *testing.T) {
	testKey1 := Key("test_key1")
	ke1 := KeyableError(testKey1, nil)
	require.Equalf(t, testKey1, ke1.Key(), "%#v", ke1)

	testKey2 := Key("test_key2")
	ke2 := KeyableError(testKey2, common.Map{"error": "q"})
	require.Equalf(t, testKey2, ke2.Key(), "%#v", ke2)

	testKey3 := Key("test_key3")
	ke3 := KeyableError(testKey3, common.Map{"error": "q"})
	require.Equalf(t, testKey3, ke3.Key(), "%#v", ke3)

	testKey4 := Key("test_key4")
	ke4 := KeyableError(testKey4, common.Map{"error": "q"})
	require.Equalf(t, testKey4, ke4.Key(), "%#v", ke4)

}
