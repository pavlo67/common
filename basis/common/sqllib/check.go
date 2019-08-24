package sqllib

import (
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pkg/errors"
)

const FieldNotFound = "field not found"
const IndexNotFound = "index not found"
const FieldNotUsed = "field not used"
const IndexNotUsed = "index not used"

const BadField = "bad field struct"
const BadIndex = "bad index"

func CheckTables(sqlOp Operator, tablesConfig map[string]config.SQLTable) ([]common.Info, error) {

	tablesConfig = PrepareTables(tablesConfig)
	isErr := false
	info := []common.Info{}

	for _, table := range tablesConfig {

		ok, err := TableExists(sqlOp, table.Name)
		if err != nil {
			info = append(info, common.Info{"check if table exists": table, "status": err.Error()})
			isErr = true
			continue
		}

		if !ok {
			info = append(info, common.Info{"check if table exists": table, "status": "does not exist"})
			isErr = true
			continue
		}

		//fields, err := TableFields(dbh, table)
		//if err != nil {
		//	info = append(info, basis.Info{"check table fields": table, "status": err.Error()})
		//	isErr = true
		//	continue
		//}

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

		//for _, idx := range table.Indexes {
		//	err = AddTableIndex(sqlOp, table.Title, idx.Title, idx.Type, idx.Fields)
		//
		//	if err != nil {
		//		return err
		//	}
		//
		//	log.Println("index '" + idx.Title + "' is created")
		//
		//}
	}

	if isErr {
		return info, errors.New("check isn't ok")
	}

	return info, nil
}

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
//		if err = rows.Scan(&r.Title, &r.Type, &null, &key, &rDefault, &r.Extra); err != nil {
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
//		if len(indexes) > 0 && indexes[len(indexes)-1].Title == name {
//			indexes[len(indexes)-1].Fields = append(indexes[len(indexes)-1].Fields, column)
//			continue
//		}
//
//		ii := sqllib.SQLIndex{Title: name, Fields: []string{column}, IndexType: string(iType)}
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
//			if f.Title == f1.Title {
//				found = true
//				if strings.ToLower(f.Type) != strings.ToLower(f1.Type) {
//					res[BadField] += f.Title + ".Genus is: " + f1.Type + ", must be: " + f.Type + "; "
//				}
//				if f.Null != f1.Null {
//					res[BadField] += f.Title + ".Null is: " + strconv.FormatBool(f1.Null) + ", must be: " + strconv.FormatBool(f.Null) + "; "
//				}
//				if f.Default != f1.Default {
//					res[BadField] += f.Title + ".Default is: " + f1.Default + ", must be: " + f.Default + "; "
//				}
//				if f.Extra != f1.Extra {
//					// old mysql version does not see 'on update CURRENT_TIMESTAMP'
//					if version < MySQLVersionLimit && f.Extra == "on update CURRENT_TIMESTAMP" {
//						log.Printf("Ignore field: '%s' check  error: 'on update CURRENT_TIMESTAMP' for mysql version: %s", f.Title, version)
//					} else {
//						res[BadField] += f.Title + ".Extra is: " + f1.Extra + ", must be: " + f.Extra + "; "
//					}
//				}
//				break
//			}
//		}
//		if !found {
//			res[FieldNotFound] += f.Title + ","
//		}
//	}
//	//find not used fields
//	for _, f1 := range is {
//		found := false
//		for _, f := range need {
//			if f.Title == f1.Title {
//				found = true
//				break
//			}
//		}
//		if !found {
//			res[FieldNotUsed] += f1.Title + ","
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
//			if i.Title == i1.Title {
//				found = true
//				if i.Type != i1.Type && i.Type != i1.IndexType {
//					res[BadIndex] += i.Title + ".Genus is: " + i1.Type + ", must be: " + i.Type + "; "
//				}
//				if len(i.Fields) != len(i1.Fields) {
//					res[BadIndex] += i.Title + ".FieldsArr is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
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
//						res[BadIndex] += i.Title + ".FieldsArr is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
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
//				log.Println(IndexNotFound+": table=", table, "; index=", i.Title)
//				if err := AddTableIndex(dbh, table, i.Title, i.Type, i.Fields); err != nil {
//					log.Println(err.Error(), "can't add new index:", i.Title, "for table:", table)
//				}
//			}
//		}
//	}
//	//find not used indexes
//	for _, i1 := range is {
//		found := false
//		for _, i := range need {
//			if i.Title == i1.Title {
//				found = true
//				break
//			}
//		}
//		if !found {
//			log.Println(IndexNotUsed+": table=", table, "; index=", i1.Title)
//			if err := DropTableIndex(dbh, table, i1.Title); err != nil {
//				log.Println(err.Error(), "can't drop unused index:", i1.Title, "for table:", table)
//			}
//			//res[IndexNotUsed] += i1.Nick +": type="+ i1.Genus + "; fields=[" + strings.Join123(i1.FieldsArr, ", ") + "]; indexType=" + i1.IndexType  + ";   "
//		}
//	}
//	return res
//}
