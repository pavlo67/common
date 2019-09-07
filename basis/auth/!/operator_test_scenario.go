package a

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/partes/connector/receiver"
	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/basis/strlib"
)

var testSalt = "$5$1234"
var testPasswordMinLength = 6

var testDomain = "test.com"
var testID = strlib.RandomString(10)
var testNickname = "user_" + testID
var testEmail = testNickname + "@" + testDomain
var testPassword = strlib.RandomString(testPasswordMinLength)

type ToUseStep struct {
	ToUse, ToAuth Creds
	ToSet         []Creds
	ExpectedUser  *User
	ExpectedErr   error
}

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	ReceiverOp receiver.Operator

	ToCreate        []Creds
	ExpectCreateErr bool

	// ToUseSteps []ToUseStep
}

func TestRegistrationAndPasswordUpdating(t *testing.T, testCases []OperatorTestCase) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	var user *User

	for i, tc := range testCases {
		fmt.Println(i)

		// clear database ------------------------------------------------------------------------------

		if tc.Cleaner != nil {
			err := tc.Cleaner()
			require.NoError(t, err, "what is the error on .Cleaner()?")
		}

		// check the user to be created-----------------------------------------------------------------

		userToCreate := NewUserToCreate(tc.ToCreate)
		toUseNickname := Creds{
			Type:   CredsNickname,
			Values: []string{userToCreate.Nickname},
		}
		toUseEmail := Creds{
			Type:   CredsEmail,
			Values: []string{userToCreate.Email},
		}
		toAuth := Creds{
			Type:   CredsPassword,
			Values: []string{userToCreate.Password},
		}

		// authenication before registration - error ---------------------------------------------------

		user, _, _ = tc.Use(toUseNickname, toAuth)
		require.Nil(t, user)
		// require.Error(t, err) // it depends...

		user, _, _ = tc.Use(toUseEmail, toAuth)
		require.Nil(t, user)
		// require.Error(t, err) // it depends...

		// registration code - ok ----------------------------------------------------------------------

		listener, err := receiver.NewListener(tc.ReceiverOp)
		require.NoError(t, err)

		_, err = tc.Create(tc.ToCreate...)
		if tc.ExpectCreateErr {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)

		messages, err := listener.ReadNext()
		require.NoError(t, err)
		require.Equal(t, 1, len(messages))

		sentCode := strings.TrimSpace(messages[0].Body)
		require.True(t, len(sentCode) > 0)

		user, _, err = tc.Use(Creds{}, Creds{Type: CredsSentCode, Values: []string{sentCode}})
		require.NoError(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user.ID)

		if userToCreate.Nickname != "" {
			require.Equal(t, userToCreate.Nickname, user.Nickname)
		} else {
			require.NotEmpty(t, user.Nickname)
			toUseNickname.Values = []string{user.Nickname}
		}

		// nickname authenication after registration - ok ----------------------------------------------

		user, _, _ = tc.Use(toUseNickname, toAuth)
		require.NotNil(t, user)
		require.NoError(t, err)
		require.Equal(t, userToCreate.Nickname, user.Nickname)

		// email authenication after registration - ok -------------------------------------------------

		user, _, err = tc.Use(toUseEmail, toAuth)
		require.NotNil(t, user)
		require.NoError(t, err)
		require.Equal(t, userToCreate.Nickname, user.Nickname)

		// TODO: fails for registration with incorrect data or clones

	}

	//
	//// new password creation - ok
	//encodedPasswordNew, err := encrlib.GetEncodedPassword(testPassword+"new", []byte(testSalt), encrlib.testCryptype, passwordMinLength, false)
	//require.NoError(t, err)
	//require.NotNil(t, encodedPasswordNew)
	//
	//// password updating token - ok
	//err = credentialsOp.QueryCodeToUpdatePassword(user.Nick)
	//require.NoError(t, err, "wrong query password")
	//message, err = receiverOp.ReceiveNext() // !!! confirmestubsender
	//token = message.Body
	//
	//// password updating confirmation - ok
	//err = credentialsOp.UpdatePasswordWithCode(token, *encodedPasswordNew)
	//require.NoError(t, err)
	//
	//// authenication with old creds after password updating - error
	//user, err = credentialsOp.AuthenticateWithCreds(creds)
	//require.Nil(t, user)
	//require.Error(t, err)
	//
	//// authenication with new creds after password updating - ok
	//credsNew := Creds{testNickname, *encodedPasswordNew}
	//user, err = credentialsOp.AuthenticateWithCreds(credsNew)
	//require.NotNil(t, user)
	//require.NoError(t, err)
	//require.Equal(t, user0.Nick, user.Nick)
	//require.Equal(t, user0.Email, user.Email)
	//require.Equal(t, *encodedPasswordNew, *user.Hash)
	//
	//// clearing the database - ok
	//err = credentialsOp.Clean(nil)
	//require.NoError(t, err)
}

func TestCases(operator Operator, cleaner crud.Cleaner, ReceiverOp receiver.Operator) ([]OperatorTestCase, error) {

	testNickname := "testNickname"
	testEmail := "test@Email.com"
	testPassword := "testPassword"

	toUseNickname := Creds{
		Type:   CredsNickname,
		Values: []string{testNickname},
	}
	toUseEmail := Creds{
		Type:   CredsEmail,
		Values: []string{testEmail},
	}
	toAuth := Creds{
		Type:   CredsPassword,
		Values: []string{testPassword},
	}

	testCases := []OperatorTestCase{
		{
			Operator:        operator,
			Cleaner:         cleaner,
			ReceiverOp:      ReceiverOp,
			ToCreate:        []Creds{toUseNickname, toUseEmail, toAuth},
			ExpectCreateErr: false,
			// ToUseSteps:      nil,
		},
	}

	return testCases, nil
}
