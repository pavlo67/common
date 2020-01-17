package auth

import (
	"log"
	"os"
	"testing"

	"github.com/btcsuite/btcutil/base58"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/libraries/encrlib"

	"github.com/pavlo67/workshop/common/logger"
)

type OperatorTestCase struct {
	Operator
	User   *User
	ToSet  Creds
	ToInit Creds
}

const testIP = "1.2.3.4"
const testNick = "nick1"

func TestCases(authOp Operator) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: authOp,

			User: &User{
				Key:      "nick1@aaa",
				Nickname: testNick,
			},

			ToSet: Creds{
				Cryptype: encrlib.NoCrypt,
				Values: Values{
					CredsPassword: "pass1",
				},
			},
			ToInit: Creds{
				Cryptype: encrlib.NoCrypt,
				Values: Values{
					CredsIP: testIP,
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

		userSet, userCreds, err := tc.SetCreds(tc.User, tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userSet)
		require.Nil(t, userCreds)

		log.Printf("            user: %#v", *userSet)

		require.Equal(t, tc.User.Nickname, userSet.Nickname)
		require.Equal(t, tc.User.Key, userSet.Key)

		// .InitAuthSession() -----------------------------------

		sessionCreds, err := tc.InitAuthSession(tc.ToInit)
		require.NoError(t, err)
		require.Nil(t, sessionCreds)

		// .Authorize() -----------------------------------------

		userCreds = &Creds{
			Cryptype: encrlib.NoCrypt,
			Values: Values{
				CredsIP:       testIP,
				CredsLogin:    tc.User.Nickname,
				CredsPassword: tc.ToSet.Values[CredsPassword],
			},
		}

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, userSet.Nickname, user.Nickname)
		require.Equal(t, userSet.Key, user.Key)
	}
}

func OperatorTestScenarioToken(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userSet, userCreds, err := tc.SetCreds(tc.User, tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userSet)
		require.NotNil(t, userCreds)

		log.Printf("            user: %#v", *userSet)

		require.Equal(t, tc.User.Nickname, userSet.Nickname)
		require.Equal(t, tc.User.Key, userSet.Key)

		// .InitAuthSession() -----------------------------------

		sessionCreds, err := tc.InitAuthSession(tc.ToInit)
		require.NoError(t, err)
		require.Nil(t, sessionCreds)

		// .Authorize() -----------------------------------------

		userCreds.Values[CredsIP] = testIP

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, userSet.Nickname, user.Nickname)
		require.Equal(t, userSet.Key, user.Key)
	}
}

func OperatorTestScenarioPublicKey(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userSet, userCreds, err := tc.SetCreds(tc.User, tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userSet)
		require.NotNil(t, userCreds)

		log.Printf("             user: %#v", *userSet)

		require.Equal(t, tc.User.Nickname, userSet.Nickname)
		require.NotEmpty(t, userSet.Key)

		// .InitAuthSession() -----------------------------------

		sessionCreds, err := tc.InitAuthSession(tc.ToInit)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		privKeySerialization := []byte(userCreds.Values[CredsPrivateKey])
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := userCreds.Values[CredsPublicKeyBase58]
		log.Printf("public key base58: %s", publicKeyBase58)

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
		userCreds.Values[CredsNickname] = userSet.Nickname
		userCreds.Values[CredsKeyToSignature] = keyToSignature
		userCreds.Values[CredsSignature] = string(signature)

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, userSet.Nickname, user.Nickname)
		require.Equal(t, userSet.Key, user.Key)
	}
}
