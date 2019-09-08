package auth_users

import (
	"testing"

	"os"

	"github.com/pavlo67/partes/connector/receiver"
	"github.com/pavlo67/partes/connector/senderreceiver"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/groups/groupsstub"
	"github.com/pavlo67/punctum/confidenter/users/userscrud"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/stretchr/testify/require"
)

func TestAuthUsers(t *testing.T) {

	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment set!!!")
	}

	_, conf, err := joiner.Init(filelib.CurrentPath() + "../../../cfg.json5")
	require.NotNil(t, conf)
	require.NoError(t, err)

	starters := []starter.Starter{
		{senderreceiver.Starter(), ""},
		{groupsstub.Starter(), ""},

		{userscrud.Starter(), ""},
		{Starter(false), ""},
	}

	err = starter.Run(conf, starters, "TEST AUTHSUSERS BUILD", false, false)

	receiverOp, ok := joiner.Component(receiver.InterfaceKey).(receiver.Operator)
	if !ok {
		t.Fatal("no receiver.Operator for test")
	}

	authOp, ok := joiner.Component(auth.InterfaceKey).(auth.Operator)
	if !ok {
		t.Fatal("no auth.Operator for test")
	}

	testCases, err := auth.TestCases(authOp, nil, receiverOp)
	require.NoError(t, err)

	auth.TestRegistrationAndPasswordUpdating(t, testCases)
}
