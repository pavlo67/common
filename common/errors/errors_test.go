package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrappedError(t *testing.T) {
	testKey1 := Key("test_key1")
	ke1 := KeyableError(nil, testKey1, nil)
	require.Equalf(t, testKey1, ke1.Key(), "%#v", ke1)

	testKey2 := Key("test_key2")
	ke2 := KeyableError(errors.New("q"), testKey2, nil)
	require.Equalf(t, testKey2, ke2.Key(), "%#v", ke2)

	testKey3 := Key("test_key3")
	ke3 := KeyableError(New("q"), testKey3, nil)
	require.Equalf(t, testKey3, ke3.Key(), "%#v", ke3)

	testKey4 := Key("test_key4")
	ke4 := KeyableError(errors.New("q"), testKey4, nil)
	require.Equalf(t, testKey4, ke4.Key(), "%#v", ke4)

}
