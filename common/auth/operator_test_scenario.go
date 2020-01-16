package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/libraries/encrlib"

	"github.com/pavlo67/workshop/common/logger"
)

type OperatorTestCase struct {
	Operator
	ToSet  Creds
	ToInit Creds
}

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
					CredsIP: "1.2.3.4",
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

		userCreds.Values[CredsKeyToSignature] = sessionCreds.Values[CredsKeyToSignature]
		user, err := tc.Authorize(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, identityKey, user.Key)
	}
}
