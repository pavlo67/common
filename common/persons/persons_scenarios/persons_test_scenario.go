package persons_scenarios

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
)

func OperatorTestScenario(t *testing.T, personsOp persons.Operator, personsCleanerOp crud.Cleaner) {

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

	// prepare data --------------------------------------------

	dataToTest := common.Map{}
	passwordToTestWithID := "passwordToTestWithID"

	identityToTestWithID := auth.Identity{
		ID:       "test_id",
		Nickname: "test_nickname",
		Roles:    rbac.Roles{rbac.RoleUser},
		Creds:    auth.Creds{},
	}

	passwordToTestWithoutID := "passwordToTestWithoutID"

	identityToTestWithoutID := auth.Identity{
		Nickname: "test_nickname",
		Roles:    rbac.Roles{rbac.RoleUser},
		Creds:    auth.Creds{},
	}

	// add person with ID --------------------------------------

	personID, err := personsOp.Add(identityToTestWithID, dataToTest, nil)
	require.Error(t, err)
	require.Empty(t, personID)

	identityToTestWithID.Creds[auth.CredsPassword] = passwordToTestWithID
	personID, err = personsOp.Add(identityToTestWithID, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, identityToTestWithID.ID, personID)

	identityToTestWithID.Creds[auth.CredsPassword] = passwordToTestWithID
	personID, err = personsOp.Add(identityToTestWithID, dataToTest, adminOptions)
	require.Errorf(t, err, "%#v", err)
	require.Empty(t, personID)

	// add person without ID -----------------------------------

	identityToTestWithoutID.Creds[auth.CredsPassword] = passwordToTestWithoutID
	personID, err = personsOp.Add(identityToTestWithoutID, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personID)

	identityToTestWithoutID.Creds[auth.CredsPassword] = passwordToTestWithoutID
	personID, err = personsOp.Add(identityToTestWithoutID, dataToTest, adminOptions)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personID)

	// clean old data ------------------------------------------
	//  ------------------------------------------
	//  ------------------------------------------
	//  ------------------------------------------
	//  ------------------------------------------
	//  ------------------------------------------
	//  ------------------------------------------

}
