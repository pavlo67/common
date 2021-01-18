package auth_ecdsa

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
)

const serviceName = ""

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

	authOp, err := New()
	require.NoError(t, err)
	require.NotNil(t, authOp)

	auth.OperatorTestScenarioPublicKey(t, authOp, l)
}
