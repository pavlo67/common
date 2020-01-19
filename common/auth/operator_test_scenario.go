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
				Cryptype: encrlib.NoCrypt,
				Values: Values{
					CredsPassword: "pass1",
					CredsNickname: testNick,
				},
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

		userKey, userCreds, err := tc.SetCreds(tc.UserKey, tc.ToSet, "")
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)

		require.Equal(t, tc.ToSet.Values[CredsNickname], userCreds.Values[CredsNickname])
		require.NotNil(t, userKey)

		// .Authorize() -----------------------------------------

		userCreds = &Creds{
			Cryptype: encrlib.NoCrypt,
			Values: Values{
				CredsIP:       testIP,
				CredsLogin:    tc.ToSet.Values[CredsNickname],
				CredsPassword: tc.ToSet.Values[CredsPassword],
			},
		}

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, tc.ToSet.Values[CredsNickname], user.Creds.Values[CredsNickname])
		require.Equal(t, userKey, user.Key)
	}
}

func OperatorTestScenarioToken(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userKey, userCreds, err := tc.SetCreds(testUserKey, tc.ToSet, CredsJWT)
		require.NoError(t, err)
		require.Equal(t, userKey, testUserKey)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)
		require.Equal(t, tc.ToSet.Values[CredsNickname], userCreds.Values[CredsNickname])

		// .Authorize() -----------------------------------------

		userCreds.Values[CredsIP] = testIP

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, tc.ToSet.Values[CredsNickname], user.Creds.Values[CredsNickname])
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

		userKey, userCreds, err := tc.SetCreds("", tc.ToSet, CredsPrivateKey)
		require.NoError(t, err)
		require.NotNil(t, userKey)
		require.NotNil(t, userCreds)

		log.Printf("   userKey, creds: %s, %#v", userKey, userCreds)

		require.NotEmpty(t, userKey)

		require.Equal(t, tc.ToSet.Values[CredsNickname], userCreds.Values[CredsNickname])
		nickname := userCreds.Values[CredsNickname]

		// .InitAuth() -----------------------------------

		privKeySerialization := []byte(userCreds.Values[CredsPrivateKey])
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := userCreds.Values[CredsPublicKeyBase58]
		log.Printf("public key base58: %s", publicKeyBase58)

		credsToInit := Creds{
			Cryptype: encrlib.NoCrypt,
			Values: Values{
				CredsIP: testIP,
			},
		}

		_, sessionCreds, err := tc.SetCreds("", credsToInit, CredsKeyToSignature)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		// ---------------------------------------------------------------------

		keyToSignature := sessionCreds.Values[CredsKeyToSignature]
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

		userCreds.Values[CredsIP] = testIP
		userCreds.Values[CredsNickname] = nickname
		userCreds.Values[CredsKeyToSignature] = keyToSignature
		userCreds.Values[CredsSignature] = string(signature)

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, nickname, user.Creds.Values[CredsNickname])
		require.Equal(t, userKey, user.Key)
	}
}
