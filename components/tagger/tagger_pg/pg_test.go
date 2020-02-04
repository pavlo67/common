package tagger_pg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/components/tagger"
)

type Test struct {
	AAA string
	BBB int
}

const serviceName = "gatherer"

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../apps/_environments/" + serviceName + "." + env + ".yaml"
	cfg, err := config.Get(configPath, serviceName, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgPg := config.Access{}
	err = cfg.Value("pg", &cfgPg)
	require.NoError(t, err)

	l.Infof("%#v", cfgPg)

	taggerOp, cleanerOp, err := New(cfgPg, tagger.InterfaceKey)
	require.NoError(t, err)

	l.Debugf("%#v", taggerOp)

	testCases := tagger.QueryTagsTestCases(taggerOp)

	tagger.OperatorTestScenario(t, testCases, cleanerOp, l)
}
