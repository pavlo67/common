package packs_pg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"

	"github.com/pavlo67/workshop/components/packs"
)

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/gatherer." + env + ".yaml"
	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgPostgres := config.Access{}
	err = cfg.Value("postgres", &cfgPostgres)
	require.NoError(t, err)

	l.Infof("%#v", cfgPostgres)

	packsOp, cleanerOp, err := New(cfgPostgres, "", "")
	require.NoError(t, err)

	testCases := packs.TestCases(packsOp, cleanerOp)

	packs.OperatorTestScenario(t, testCases, l)
}
