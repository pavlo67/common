package data_pg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"

	"github.com/pavlo67/workshop/components/data"
)

const serviceName = "gatherer"

func TestCRUD(t *testing.T) {
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

	//taggerOp, taggerCleanerOp, err := tagger_sqlite.New(cfgPg, "")
	//require.NoError(t, err)

	dataOp, cleanerOp, err := New(cfgPg, "storage", "", nil, nil) // taggerOp, taggerCleanerOp
	require.NoError(t, err)
	require.NotNil(t, dataOp)
	require.NotNil(t, cleanerOp)

	testCases := data.TestCases(dataOp, cleanerOp)

	data.OperatorTestScenario(t, testCases, l)
}
