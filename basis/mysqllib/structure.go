package mysqllib

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/starter/config"
)

const FieldNotFound = "field not found"
const IndexNotFound = "index not found"
const FieldNotUsed = "field not used"
const IndexNotUsed = "index not used"

const BadField = "bad field struct"
const BadIndex = "bad index"

//func SetupMySQLTables(mysqlConfig config.ServerAccess, tablesConfig map[string]config.MySQLTableComponent, tables []config.Table) error {
//	dbh, userTablesJSON, userIndexesJSON, err := prepareTableComponents(mysqlConfig, tablesConfig)
//	if err != nil {
//		return err
//	}
//	defer dbh.Close()
//
//	for _, t := range tables {
//		err = DropTable(dbh, t.Name)
//		if err != nil {
//			return err
//		}
//		log.Println("table '" + t.Name + "' is dropped ")
//
//		err = CreateTable(dbh, t.Name, userTablesJSON[t.Key], userIndexesJSON[t.Key])
//		if err != nil {
//			return err
//		}
//		log.Println("table '" + t.Name + "' is created")
//	}
//	return nil
//}

func TableFields(dbh *sql.DB, table string) ([]config.MySQLField, error) {
	var stmt *sql.Stmt
	sqlQuery := "desc `" + table + "`"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
	}
	if rows != nil {
		defer rows.Close()
	}

	var null, key string
	var rDefault []byte
	var fields []config.MySQLField
	for rows.Next() {
		r := config.MySQLField{}
		if err = rows.Scan(&r.Name, &r.Type, &null, &key, &rDefault, &r.Extra); err != nil {
			return nil, errors.Wrapf(err, "can't scan query (sql='%v')", sqlQuery)
		}
		r.Default = string(rDefault)
		if strings.ToUpper(null) == "NO" {
			r.Null = false
		} else {
			r.Null = true
		}
		fields = append(fields, r)
	}
	return fields, nil
}

func TableIndexes(dbh *sql.DB, table string) ([]config.MySQLIndex, error) {
	version, err := MySQLVersion(dbh)
	if err != nil {
		return nil, err
	}
	var stmt *sql.Stmt
	sqlQuery := "show index from `" + table + "`"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
	}
	if rows != nil {
		defer rows.Close()
	}

	var t, notUnique, name, null, sec, column string
	var coll, card, sub, pack, iType, com, iCom, visible []byte
	var indexes []config.MySQLIndex

	for rows.Next() {
		if version < MySQLVersionLimit {
			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com)
		} else if version < MySQLVersionLimit2 {
			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com, &iCom)
		} else {
			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com, &iCom, &visible)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "can't scan query (sql='%v')", sqlQuery)
		}
		if sub != nil {
			column += "(" + string(sub) + ")"
		}
		if len(indexes) > 0 && indexes[len(indexes)-1].Name == name {
			indexes[len(indexes)-1].Fields = append(indexes[len(indexes)-1].Fields, column)
			continue
		}

		ii := config.MySQLIndex{Name: name, Fields: []string{column}, IndexType: string(iType)}

		if strings.ToUpper(name) == "PRIMARY" {
			ii.Type = "PRIMARY"
		} else if notUnique == "0" {
			ii.Type = "UNIQUE"
		}

		indexes = append(indexes, ii)
	}

	return indexes, nil
}

func DropTableIndex(dbh *sql.DB, table, indexName string) error {

	var stmt *sql.Stmt
	var err error
	if strings.ToUpper(indexName) == "PRIMARY" {
		return errors.New("reindex does not change PRIMARY index ")
	}
	sqlQuery := "alter table `" + table + "` drop index `" + indexName + "`"
	if stmt, err = dbh.Prepare(sqlQuery); err != nil {
		return errors.Wrapf(err, "can't prepare sql: %v", sqlQuery)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		return errors.Wrapf(err, "can't exec sql: %v, values=%v", sqlQuery, indexName)
	}
	log.Println(table, ": drop index `"+indexName+"`")
	return nil
}

func AddTableIndex(dbh *sql.DB, table, indexName, indexType string, indexFields []string) error {

	var stmt *sql.Stmt
	sqlQuery := "alter table `" + table + "` add "
	if strings.ToUpper(indexType) == "PRIMARY" {
		return errors.New("reindex does not change PRIMARY index ")
	}
	if strings.ToUpper(indexType) == "UNIQUE" || indexType == "" || strings.ToUpper(indexType) == "FULLTEXT" {
		if indexType == strings.ToUpper("UNIQUE") {
			sqlQuery += "UNIQUE KEY `" + indexName + "` ("
		} else if strings.ToUpper(indexType) == "FULLTEXT" {
			sqlQuery += "FULLTEXT KEY `" + indexName + "` ("
		} else {
			sqlQuery += "INDEX `" + indexName + "` USING BTREE ("
		}
		liF := 0
		for _, f := range indexFields {
			if liF > 0 {
				sqlQuery += ", "
			}
			sqlQuery += "`" + f + "`"
			liF++
		}
		sqlQuery += ")"
	}
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return err
	}
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "can't exec sql:%v", sqlQuery)
	}
	log.Println(table, ": add index `"+indexName+"`: ", sqlQuery)
	return nil
}

var reIndexLimit = regexp.MustCompile(`\(\d+\)`)

func CreateTable(dbh *sql.DB, table string, tableFields []config.MySQLField, tableIndexes []config.MySQLIndex) error {
	var stmt *sql.Stmt
	reText := regexp.MustCompile("text")
	reTime := regexp.MustCompile("(datetime|timestamp)")
	sqlQuery := "create table `" + table + "` ( "
	firsFieldAdded := false
	for _, f := range tableFields {
		if firsFieldAdded {
			sqlQuery += ", \n"
		}
		sqlQuery += "`" + f.Name + "` " + f.Type + " "
		if !f.Null {
			sqlQuery += " NOT NULL "
		} else {
			sqlQuery += " NULL "
		}
		if strings.ToUpper(f.Default) != "CURRENT_TIMESTAMP" && strings.ToUpper(f.Default) != "NOW()" {
			if f.Default != "" {
				sqlQuery += " default '" + f.Default + "' "
			} else if !reText.MatchString(f.Type) && !reTime.MatchString(f.Type) && f.Extra != "auto_increment" && !f.Null {
				sqlQuery += " default '' "
			}
		} else {
			sqlQuery += " default CURRENT_TIMESTAMP "
		}
		if f.Extra != "" {
			sqlQuery += f.Extra + " "
		}
		firsFieldAdded = true
	}
	for _, i := range tableIndexes {
		sqlQuery += ", \n"
		if strings.ToUpper(i.Type) == "PRIMARY" {
			//sqlQuery += "PRIMARY KEY  (`" + i.Fields[0] + "`)"
			sqlQuery += "PRIMARY KEY  ("
			liF := 0
			for _, f := range i.Fields {
				if liF > 0 {
					sqlQuery += ", "
				}
				sqlQuery += "`" + f + "`"
				liF++
			}
			sqlQuery += ")"
		}
		if strings.ToUpper(i.Type) == "UNIQUE" || i.Type == "" || strings.ToUpper(i.Type) == "FULLTEXT" {
			if i.Type == strings.ToUpper("UNIQUE") {
				sqlQuery += "UNIQUE KEY `" + i.Name + "` ("
			} else if strings.ToUpper(i.Type) == "FULLTEXT" {
				sqlQuery += "FULLTEXT KEY `" + i.Name + "` ("
			} else {
				sqlQuery += "KEY `" + i.Name + "` USING BTREE ("
			}
			liF := 0
			for _, f := range i.Fields {
				if liF > 0 {
					sqlQuery += ", "
				}
				if reIndexLimit.MatchString(f) {
					sqlQuery += f
				} else {
					sqlQuery += "`" + f + "`"
				}
				liF++
			}
			sqlQuery += ")"
		}
	}
	sqlQuery += ") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return err
	}
	_, err := stmt.Exec()
	if err != nil {
		return errors.Wrapf(err, "can't exec SQL: %s", sqlQuery)
	}

	return nil
}
