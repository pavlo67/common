package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrappedError(t *testing.T) {
	testKey1 := ErrorKey("test_key1")
	ke1 := KeyableError(testKey1, nil, nil)
	we1 := WrappedError(ke1)
	require.Equalf(t, testKey1, we1.Key(), "%#v", we1)

	testKey2 := ErrorKey("test_key2")
	ke2 := KeyableError(testKey2, nil, errors.New("q"))
	we2 := WrappedError(ke2)
	require.Equalf(t, testKey2, we2.Key(), "%#v", we2)

	testKey3 := ErrorKey("test_key3")
	ke3 := KeyableError(testKey3, nil, errors.New("q"))
	we3 := WrappedError(ke3)
	require.Equalf(t, testKey3, we3.Key(), "%#v", we3)

	testKey4 := ErrorKey("test_key4")
	ke4 := KeyableError(testKey4, nil, errors.New("q"))
	we4 := WrappedError(ke4)
	require.Equalf(t, testKey4, we4.Key(), "%#v", we4)

}
