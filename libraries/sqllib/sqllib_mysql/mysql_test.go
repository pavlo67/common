package sqllib_mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// import cycle not allowed - this file can't be placed in mysqllib directoty

func TestDesc(t *testing.T) {
	t.Skip()

	//	conf, err := config.ReadList(filelib.CurrentPath() + "../../cfg.json5")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	if conf == nil {
	//		log.Fatal(nil.New("no config data after setup.Prepare()"))
	//	}
	//	partKeys := config.PartKeys{
	//		"mysql": "items",
	//	}
	//	mysqlConfig, errs := conf.MySQL("", nil)
	//	err = errs
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	//err := config.LoadContext("../../../cfg.json5")
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//	//mysqlConfig, ok := config.Mysql["items"]
	//	//if !ok {
	//	//	log.Fatal(nil.Errorf("no mysql[items.comp] section in config: %v", config.Mysql))
	//	//}
	//
	//	dbh, err := ConnectToMysql(mysqlConfig)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	stmt, err := dbh.Init("desc user")
	//	if err != nil {
	//		log.Println("err .Init:", err)
	//	}
	//	defer stmt.Close()
	//
	//	rows, err := stmt.Query()
	//	if err != nil {
	//		log.Println("err .Query:", err)
	//	}
	//
	//	var a, b, c, d, f string
	//	var e []byte
	//
	//	for rows.Next() {
	//		err = rows.Scan(&a, &b, &c, &d, &e, &f)
	//		log.Println("a:", a, "b:", b, "c:", c, "d:", d, "e:", string(e), "f:", f, err)
	//	}
	//
	//}
	//
	//func TestCreateTable(t *testing.T) {
	//
	//	conf, err := config.ReadList(filelib.CurrentPath() + "../../cfg.json5")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	if conf == nil {
	//		log.Fatal(nil.New("no config data after setup.Prepare()"))
	//	}
	//	partKeys := config.PartKeys{
	//		"mysql": "items",
	//	}
	//	mysqlConfig, errs := conf.MySQL("", nil)
	//	err = errs
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	//err := config.LoadContext("../../../cfg.json5")
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//	//mysqlConfig, ok := config.Mysql["items"]
	//	//if !ok {
	//	//	log.Fatal(nil.Errorf("no mysql[items.comp] section in config: %v", config.Mysql))
	//	//}
	//
	//	dbh, err := ConnectToMysql(mysqlConfig)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	tableName := "test_create_table"
	//	testTable := []MySQLField{
	//		{"id", "int(10) unsigned", false, "", "auto_increment"},
	//		{"name", "varchar(100)", false, "", ""},
	//		{"r_view", "varchar(255)", false, "", ""},
	//		{"created_at", "timestamp", false, "CURRENT_TIMESTAMP", ""},
	//		{"updated_at", "timestamp", true, "", "on update CURRENT_TIMESTAMP"},
	//	}
	//	testIndex := []MySQLIndex{
	//		{"PRIMARY", "PRIMARY", []string{"id"}, ""},
	//		{"name", "UNIQUE", []string{"name", "r_view"}, ""},
	//		{"r_view", "", []string{"r_view"}, ""},
	//	}
	//
	//	dbh.Exec("drop table if exists " + tableName)
	//	// test not exists table
	//	err = GetTableExists(dbh, tableName)
	//	require.Equal(t, err, ErrNoTable, fmt.Sprintf("table: '%s' should not exist", tableName))
	//
	//	// test create table
	//	err = CreateTable(dbh, tableName, testTable, testIndex)
	//	require.Equal(t, err, nil, fmt.Sprintf("wrong create table: %s", tableName))
	//
	//	// test exists table
	//	err = GetTableExists(dbh, tableName)
	//	require.Equal(t, err, nil, fmt.Sprintf("table: '%s' must exist", tableName))
	//
	//	// test get table fields
	//	fields, err := GetMySQLTableFields(dbh, tableName)
	//	require.Equal(t, err, nil, fmt.Sprintf("wrong get fields from table: %s ", tableName))
	//	require.Equal(t, fields[1].Label, "name", fmt.Sprintf("wrong table field in %s (fields = %+v)", tableName, fields))
	//
	//	// test get table indexes
	//	indexes, err := GetMySQLTableIndexes(dbh, tableName)
	//	require.Equal(t, err, nil, fmt.Sprintf("wrong get indexes from table: %s ", tableName))
	//	require.Equal(t, indexes[0].Type, "PRIMARY", fmt.Sprintf("table %s indexes[0].Genus is incorrect (indexes: %+v)", tableName, indexes))
	//	require.Equal(t, indexes[1].Type, "UNIQUE", fmt.Sprintf("table %s indexes[1].Genus is incorrect (indexes: %+v)", tableName, indexes))
	//	require.Equal(t, indexes[2].Type, "", fmt.Sprintf("table %s indexes[2].Genus is incorrect (indexes: %+v)", tableName, indexes))
	//
	//	// test drop table
	//	err = DropTable(dbh, tableName)
	//	require.Equal(t, err, nil, "error drop table")
	//
}
