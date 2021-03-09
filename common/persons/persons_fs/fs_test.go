package persons_fs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/starter"
)

func TestPersonsFSAsPersonsOperator(t *testing.T) {
	_, cfgService, l := apps.PrepareTests(
		t,
		"../../../"+apps.AppsSubpathDefault,
		"test",
		"", // "persons_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	components := []starter.Starter{
		{Starter(), nil},
	}

	label := "PERSONS_FS/TEST BUILD"
	joinerOp, err := starter.Run(components, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(crud.Cleaner)
	require.NotNil(t, personsCleanerOp)

	persons.OperatorTestScenario(t, personsOp, personsCleanerOp)
}
