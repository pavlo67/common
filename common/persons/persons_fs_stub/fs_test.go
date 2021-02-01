package persons_fs_stub

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/common/common/logger"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/persons/persons_scenarios"
)

func TestOperator(t *testing.T) {

	os.Setenv("ENV", "test")

	var err error
	l, err = logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	personsOp, personsCleanerOp, err := New(config.Access{Path: "./test/"})
	require.NoError(t, err)

	persons_scenarios.OperatorTestScenario(t, personsOp, personsCleanerOp)
}
