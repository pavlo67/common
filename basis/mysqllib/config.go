package mysqllib

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"database/sql"

	"github.com/pavlo67/associatio/starter/config"
)

func prepareTableComponents(mysqlConfig config.ServerAccess, tablesConfig map[string]config.MySQLTableComponent) (*sql.DB, map[string][]config.MySQLField, map[string][]config.MySQLIndex, error) {
	dbh, err := ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error connect to mySQL")
	}

	version, err := MySQLVersion(dbh)

	tables := map[string][]config.MySQLField{}
	indexes := map[string][]config.MySQLIndex{}

	for key, t := range tablesConfig {
		tables[key] = SetTableFields(t.Fields, version)
		if version < MySQLVersionFullTextLimit {
			// ignore FULLTEXT index for old mySQL version
			for _, in := range t.Indexes {
				if strings.ToUpper(in.Type) != "FULLTEXT" {
					indexes[key] = append(indexes[key], in)
				}
			}
		} else {
			indexes[key] = t.Indexes
		}
	}
	return dbh, tables, indexes, nil
}

func SetTableFields(arr [][]string, version string) []config.MySQLField {
	var res []config.MySQLField
	useCURRENTTIMESTAMP := false

	if version < MySQLVersionLimit {
		// availability CURRENT_TIMESTAMP
		for _, f := range arr {
			extra := strings.Trim(f[4], " \n\r")
			if strings.ToUpper(extra) == "ON UPDATE CURRENT_TIMESTAMP" {
				useCURRENTTIMESTAMP = true
				break
			}
		}
	}
	for _, f := range arr {
		n := false
		if f[2] == "true" {
			n = true
		}

		t := strings.Trim(f[1], " \n\r")
		def := strings.Trim(f[3], " \n\r")
		extra := strings.Trim(f[4], " \n\r")

		if strings.ToLower(def) == "now()" {
			def = "CURRENT_TIMESTAMP"
		}
		if version >= MySQLVersionLimit {
			if strings.ToLower(t) == "timestamp" {
				t = "datetime"
			}
			if strings.ToLower(def) == "0000-00-00 00:00:00" {
				def = "CURRENT_TIMESTAMP"
			}

		} else {
			if strings.ToLower(t) == "datetime" {
				t = "timestamp"
			}
			if strings.ToUpper(def) == "CURRENT_TIMESTAMP" {
				if useCURRENTTIMESTAMP {
					def = "0000-00-00 00:00:00"
				}
				useCURRENTTIMESTAMP = true
			}
		}

		res = append(res, config.MySQLField{f[0], t, n, def, extra})
	}
	return res
}

func SetupMySQLTables(mysqlConfig config.ServerAccess, tablesConfig map[string]config.MySQLTableComponent, tables []config.Table) error {
	dbh, userTablesJSON, userIndexesJSON, err := prepareTableComponents(mysqlConfig, tablesConfig)
	if err != nil {
		return err
	}
	defer dbh.Close()

	for _, t := range tables {
		err = DropTable(dbh, t.Name)
		if err != nil {
			return err
		}
		log.Println("table '" + t.Name + "' is dropped ")

		err = CreateTable(dbh, t.Name, userTablesJSON[t.Key], userIndexesJSON[t.Key])
		if err != nil {
			return err
		}
		log.Println("table '" + t.Name + "' is created")
	}
	return nil
}
