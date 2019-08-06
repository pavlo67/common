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

//
//func CheckSQLTables(cfg config.ServerAccess, tablesConfig map[string]SQLTable) ([]basis.Info, error) {
//	dbh, userTablesJSON, userIndexesJSON, err := PrepareTables(cfg, tablesConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	var userTables = map[string][]sqllib.SQLField{}
//	var userIndexes = map[string][]sqllib.SQLIndex{}
//	var ok bool
//	var info []basis.Info
//	isErr := false
//	for _, t := range tables {
//		if userTables[t.Name], ok = userTablesJSON[t.Key]; !ok {
//			info = append(info, basis.Info{"key": t.Key, "status": "can't find table structure in mysql.json5"})
//			isErr = true
//		}
//		if userIndexes[t.Name], ok = userIndexesJSON[t.Key]; !ok {
//			info = append(info, basis.Info{"key": t.Key, "status": "can't find table indexes in mysql.json5"})
//		}
//	}
//
//	if !isErr {
//		dbh, err := ConnectToMysql(cfg)
//		if err != nil {
//			return nil, errors.Wrap(err, "error connect to mySQL")
//		}
//		defer dbh.Close()
//
//		info, isErr = CheckTables(dbh, userTables, userIndexes)
//	}
//
//	if isErr {
//		return info, errors.New("check isn't ok")
//	}
//
//	return info, nil
//}
//
//func CheckTables(dbh *sql.DB, tablesFields map[string][]sqllib.SQLField, tablesIndexes map[string][]sqllib.SQLIndex) ([]basis.Info, bool) {
//	isErr := false
//	info := []basis.Info{}
//	for table := range tablesFields {
//		err := TableExists(dbh, table)
//		if err != nil {
//			info = append(info, basis.Info{"key": table, "status": err.Error()})
//			isErr = true
//			continue
//		}
//		fields, err := TableFields(dbh, table)
//		if err != nil {
//			info = append(info, basis.Info{"key": table, "status": err.Error()})
//			isErr = true
//			continue
//		}
//
//		indexes, err := TableIndexes(dbh, table)
//		if err != nil {
//			info = append(info, basis.Info{"key": table, "status": err.Error()})
//			isErr = true
//			continue
//		}
//
//		resF, err := CheckTableFields(dbh, fields, tablesFields[table])
//		if err != nil {
//			info = append(info, basis.Info{"key": table, "status": err.Error()})
//			isErr = true
//		}
//		if len(resF) > 0 {
//			for e, v := range resF {
//				i := basis.Info{"key": table, "status": e, "details": v}
//				info = append(info, i)
//			}
//			isErr = true
//			//continue
//		}
//		if _, ok := tablesIndexes[table]; ok {
//			resI := CheckTableIndexes(dbh, table, indexes, tablesIndexes[table])
//			if len(resI) > 0 {
//				for e, v := range resI {
//					i := basis.Info{"key": table, "status": e, "details": v}
//					info = append(info, i)
//				}
//				isErr = true
//				continue
//			}
//		}
//	}
//	return info, isErr
//}
//
//func CheckTableFields(dbh *sql.DB, is, need []sqllib.SQLField) (map[string]string, error) {
//	version, err := MySQLVersion(dbh)
//	if err != nil {
//		return nil, err
//	}
//	if err != nil {
//		return nil, err
//	}
//	res := map[string]string{}
//	//check used fields
//	for _, f := range need {
//		found := false
//		for _, f1 := range is {
//			if f.Name == f1.Name {
//				found = true
//				if strings.ToLower(f.Type) != strings.ToLower(f1.Type) {
//					res[BadField] += f.Name + ".Genus is: " + f1.Type + ", must be: " + f.Type + "; "
//				}
//				if f.Null != f1.Null {
//					res[BadField] += f.Name + ".Null is: " + strconv.FormatBool(f1.Null) + ", must be: " + strconv.FormatBool(f.Null) + "; "
//				}
//				if f.Default != f1.Default {
//					res[BadField] += f.Name + ".Default is: " + f1.Default + ", must be: " + f.Default + "; "
//				}
//				if f.Extra != f1.Extra {
//					// old mysql version does not see 'on update CURRENT_TIMESTAMP'
//					if version < MySQLVersionLimit && f.Extra == "on update CURRENT_TIMESTAMP" {
//						log.Printf("Ignore field: '%s' check  error: 'on update CURRENT_TIMESTAMP' for mysql version: %s", f.Name, version)
//					} else {
//						res[BadField] += f.Name + ".Extra is: " + f1.Extra + ", must be: " + f.Extra + "; "
//					}
//				}
//				break
//			}
//		}
//		if !found {
//			res[FieldNotFound] += f.Name + ","
//		}
//	}
//	//find not used fields
//	for _, f1 := range is {
//		found := false
//		for _, f := range need {
//			if f.Name == f1.Name {
//				found = true
//				break
//			}
//		}
//		if !found {
//			res[FieldNotUsed] += f1.Name + ","
//		}
//	}
//	return res, nil
//}
//
//func CheckTableIndexes(dbh *sql.DB, table string, is, need []sqllib.SQLIndex) map[string]string {
//	version, err := MySQLVersion(dbh)
//	if err != nil {
//		log.Println(err.Error(), "can't get version of mySQL")
//	}
//	res := map[string]string{}
//	//check used indexes
//	for _, i := range need {
//		found := false
//		for _, i1 := range is {
//			if i.Name == i1.Name {
//				found = true
//				if i.Type != i1.Type && i.Type != i1.IndexType {
//					res[BadIndex] += i.Name + ".Genus is: " + i1.Type + ", must be: " + i.Type + "; "
//				}
//				if len(i.Fields) != len(i1.Fields) {
//					res[BadIndex] += i.Name + ".FieldsArr is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
//					break
//				}
//				for _, f := range i.Fields {
//					found1 := false
//					for _, f1 := range i1.Fields {
//						if f == f1 {
//							found1 = true
//							break
//						}
//					}
//					if !found1 {
//						res[BadIndex] += i.Name + ".FieldsArr is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
//						break
//					}
//				}
//				break
//			}
//		}
//		if !found {
//			if version < MySQLVersionFullTextLimit && strings.ToUpper(i.Type) == "FULLTEXT" {
//				log.Println("ignore FULLTEXT index for table:", table, "; current mySQL:", version)
//			} else {
//				log.Println(IndexNotFound+": table=", table, "; index=", i.Name)
//				if err := AddTableIndex(dbh, table, i.Name, i.Type, i.Fields); err != nil {
//					log.Println(err.Error(), "can't add new index:", i.Name, "for table:", table)
//				}
//			}
//		}
//	}
//	//find not used indexes
//	for _, i1 := range is {
//		found := false
//		for _, i := range need {
//			if i.Name == i1.Name {
//				found = true
//				break
//			}
//		}
//		if !found {
//			log.Println(IndexNotUsed+": table=", table, "; index=", i1.Name)
//			if err := DropTableIndex(dbh, table, i1.Name); err != nil {
//				log.Println(err.Error(), "can't drop unused index:", i1.Name, "for table:", table)
//			}
//			//res[IndexNotUsed] += i1.Nick +": type="+ i1.Genus + "; fields=[" + strings.Join123(i1.FieldsArr, ", ") + "]; indexType=" + i1.IndexType  + ";   "
//		}
//	}
//	return res
//}
