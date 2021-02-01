package persons_fs

import (
	"testing"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/starter"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/persons/persons_scenarios"
)

func TestOperator(t *testing.T) {

	_, cfgService := apps.PrepareTests(t, "test_service", "../../../"+apps.AppsSubpathDefault, "test")

	components := []starter.Starter{
		{Starter(), nil},
	}

	label := "PERSONS_FS/TEST BUILD"
	joinerOp, err := starter.Run(components, cfgService, label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(crud.Cleaner)
	require.NotNil(t, personsCleanerOp)

	persons_scenarios.OperatorTestScenario(t, personsOp, personsCleanerOp)
}
