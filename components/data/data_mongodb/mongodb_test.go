package data_mongodb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/data"
)

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + env + ".yaml"
	cfg, err := config.Get(configPath, encodelib.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgMongoDB := config.Access{}
	err = cfg.Value("mongodb", &cfgMongoDB)
	require.NoError(t, err)

	dataOp, cleanerOp, mgoClient, err := NewData(&cfgMongoDB, 5*time.Second, "test", "crud", data.Item{Details: data.Test{}})
	require.NoError(t, err)

	testCases := data.TestCases(dataOp, cleanerOp)

	data.OperatorTestScenario(t, testCases, l)

	mgoClient.Disconnect(nil)
}
