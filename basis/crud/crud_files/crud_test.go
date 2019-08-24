package crud_files

//import (
//	"log"
//	"testing"
//
//	"github.com/pavlo67/partes/confidenter/groups"
//	"github.com/pavlo67/partes/confidenter/groups/groupsmysql"
//	"github.com/pavlo67/partes/things_old/generastub"
//	"github.com/pavlo67/workshop/crud"
//	"github.com/pavlo67/workshop/dataspace/links"
//	"github.com/pavlo67/workshop/dataspace/links/links_mysql"
//	"github.com/pavlo67/workshop/starter"
//	"github.com/pavlo67/workshop/starter/joiner"
//)
//
//func TestCRUD(t *testing.T) {
//
//	starters := []starter.Starter{
//		{groupsmysql.Starter(), ""},
//		{links_mysql.Starter(), ""},
//		{generastub.Starter(), ""},
//	}
//	err = starter.Run(conf, starters, "TEST BUILD", false, false)
//	if err != nil {
//		log.Println(err)
//	}
//
//	info, err := Starter(false).Check(*conf, "./")
//	if err != nil {
//		t.Fatal(err, info)
//	}
//
//	ctrlOp, ok := joiner.Component(groups.InterfaceKey).(groups.Operator)
//	if !ok {
//		t.Fatalf("no groups.Operator found for objecstmysql")
//	}
//
//	linksOp, ok := joiner.Component(links.InterfaceKey).(links.Operator)
//	if !ok {
//		t.Fatalf("no tags.Operator found for objecstmysql")
//	}
//
//	generaOp, ok := joiner.Component(genera.InterfaceKey).(genera.Operator)
//	if !ok {
//		t.Fatalf("no genera.Operator found for objecstmysql")
//	}
//
//	objectMySQL, err := New(
//		mysqlConfig,
//		"object",
//		false,
//		ctrlOp,
//		linksOp,
//		generaOp,
//		nil,
//	)
//
//	if err != nil {
//		t.Fatalf("can't init notesMySQL for tests: %s", err)
//	}
//
//	operatorCRUD := notes.OperatorCRUD{objectMySQL}
//	testCases, err := operatorCRUD.TestCases()
//	if err != nil {
//		t.Fatalf("can't operatorCRUD.TestCases(): %s", err)
//	}
//
//	crud.OperatorTest(t, testCases)
//
//	starter.CloseAll()
//}
//
//var tc = crud.OperatorTestCase{
//
//	ToSave: crud.StringMap{"a": "b1", "c": "d1"},
//	ToUpdate: crud.StringMap{"a": "b2", "c": "d2"},
//
//	//ExcludeListTest bool
//	//ExcludeUpdateTest   bool
//	//ExcludeRemoveTest   bool
//}
