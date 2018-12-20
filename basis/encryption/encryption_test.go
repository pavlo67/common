package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/basis/strlib"
)

var passwordMinLength = 6
var testPassword = strlib.RandomString(passwordMinLength)
var testPasswordBad = strlib.RandomString(passwordMinLength - 1)

const testSalt = "$5$1234"
const testCryptype = SHA256
const testCryptypeBad1 = NoCrypt
const testCryptypeBad2 = Provos

func TestGetEncodedPassword(t *testing.T) {

	var err error
	var encodedPassword *Hash

	// too short password - error
	encodedPassword, err = GetEncodedPassword(testPasswordBad, []byte(testSalt), testCryptype, passwordMinLength, false)
	require.Error(t, err)
	require.Nil(t, encodedPassword)

	// bad cryptype - changed
	encodedPassword, err = GetEncodedPassword(testPassword, []byte(testSalt), testCryptypeBad1, passwordMinLength, false)
	require.NoError(t, err)
	require.NotNil(t, encodedPassword)
	require.Equal(t, testCryptype, encodedPassword.Cryptype)

	// another bad cryptype - changed
	encodedPassword, err = GetEncodedPassword(testPassword, []byte(testSalt), testCryptypeBad2, passwordMinLength, false)
	require.NoError(t, err)
	require.NotNil(t, encodedPassword)
	require.Equal(t, testCryptype, encodedPassword.Cryptype)

	// all ok
	encodedPassword, err = GetEncodedPassword(testPassword, []byte(testSalt), testCryptype, passwordMinLength, false)
	require.NoError(t, err)
	require.NotNil(t, encodedPassword)
	require.Equal(t, testCryptype, encodedPassword.Cryptype)

}
