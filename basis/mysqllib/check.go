package mysqllib

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
)

func CheckMySQLTables(mysqlConfig config.ServerAccess, tablesConfig map[string]config.MySQLTableComponent, tables []config.Table) ([]starter.Info, error) {
	dbh, userTablesJSON, userIndexesJSON, err := prepareTableComponents(mysqlConfig, tablesConfig)
	if err != nil {
		return nil, err
	}
	defer dbh.Close()

	var userTables = map[string][]config.MySQLField{}
	var userIndexes = map[string][]config.MySQLIndex{}
	var ok bool
	var info []starter.Info
	isErr := false
	for _, t := range tables {
		if userTables[t.Name], ok = userTablesJSON[t.Key]; !ok {
			info = append(info, starter.Info{Path: t.Key, Status: "can't find table structure in mysql.json5"})
			isErr = true
		}
		if userIndexes[t.Name], ok = userIndexesJSON[t.Key]; !ok {
			info = append(info, starter.Info{Path: t.Key, Status: "can't find table indexes in mysql.json5"})
		}
	}

	if !isErr {
		dbh, err := ConnectToMysql(mysqlConfig)
		if err != nil {
			return nil, errors.Wrap(err, "error connect to mySQL")
		}
		defer dbh.Close()

		info, isErr = CheckTables(dbh, userTables, userIndexes)
	}

	if isErr {
		return info, errors.New("check isn't ok")
	}

	return info, nil
}

func CheckTables(dbh *sql.DB, tablesFields map[string][]config.MySQLField, tablesIndexes map[string][]config.MySQLIndex) ([]starter.Info, bool) {
	isErr := false
	info := []starter.Info{}
	for table := range tablesFields {
		err := TableExists(dbh, table)
		if err != nil {
			info = append(info, starter.Info{Path: table, Status: err.Error()})
			isErr = true
			continue
		}
		fields, err := TableFields(dbh, table)
		if err != nil {
			info = append(info, starter.Info{Path: table, Status: err.Error()})
			isErr = true
			continue
		}

		indexes, err := TableIndexes(dbh, table)
		if err != nil {
			info = append(info, starter.Info{Path: table, Status: err.Error()})
			isErr = true
			continue
		}

		resF, err := CheckTableFields(dbh, fields, tablesFields[table])
		if err != nil {
			info = append(info, starter.Info{Path: table, Status: err.Error()})
			isErr = true
		}
		if len(resF) > 0 {
			for e, v := range resF {
				i := starter.Info{Path: table, Status: e, Details: v}
				info = append(info, i)
			}
			isErr = true
			//continue
		}
		if _, ok := tablesIndexes[table]; ok {
			resI := CheckTableIndexes(dbh, table, indexes, tablesIndexes[table])
			if len(resI) > 0 {
				for e, v := range resI {
					i := starter.Info{Path: table, Status: e, Details: v}
					info = append(info, i)
				}
				isErr = true
				continue
			}
		}
	}
	return info, isErr
}

func CheckTableFields(dbh *sql.DB, is, need []config.MySQLField) (map[string]string, error) {
	version, err := MySQLVersion(dbh)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	//check used fields
	for _, f := range need {
		found := false
		for _, f1 := range is {
			if f.Name == f1.Name {
				found = true
				if strings.ToLower(f.Type) != strings.ToLower(f1.Type) {
					res[BadField] += f.Name + ".Genus is: " + f1.Type + ", must be: " + f.Type + "; "
				}
				if f.Null != f1.Null {
					res[BadField] += f.Name + ".Null is: " + strconv.FormatBool(f1.Null) + ", must be: " + strconv.FormatBool(f.Null) + "; "
				}
				if f.Default != f1.Default {
					res[BadField] += f.Name + ".Default is: " + f1.Default + ", must be: " + f.Default + "; "
				}
				if f.Extra != f1.Extra {
					// old mysql version does not see 'on update CURRENT_TIMESTAMP'
					if version < MySQLVersionLimit && f.Extra == "on update CURRENT_TIMESTAMP" {
						log.Printf("Ignore field: '%s' check  error: 'on update CURRENT_TIMESTAMP' for mysql version: %s", f.Name, version)
					} else {
						res[BadField] += f.Name + ".Extra is: " + f1.Extra + ", must be: " + f.Extra + "; "
					}
				}
				break
			}
		}
		if !found {
			res[FieldNotFound] += f.Name + ","
		}
	}
	//find not used fields
	for _, f1 := range is {
		found := false
		for _, f := range need {
			if f.Name == f1.Name {
				found = true
				break
			}
		}
		if !found {
			res[FieldNotUsed] += f1.Name + ","
		}
	}
	return res, nil
}

func CheckTableIndexes(dbh *sql.DB, table string, is, need []config.MySQLIndex) map[string]string {
	version, err := MySQLVersion(dbh)
	if err != nil {
		log.Println(err.Error(), "can't get version of mySQL")
	}
	res := map[string]string{}
	//check used indexes
	for _, i := range need {
		found := false
		for _, i1 := range is {
			if i.Name == i1.Name {
				found = true
				if i.Type != i1.Type && i.Type != i1.IndexType {
					res[BadIndex] += i.Name + ".Genus is: " + i1.Type + ", must be: " + i.Type + "; "
				}
				if len(i.Fields) != len(i1.Fields) {
					res[BadIndex] += i.Name + ".Fields is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
					break
				}
				for _, f := range i.Fields {
					found1 := false
					for _, f1 := range i1.Fields {
						if f == f1 {
							found1 = true
							break
						}
					}
					if !found1 {
						res[BadIndex] += i.Name + ".Fields is: " + strings.Join(i1.Fields, ",") + ", must be: " + strings.Join(i.Fields, ",") + "; "
						break
					}
				}
				break
			}
		}
		if !found {
			if version < MySQLVersionFullTextLimit && strings.ToUpper(i.Type) == "FULLTEXT" {
				log.Println("ignore FULLTEXT index for table:", table, "; current mySQL:", version)
			} else {
				log.Println(IndexNotFound+": table=", table, "; index=", i.Name)
				if err := AddTableIndex(dbh, table, i.Name, i.Type, i.Fields); err != nil {
					log.Println(err.Error(), "can't add new index:", i.Name, "for table:", table)
				}
			}
		}
	}
	//find not used indexes
	for _, i1 := range is {
		found := false
		for _, i := range need {
			if i.Name == i1.Name {
				found = true
				break
			}
		}
		if !found {
			log.Println(IndexNotUsed+": table=", table, "; index=", i1.Name)
			if err := DropTableIndex(dbh, table, i1.Name); err != nil {
				log.Println(err.Error(), "can't drop unused index:", i1.Name, "for table:", table)
			}
			//res[IndexNotUsed] += i1.Nick +": type="+ i1.Genus + "; fields=[" + strings.Join123(i1.Fields, ", ") + "]; indexType=" + i1.IndexType  + ";   "
		}
	}
	return res
}
