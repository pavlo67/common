package sqllib

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pavlo67/constructor/starter/config"
)

func SetupTables(sqlOp Operator, tablesConfig map[string]config.SQLTable) error {

	db, err := sqlOp.DB()
	if err != nil {
		return err
	}

	tablesConfig = PrepareTables(tablesConfig)

	for _, table := range tablesConfig {
		err := DropTable(db, table.Name)
		if err != nil {
			return err
		}
		log.Println("table '" + table.Name + "' is dropped ")

		sqlQuery, err := sqlOp.CreateSQLQuery(table)
		if err != nil {
			return err
		}

		_, err = Exec(db, sqlQuery)
		if err != nil {
			return err
		}

		log.Println("table '" + table.Name + "' is created")
	}
	return nil
}

func PrepareTables(tablesConfig map[string]config.SQLTable) map[string]config.SQLTable {

	for key, table := range tablesConfig {

		table.Fields = []config.SQLField{}

		for _, f := range table.FieldsArr {
			n := false
			if f[2] == "true" {
				n = true
			}

			t := strings.Trim(f[1], " \n\r")
			def := strings.Trim(f[3], " \n\r")
			extra := strings.Trim(f[4], " \n\r")

			table.Fields = append(table.Fields, config.SQLField{f[0], t, n, def, extra})

		}

		tablesConfig[key] = table

	}

	return tablesConfig
}

//func DropTableIndex(dbh *sql.DB, table, indexName string) error {
//
//	var stmt *sql.Stmt
//	var err error
//	if strings.ToUpper(indexName) == "PRIMARY" {
//		return errors.New("reindex does not change PRIMARY index ")
//	}
//	sqlQuery := "alter table `" + table + "` drop index `" + indexName + "`"
//	if stmt, err = dbh.Prepare(sqlQuery); err != nil {
//		return errors.Wrapf(err, "can't prepare sql: %v", sqlQuery)
//	}
//	defer stmt.Close()
//
//	if _, err = stmt.Exec(); err != nil {
//		return errors.Wrapf(err, "can't exec sql: %v, values=%v", sqlQuery, indexName)
//	}
//	log.Println(table, ": drop index `"+indexName+"`")
//	return nil
//}
//
//func AddTableIndex(dbh *sql.DB, table, indexName, indexType string, indexFields []string) error {
//
//	var stmt *sql.Stmt
//	sqlQuery := "alter table `" + table + "` add "
//	if strings.ToUpper(indexType) == "PRIMARY" {
//		return errors.New("reindex does not change PRIMARY index ")
//	}
//	if strings.ToUpper(indexType) == "UNIQUE" || indexType == "" || strings.ToUpper(indexType) == "FULLTEXT" {
//		if indexType == strings.ToUpper("UNIQUE") {
//			sqlQuery += "UNIQUE KEY `" + indexName + "` ("
//		} else if strings.ToUpper(indexType) == "FULLTEXT" {
//			sqlQuery += "FULLTEXT KEY `" + indexName + "` ("
//		} else {
//			sqlQuery += "INDEX `" + indexName + "` USING BTREE ("
//		}
//		liF := 0
//		for _, f := range indexFields {
//			if liF > 0 {
//				sqlQuery += ", "
//			}
//			sqlQuery += "`" + f + "`"
//			liF++
//		}
//		sqlQuery += ")"
//	}
//	if err := Exec(dbh, sqlQuery, &stmt); err != nil {
//		return err
//	}
//	if _, err := stmt.Exec(); err != nil {
//		return errors.Wrapf(err, "can't exec sql:%v", sqlQuery)
//	}
//	log.Println(table, ": add index `"+indexName+"`: ", sqlQuery)
//	return nil
//}
