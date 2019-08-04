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
