package sqllib

const FieldNotFound = "field not found"
const IndexNotFound = "index not found"
const FieldNotUsed = "field not used"
const IndexNotUsed = "index not used"

const BadField = "bad field struct"
const BadIndex = "bad index"

//func TableFields(dbh *sql.DB, table string) ([]sqllib.SQLField, error) {
//	var stmt *sql.Stmt
//	sqlQuery := "desc `" + table + "`"
//	if err := Exec(dbh, sqlQuery, &stmt); err != nil {
//		return nil, err
//	}
//	rows, err := stmt.Query()
//	if err != nil {
//		return nil, errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
//	}
//	if rows != nil {
//		defer rows.Close()
//	}
//
//	var null, key string
//	var rDefault []byte
//	var fields []sqllib.SQLField
//	for rows.Next() {
//		r := sqllib.SQLField{}
//		if err = rows.Scan(&r.Name, &r.Type, &null, &key, &rDefault, &r.Extra); err != nil {
//			return nil, errors.Wrapf(err, "can't scan query (sql='%v')", sqlQuery)
//		}
//		r.Default = string(rDefault)
//		if strings.ToUpper(null) == "NO" {
//			r.Null = false
//		} else {
//			r.Null = true
//		}
//		fields = append(fields, r)
//	}
//	return fields, nil
//}
//
//func TableIndexes(dbh *sql.DB, table string) ([]sqllib.SQLIndex, error) {
//	version, err := MySQLVersion(dbh)
//	if err != nil {
//		return nil, err
//	}
//	var stmt *sql.Stmt
//	sqlQuery := "show index from `" + table + "`"
//	if err := Exec(dbh, sqlQuery, &stmt); err != nil {
//		return nil, err
//	}
//	rows, err := stmt.Query()
//	if err != nil {
//		return nil, errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
//	}
//	if rows != nil {
//		defer rows.Close()
//	}
//
//	var t, notUnique, name, null, sec, column string
//	var coll, card, sub, pack, iType, com, iCom, visible []byte
//	var indexes []sqllib.SQLIndex
//
//	for rows.Next() {
//		if version < MySQLVersionLimit {
//			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com)
//		} else if version < MySQLVersionLimit2 {
//			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com, &iCom)
//		} else {
//			err = rows.Scan(&t, &notUnique, &name, &sec, &column, &coll, &card, &sub, &pack, &null, &iType, &com, &iCom, &visible)
//		}
//		if err != nil {
//			return nil, errors.Wrapf(err, "can't scan query (sql='%v')", sqlQuery)
//		}
//		if sub != nil {
//			column += "(" + string(sub) + ")"
//		}
//		if len(indexes) > 0 && indexes[len(indexes)-1].Name == name {
//			indexes[len(indexes)-1].Fields = append(indexes[len(indexes)-1].Fields, column)
//			continue
//		}
//
//		ii := sqllib.SQLIndex{Name: name, Fields: []string{column}, IndexType: string(iType)}
//
//		if strings.ToUpper(name) == "PRIMARY" {
//			ii.Type = "PRIMARY"
//		} else if notUnique == "0" {
//			ii.Type = "UNIQUE"
//		}
//
//		indexes = append(indexes, ii)
//	}
//
//	return indexes, nil
//}
//
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
