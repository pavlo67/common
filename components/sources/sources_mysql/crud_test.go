package sources_mysql

import (
	"testing"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/joiner"

	"os"

	"github.com/stretchr/testify/require"
)

func TestCRUD(t *testing.T) {

	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment set!!!")
	}

	_, conf, err := joiner.Init(filelib.CurrentPath() + "../../../cfg.json5")
	if err != nil {
		require.NoError(t, err)
	}

	mysqlConfig, errs := conf.MySQL("processor", nil)
	require.NoError(t, errs.Err())

	//starters := []starter.Starter{
	//	//{groupsmysql.Starter(), ""},
	//}
	//err = starter.Run(conf, starters, "TEST BUILD", false, false)
	//if err != nil {
	//	t.Fatal(err)
	//}

	srcOp, err := New(
		nil,
		mysqlConfig,
		"sources",
		nil,
	)
	require.NoError(t, err)

	operatorCRUD := sources.OperatorCRUD{srcOp}
	testCases, err := operatorCRUD.TestCases(func() error { return srcOp.clean() })
	require.NoError(t, err)

	crud.OperatorTest(t, testCases)

	starter.CloseAll()

}

//info, err := Starter(false).Check(*conf, "./")
//if err != nil {
//	t.Fatal(err, info)
//}
//
//ctrlOp, ok := joiner.Component(groups.InterfaceKey).(groups.Operator)
//if !ok {
//	t.Fatalf("no groups.Operator found for objecstmysql")
//}
//
//linksOp, ok := joiner.Component(links.InterfaceKey).(links.Operator)
//if !ok {
//	t.Fatalf("no tags.Operator found for objecstmysql")
//}
//
//generaOp, ok := joiner.Component(genera.InterfaceKey).(genera.Operator)
//if !ok {
//	t.Fatalf("no genera.Operator found for objecstmysql")
//}
//
//objectMySQL, err := New(
//	mysqlConfig,
//	"object",
//	false,
//	ctrlOp,
//	linksOp,
//	generaOp,
//	nil,
//)
//
//if err != nil {
//	t.Fatalf("can't init objectsMySQL for tests: %s", err)
//}
