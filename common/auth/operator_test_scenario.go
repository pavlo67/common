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
	ToSet  Creds
	ToInit Creds
}

const testIP = "1.2.3.4"

func TestCases(authOp Operator) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: authOp,
			ToSet: Creds{
				Cryptype: encrlib.NoCrypt,
				Values: Values{
					CredsNickname: "nick1",
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

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Info(i)

		userCreds, err := tc.SetCreds(nil, tc.ToSet)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		nickname := userCreds.Values[CredsNickname]
		identityKey := userCreds.Values[CredsIentityKey]

		if nicknameToSet := tc.ToSet.Values[CredsNickname]; nicknameToSet != "" {
			require.Equal(t, nicknameToSet, nickname)
		} else {
			require.NotEmpty(t, nickname)
		}

		sessionCreds, err := tc.InitAuthSession(tc.ToInit)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		privKeySerialization := []byte(userCreds.Values[CredsPrivateKey])
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := userCreds.Values[CredsPublicKeyBase58]
		log.Printf("         address: %s", publicKeyBase58)

		// ---------------------------------------------------------------------

		keyToSignature := sessionCreds.Values[CredsKeyToSignature]
		log.Printf("key to signature: %s", keyToSignature)
		require.True(t, len(keyToSignature) > 0)

		signature, err := encrlib.ECDSASign(keyToSignature, *privKey)
		require.NoError(t, err)
		require.True(t, len(signature) > 0)

		log.Printf("     private key: %s", privKeySerialization)
		log.Printf("       signature: %s", base58.Encode(signature))

		publKey := base58.Decode(publicKeyBase58)
		ok := encrlib.ECDSAVerify(keyToSignature, publKey, signature)
		require.True(t, ok)

		// ---------------------------------------------------------------------

		userCreds.Values[CredsIP] = testIP
		userCreds.Values[CredsKeyToSignature] = keyToSignature
		userCreds.Values[CredsSignature] = string(signature)

		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, identityKey, string(user.Key))
	}
}
