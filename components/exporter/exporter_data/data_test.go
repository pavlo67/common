package exporter_data

import (
	"os"
	"testing"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/components/exporter"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/data/data_pg"
	"github.com/pavlo67/workshop/components/storage"
)

const serviceName = "notebook"

func TestExporterData(t *testing.T) {
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

	label := "NOTEBOOK/PG EXPORTER TEST CLI BUILD"

	dataKey := storage.DataInterfaceKey
	cleanerKey := data.CleanerInterfaceKey

	starters := []starter.Starter{
		{data_pg.Starter(), common.Map{"table": storage.CollectionDefault, "interface_key": dataKey, "cleaner_key": cleanerKey, "no_tagger": true}},
		{Starter(), common.Map{"data_key": dataKey}},
	}

	joinerOp, err := starter.Run(starters, nil, cfg, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	exporterOp, ok := joinerOp.Interface(exporter.InterfaceKey).(exporter.Operator)
	require.True(t, ok)
	require.NotNil(t, exporterOp)

	dataOp, ok := joinerOp.Interface(dataKey).(data.Operator)
	require.True(t, ok)
	require.NotNil(t, dataOp)

	cleanerOp, ok := joinerOp.Interface(cleanerKey).(crud.Cleaner)
	require.True(t, ok)
	require.NotNil(t, cleanerOp)

	testCases := exporter.TestCases(exporterOp, dataOp, cleanerOp)

	exporter.OperatorTestScenario(t, testCases, l)
}
