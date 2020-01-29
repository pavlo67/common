package auth

import (
	"log"
	"os"
	"testing"

	"github.com/btcsuite/btcutil/base58"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/libraries/encrlib"
	"github.com/pavlo67/workshop/common/logger"
)

type OperatorTestCase struct {
	Operator
	UserKey identity.Key
	ToSet   Creds
}

const testIP = "1.2.3.4"
const testNick = "nick1"
const testUserKey = identity.Key("nick1@aaa")

func TestCases(authOp Operator) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: authOp,

			ToSet: Creds{
				CredsPassword: "pass1",
				CredsNickname: testNick,
			},
		},
	}
}

func OperatorTestScenarioPassword(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userCreds, err := tc.SetCreds("", tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)

		require.Equal(t, tc.ToSet[CredsNickname], (*userCreds)[CredsNickname])

		// .Authorize() ok -----------------------------------------

		userCreds = &Creds{
			CredsIP:       testIP,
			CredsLogin:    tc.ToSet[CredsNickname],
			CredsPassword: tc.ToSet[CredsPassword],
		}

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, tc.ToSet[CredsNickname], user.Creds[CredsNickname])
		require.NotEmpty(t, user.Key)

		// .Authorize() err ----------------------------------------

		userCreds = &Creds{
			CredsIP:       testIP,
			CredsLogin:    tc.ToSet[CredsNickname],
			CredsPassword: tc.ToSet[CredsPassword] + "1",
		}

		user, err = tc.Authorize(*userCreds)

		require.Error(t, err)
		require.Nil(t, user)
	}
}

func OperatorTestScenarioToken(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userCreds, err := tc.SetCreds(testUserKey, tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)
		require.Equal(t, tc.ToSet[CredsNickname], (*userCreds)[CredsNickname])

		// .Authorize() -----------------------------------------

		(*userCreds)[CredsIP] = testIP

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, tc.ToSet[CredsNickname], user.Creds[CredsNickname])
		require.Equal(t, testUserKey, user.Key)
	}
}

func OperatorTestScenarioPublicKey(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		tc.ToSet[CredsToSet] = string(CredsPrivateKey)

		userCreds, err := tc.SetCreds("", tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("            creds: %#v", userCreds)

		require.Equal(t, tc.ToSet[CredsNickname], (*userCreds)[CredsNickname])
		nickname := (*userCreds)[CredsNickname]

		// .InitAuth() -----------------------------------

		privKeySerialization := []byte((*userCreds)[CredsPrivateKey])
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := (*userCreds)[CredsPublicKeyBase58]
		log.Printf("public key base58: %s", publicKeyBase58)

		credsToInit := Creds{CredsIP: testIP, CredsToSet: string(CredsKeyToSignature)}

		sessionCreds, err := tc.SetCreds("", credsToInit)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		// ---------------------------------------------------------------------

		keyToSignature := (*sessionCreds)[CredsKeyToSignature]
		log.Printf(" key to signature: %s", keyToSignature)
		require.True(t, len(keyToSignature) > 0)

		signature, err := encrlib.ECDSASign(keyToSignature, *privKey)
		require.NoError(t, err)
		require.True(t, len(signature) > 0)

		log.Printf("      private key: %s", privKeySerialization)
		log.Printf("        signature: %s", base58.Encode(signature))

		publKey := base58.Decode(publicKeyBase58)
		ok := encrlib.ECDSAVerify(keyToSignature, publKey, signature)
		require.True(t, ok)

		// .Authorize() -----------------------------------------

		(*userCreds)[CredsIP] = testIP
		(*userCreds)[CredsNickname] = nickname
		(*userCreds)[CredsKeyToSignature] = keyToSignature
		(*userCreds)[CredsSignature] = string(signature)

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, nickname, user.Creds[CredsNickname])
		require.NotEmpty(t, user.Key)
	}
}
