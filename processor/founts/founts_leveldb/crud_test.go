package founts_leveldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/processor/founts"
)

func TestCRUD(t *testing.T) {

	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment set!!!")
	}

	//err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	//require.NoError(t, err)
	//
	//cfgPath := filelib.CurrentPath() + "../../../cfg.json5"
	//conf, err := config.Get(cfgPath, l)
	//require.NoError(t, err)
	//
	//starters := []starter.Starter{{Starter(), nil}}
	//joiner, err := starter.Run(conf, starters, "FOUNTS LEVELDB TEST", nil)
	//require.NoError(t, err)
	//
	//defer joiner.CloseAll()
	//
	//fountsOp, _ := joiner.Interface(founts.InterfaceKey).(founts.Operator)
	//require.NotNil(t, fountsOp)
	//

	fountsOp, err := New("test")

	operatorCRUD := founts.OperatorCRUD{fountsOp}
	testCases, err := operatorCRUD.TestCases(func() error { return fountsOp.clean() })
	require.NoError(t, err)

	crud.OperatorTest(t, testCases)

	err = fountsOp.Close()
	require.NoError(t, err)

}
