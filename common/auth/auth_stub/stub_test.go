package auth_stub

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
)

const serviceName = "gatherer"

func TestOperator(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	//configPath := filelib.CurrentPath() + "../../../environments/" + serviceName + "." + env + ".yaml"
	//cfg, err := config.Get(configPath, serializer.MarshalerYAML)
	//require.NoError(t, err)
	//require.NotNil(t, cfg)

	//salt := &common.Salt{
	//	SaltLenMin:    1,
	//	SaltLenMax:    100,
	//	RoundsMin:     1,
	//	RoundsMax:     100,
	//	RoundsDefault: 10,
	//}
	//saltStr := string(salt.Generate(10))
	//
	//l.Infof("salt string: %s", saltStr)

	authOp, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, authOp)

	testCases := auth.TestCases(authOp)

	auth.OperatorTestScenarioPassword(t, testCases, l)
}
