package links_mysql

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/basis/filelib"
	"github.com/pavlo67/associatio/starter/config"
	"github.com/pavlo67/associatio/starter/joiner"
)

var conf *config.Config
var mysqlConfig config.ServerAccess

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalln("No test environment!!!")
	}

	var err error
	_, conf, err := joiner.Init(filelib.CurrentPath()+"../../../cfg.json5", false)
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal("no config data after setup.Run()")
	}

	var errs basis.Errors
	mysqlConfig, errs = conf.MySQL("notebook", nil)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	//starters := []starter.Starter{
	//	{groupsmysql.Starter(), "groupsmysql", ""},
	//}
	//
	//for _, c := range starters {
	//	log.Println("  -------------   check component: ", c.Nick, "   ---------------")
	//	info, err := c.Check(*conf, PartKeys, c.IndexPath)
	//	if err != nil {
	//		for _, i := range info {
	//			log.Println(i)
	//		}
	//		log.Fatalf("error calling Check() for component (%s): %s", c.Nick, err)
	//	}
	//	err = c.Run()
	//	if err != nil {
	//		log.Fatalf("error calling Run() for component (%s): %s", c.Nick, err)
	//	}
	//}
	//
	//info, err := Starter().Check(*conf, PartKeys, "./")
	//if err != nil {
	//	log.Fatal(err, info)
	//}

	os.Exit(m.Run())
}

//func TestQuery(t *testing.T) {
//	identity, identityAnother, managers, ctrlOp, err := groupsstub.IdentitiesForTestsOld()
//	if err != nil {
//		log.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
//	}
//	linksMySQL, err := New(mysqlConfig, linkTableDefault, ctrlOp, managers)
//	if err != nil {
//		t.Fatalf("can't init linksMySQL for tests: %s", err)
//	}
//
//	testCases := links.QueryTestCases(linksMySQL, identity, identityAnother)
//	links.QueryTest(t, testCases)
//}
//
//func TestQueryByObjectID(t *testing.T) {
//	identity, identityAnother, managers, ctrlOp, err := groupsstub.IdentitiesForTestsOld()
//	if err != nil {
//		log.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
//	}
//	linksMySQL, err := New(mysqlConfig, linkTableDefault, ctrlOp, managers)
//	if err != nil {
//		t.Fatalf("can't init linksMySQL for tests: %s", err)
//	}
//
//	testCases := links.QueryByObjectIDTestCases(linksMySQL, identity, identityAnother)
//	links.QueryByObjectIDTest(t, testCases)
//}
//
//func TestQueryByTag(t *testing.T) {
//	identity, identityAnother, managers, ctrlOp, err := groupsstub.IdentitiesForTestsOld()
//	if err != nil {
//		log.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
//	}
//	linksMySQL, err := New(mysqlConfig, linkTableDefault, ctrlOp, managers)
//	if err != nil {
//		t.Fatalf("can't init linksMySQL for tests: %s", err)
//	}
//
//	testCases := links.QueryByTagTestCases(linksMySQL, identity, identityAnother)
//	links.QueryByTagTest(t, testCases)
//}
//
//func TestQueryTags(t *testing.T) {
//	identity, identityAnother, managers, ctrlOp, err := groupsstub.IdentitiesForTestsOld()
//	if err != nil {
//		log.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
//	}
//	linksMySQL, err := New(mysqlConfig, linkTableDefault, ctrlOp, managers)
//	if err != nil {
//		t.Fatalf("can't init linksMySQL for tests: %s", err)
//	}
//
//	testCases := links.QueryTagsTestCases(linksMySQL, identity, identityAnother)
//	links.QueryTagsTest(t, testCases)
//}
//
//func TestQueryTagsByOwner(t *testing.T) {
//	identity, identityAnother, managers, ctrlOp, err := groupsstub.IdentitiesForTestsOld()
//	if err != nil {
//		log.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
//	}
//	linksMySQL, err := New(mysqlConfig, linkTableDefault, ctrlOp, managers)
//	if err != nil {
//		t.Fatalf("can't init linksMySQL for tests: %s", err)
//	}
//
//	testCases := links.QueryTagsByOwnerTestCases(linksMySQL, identity, identityAnother)
//	links.QueryTagsByOwnerTest(t, testCases)
//}
