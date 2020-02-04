package auth_jwt

import (
	"os"
	"testing"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/stretchr/testify/require"
)

const serviceName = "gatherer"

func TestOperator(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	authOp, err := New("key.test")
	require.NoError(t, err)
	require.NotNil(t, authOp)

	testCases := auth.TestCases(authOp)

	auth.OperatorTestScenarioToken(t, testCases, l)
}
