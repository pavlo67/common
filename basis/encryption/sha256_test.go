package encryption

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

var passwordMinLength = 6
var testPassword = RandomString(passwordMinLength)
var testPasswordBad = RandomString(passwordMinLength - 1)

const testSalt = "$5$1234"
const testCryptype = SHA256
const testCryptypeBad1 = NoCrypt
const testCryptypeBad2 = Provos

func TestGetEncodedPassword(t *testing.T) {

	// too short password - error
	encodedPassword, err := SHA256Hash(testPasswordBad, testSalt)
	require.NoError(t, err)

	fmt.Print(encodedPassword)

}
