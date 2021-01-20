package auth

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/libraries/encrlib"
	"github.com/pavlo67/workshop/common/logger"
)

type OperatorTestCase struct {
	UserID ID
	ToSet  Creds
}

const testIP = "1.2.3.4"

//var testCases = []OperatorTestCase{
//	{
//		ToSet: Creds{
//			CredsPassword: "pass1",
//			CredsNickname: testNick,
//		},
//	},
//}

//const testUserKey = Key("nick1@aaa")
//
//func TestCases(authOp Operator) []OperatorTestCase {
//	return []OperatorTestCase{
//		{
//			Operator: authOp,
//
//			ToSet: Creds{
//				CredsPassword: "pass1",
//				CredsNickname: testNick,
//			},
//		},
//	}
//}
//
//func OperatorTestScenarioPassword(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
//	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
//		t.Fatal("No test environment!!!")
//	}
//
//	for i, tc := range testCases {
//		l.Info(i)
//
//		// .SetCreds() ------------------------------------------
//
//		userCreds, err := operator.SetCreds("", tc)
//		require.NoError(t, err)
//		require.NotNil(t, userCreds)
//
//		log.Printf("           creds: %#v", *userCreds)
//
//		require.Equal(t, tc[CredsNickname], (*userCreds)[CredsNickname])
//
//		// .Authenticate() ok -----------------------------------------
//
//		userCreds = &Creds{
//			CredsIP:       testIP,
//			CredsLogin:    tc[CredsNickname],
//			CredsPassword: tc[CredsPassword],
//		}
//
//		user, err := operator.Authenticate(*userCreds)
//
//		require.NoError(t, err)
//		require.NotNil(t, user)
//		require.Equal(t, tc[CredsNickname], user.Creds[CredsNickname])
//		require.NotEmpty(t, user.Key)
//
//		// .Authenticate() err ----------------------------------------
//
//		userCreds = &Creds{
//			CredsIP:       testIP,
//			CredsLogin:    tc[CredsNickname],
//			CredsPassword: tc[CredsPassword] + "1",
//		}
//
//		user, err = operator.Authenticate(*userCreds)
//
//		require.Error(t, err)
//		require.Nil(t, user)
//	}
//}
//
func OperatorTestScenarioToken(t *testing.T, operator Operator, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	testCreds := []Creds{
		{
			CredsNickname: "nickname" + strconv.FormatInt(time.Now().Unix(), 10),
		},
	}

	for i, tc := range testCreds {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		userCreds, err := operator.SetCreds("", tc)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)
		require.Equal(t, tc[CredsNickname], (*userCreds)[CredsNickname])

		// .Authenticate() -----------------------------------------

		(*userCreds)[CredsIP] = testIP

		user, err := operator.Authenticate(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, tc[CredsNickname], user.Nickname)
	}
}

func OperatorTestScenarioPublicKey(t *testing.T, operator Operator, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	testCreds := []Creds{
		{},
	}

	for i, tc := range testCreds {
		l.Info(i)

		// .SetCreds() ------------------------------------------

		// tc[CredsToSet] = CredsPrivateKey

		userCreds, err := operator.SetCreds("", tc)
		require.NoError(t, err)
		require.NotNil(t, userCreds)
		require.NotEmpty(t, (*userCreds)[CredsPublicKeyBase58])
		require.NotEmpty(t, (*userCreds)[CredsPublicKeyBase58], (*userCreds)[CredsNickname])

		log.Printf("            creds: %#v", userCreds)

		// require.Equal(t, tc[CredsNickname], userCreds.StringDefault(CredsNickname, ""))
		// nickname := (*userCreds)[CredsNickname]

		// .InitAuth() -----------------------------------

		privKeySerialization := []byte(userCreds.StringDefault(CredsPrivateKey, ""))
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := userCreds.StringDefault(CredsPublicKeyBase58, "")
		log.Printf("public key base58: %s", publicKeyBase58)

		credsToSet := Creds{CredsToSet: CredsKeyToSignature} // CredsIP: testIP,
		sessionCreds, err := operator.SetCreds("", credsToSet)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		// ---------------------------------------------------------------------

		keyToSignature := testIP

		//log.Printf(" key to signature: %s", keyToSignature)
		//require.True(t, len(keyToSignature) > 0)

		signature, err := encrlib.ECDSASign(keyToSignature, *privKey)
		require.NoError(t, err)
		require.True(t, len(signature) > 0)

		log.Printf("      private key: %s", privKeySerialization)
		log.Printf("        signature: %s", base58.Encode(signature))

		publKey := base58.Decode(publicKeyBase58)
		ok := encrlib.ECDSAVerify(keyToSignature, publKey, signature)
		require.True(t, ok)

		// .Authenticate() -----------------------------------------

		(*userCreds)[CredsIP] = testIP
		(*userCreds)[CredsKeyToSignature] = keyToSignature
		(*userCreds)[CredsSignature] = string(signature)

		user, err := operator.Authenticate(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		// require.Equal(t, nickname, user.Creds[CredsNickname])
		require.NotEmpty(t, user.ID)
	}
}
