package auth_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/common/users/users_stub"
)

const serviceName = "notebook"

func TestOperator(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + serviceName + "." + env + ".yaml"
	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgPg := config.Access{}
	err = cfg.Value("pg", &cfgPg)
	require.NoError(t, err)

	l.Infof("%#v", cfgPg)

	label := "NOTEBOOK/STUB CLI BUILD"

	starters := []starter.Starter{
		{users_stub.Starter(), nil},

		{Starter(), nil},
	}

	joinerOp, err := starter.Run(starters, nil, nil, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	authOp, ok := joinerOp.Interface(InterfaceKey).(auth.Operator)
	require.True(t, ok)
	require.NotNil(t, authOp)

	testCases := auth.TestCases(authOp)

	auth.OperatorTestScenarioPassword(t, testCases, l)
}
