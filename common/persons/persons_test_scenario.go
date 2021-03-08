package persons

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/rbac"
)

func OperatorTestScenario(t *testing.T, personsOp Operator, personsCleanerOp crud.Cleaner) {

	// prepare... ----------------------------------------------

	require.Equal(t, "test", os.Getenv("ENV"))

	require.NotNil(t, personsOp)
	require.NotNil(t, personsCleanerOp)

	adminOptions := crud.OptionsWithRoles(rbac.RoleAdmin)
	require.NotNil(t, adminOptions)

	// clean old data ------------------------------------------

	err := personsCleanerOp.Clean(adminOptions)
	require.NoError(t, err)

	personItems, err := personsOp.List(adminOptions)
	require.NoError(t, err)
	require.Equal(t, 0, len(personItems))

	// add person with ID --------------------------------------

	dataToTest := common.Map{}
	passwordToTestWithID := "passwordToTestWithID"

	identityToTestWithID := auth.Identity{
		ID:       "test_id",
		Nickname: "test_nickname1",
		Roles:    rbac.Roles{rbac.RoleUser},
	}

	personIDWrong, err := personsOp.Add(identityToTestWithID, nil, dataToTest, nil)
	require.Error(t, err)
	require.Empty(t, personIDWrong)

	personID1, err := personsOp.Add(identityToTestWithID, auth.Creds{auth.CredsPassword: passwordToTestWithID}, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, identityToTestWithID.ID, personID1)

	person1, err := personsOp.Read(personID1, adminOptions)

	require.NoErrorf(t, err, "%#v", err)
	require.True(t, person1.CheckCreds(auth.CredsPassword, passwordToTestWithID))

	person1.SetCreds(auth.Creds{auth.CredsPassword: ""})
	require.Equal(t, identityToTestWithID, person1.Identity)

	person1Options := crud.Options{Identity: &person1.Identity}

	personIDWrong, err = personsOp.Add(identityToTestWithID, auth.Creds{auth.CredsPassword: passwordToTestWithID}, dataToTest, adminOptions)
	require.Errorf(t, err, "%#v", err)
	require.Empty(t, personIDWrong)

	// add person without ID -----------------------------------

	passwordToTestWithoutID := "passwordToTestWithoutID"

	identityToTestWithoutID := auth.Identity{
		Nickname: "test_nickname2",
		Roles:    rbac.Roles{rbac.RoleUser},
	}

	personID2, err := personsOp.Add(identityToTestWithoutID, auth.Creds{auth.CredsPassword: passwordToTestWithoutID}, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personID2)

	personID3, err := personsOp.Add(identityToTestWithoutID, auth.Creds{auth.CredsPassword: passwordToTestWithoutID}, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personID3)

	person2, err := personsOp.Read(personID2, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, identityToTestWithoutID.Nickname, person2.Nickname)
	require.Equal(t, personID2, person2.ID)

	person2Options := crud.Options{Identity: &person2.Identity}

	person3, err := personsOp.Read(personID3, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, identityToTestWithoutID.Nickname, person3.Nickname)
	require.Equal(t, personID3, person3.ID)

	// list persons by admin: ok -------------------------------

	personItems, err = personsOp.List(adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, 3, len(personItems))

	// list persons by itself: error ---------------------------

	personItems, err = personsOp.List(&person1Options)
	require.Errorf(t, err, "%#v", err)
	require.Empty(t, personItems)

	// change person by admin: ok ------------------------------

	person1ToChange := *person1
	person1ToChange.Nickname += "_changed"

	person1Changed, err := personsOp.Change(person1ToChange, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, person1ToChange.Identity, person1Changed.Identity)

	person1ChangedReaded, err := personsOp.Read(person1Changed.ID, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, person1ToChange.Identity, person1ChangedReaded.Identity)

	// change person by itself: ok -----------------------------

	person1ToChange.Nickname += "_again"

	person1Changed, err = personsOp.Change(person1ToChange, &person1Options)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, person1ToChange.Identity, person1Changed.Identity)

	person1ChangedReaded, err = personsOp.Read(person1Changed.ID, &person1Options)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, person1ToChange.Identity, person1ChangedReaded.Identity)

	// change/read person by another person: error -------------

	person1ToChangeAgain := *person1ChangedReaded
	person1ToChangeAgain.Nickname += "_again2"

	person1ChangedWrong, err := personsOp.Change(person1ToChangeAgain, &person2Options)
	require.Errorf(t, err, "%#v", err)
	require.Nil(t, person1ChangedWrong)

	person1ReadedWrong, err := personsOp.Read(personID1, &person2Options)
	require.Errorf(t, err, "%#v", err)
	require.Nil(t, person1ReadedWrong)

	person1Readed, err := personsOp.Read(personID1, &person1Options)
	require.NoErrorf(t, err, "%#v", err)
	require.NotNil(t, person1Readed)
	require.Equal(t, person1Changed.Identity, person1Readed.Identity)

	// remove person by admin: ok ------------------------------

	err = personsOp.Remove(personID3, adminOptions)
	require.NoErrorf(t, err, "%#v", err)

	person3Readed, err := personsOp.Read(personID3, adminOptions)
	require.Errorf(t, err, "%#v", err)
	require.Nil(t, person3Readed)

	// remove person by itself: ok -----------------------------

	require.NotNil(t, person2Options.Identity)
	err = personsOp.Remove(personID2, &person2Options)
	require.NoErrorf(t, err, "%#v / %#v", person2Options.Identity, err)

	person2Readed, err := personsOp.Read(personID2, &person2Options)
	require.Errorf(t, err, "%#v", err)
	require.Nil(t, person2Readed)

	// remove person by another person: error ------------------

	err = personsOp.Remove(personID1, &person2Options)
	require.Errorf(t, err, "%#v", err)

	person1Readed, err = personsOp.Read(personID1, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.NotNil(t, person1Readed)
	require.Equal(t, person1ChangedReaded.Identity, person1Readed.Identity)

	// list persons by admin: ok -------------------------------

	personItems, err = personsOp.List(adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, 1, len(personItems))

	// clean old data ------------------------------------------
	// TODO???

}
