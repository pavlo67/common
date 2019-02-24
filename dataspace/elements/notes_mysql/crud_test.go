package notes_mysql

import (
	"log"
	"testing"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/confidenter/groups/groupsmysql"

	"github.com/pavlo67/punctum/notebook/links"
	"github.com/pavlo67/punctum/notebook/links/links_mysql"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/things_old/genera"
	"github.com/pavlo67/punctum/things_old/genera/generastub"
)

func TestCRUD(t *testing.T) {
	_, conf, err := joiner.Init(filelib.CurrentPath() + "../../../cfg.json5")

	if err != nil {
		t.Fatal(err)
	}
	if conf == nil {
		t.Fatal("no config data after setup.Init()")
	}

	mysqlConfig, errs := conf.MySQL("notebook", nil)
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	starters := []starter.Starter{
		{groupsmysql.Starter(), ""},
		{links_mysql.Starter(), ""},
		{generastub.Starter(), ""},
	}
	err = starter.Run(conf, starters, "TEST BUILD", false, false)
	if err != nil {
		log.Println(err)
	}

	info, err := Starter(false).Check(*conf, "./")
	if err != nil {
		t.Fatal(err, info)
	}

	ctrlOp, ok := joiner.Component(groups.InterfaceKey).(groups.Operator)
	if !ok {
		t.Fatalf("no groups.Operator found for objecstmysql")
	}

	linksOp, ok := joiner.Component(links.InterfaceKey).(links.Operator)
	if !ok {
		t.Fatalf("no tags.Operator found for objecstmysql")
	}

	generaOp, ok := joiner.Component(genera.InterfaceKey).(genera.Operator)
	if !ok {
		t.Fatalf("no genera.Operator found for objecstmysql")
	}

	objectMySQL, err := New(
		mysqlConfig,
		"object",
		false,
		ctrlOp,
		linksOp,
		generaOp,
		nil,
	)

	if err != nil {
		t.Fatalf("can't init notesMySQL for tests: %s", err)
	}

	operatorCRUD := notes.OperatorCRUD{objectMySQL}
	testCases, err := operatorCRUD.TestCases()
	if err != nil {
		t.Fatalf("can't operatorCRUD.TestCases(): %s", err)
	}

	crud.OperatorTest(t, testCases)

	starter.CloseAll()
}
