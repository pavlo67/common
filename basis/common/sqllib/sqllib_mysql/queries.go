package sqllib_mysql

import (
	"regexp"
	"strings"

	"github.com/pavlo67/workshop/basis/config"
)

var reIndexLimit = regexp.MustCompile(`\(\d+\)`)

func TableExistsSQL(tableName string) string {
	return "SHOW TABLES LIKE '" + tableName + "'"
}

func CreateSQL(table config.SQLTable) string {

	//if version < MySQLVersionFullTextLimit {
	//	// ignore FULLTEXT index for old mySQL version
	//	for _, in := range table.Indexes {
	//		if strings.ToUpper(in.Type) != "FULLTEXT" {
	//			indexes[key] = append(indexes[key], in)
	//		}
	//	}
	//} else {
	//	indexes[key] = table.Indexes
	//}

	// useCURRENTTIMESTAMP := false

	//if version < MySQLVersionLimit {
	//	// availability CURRENT_TIMESTAMP
	//	for _, f := range arr {
	//		extra := strings.Trim(f[4], " \n\r")
	//		if strings.ToUpper(extra) == "ON UPDATE CURRENT_TIMESTAMP" {
	//			useCURRENTTIMESTAMP = true
	//			break
	//		}
	//	}
	//}

	reText := regexp.MustCompile("text")
	reTime := regexp.MustCompile("(datetime|timestamp)")
	sqlQuery := "create table `" + table.Name + "` ( "
	firsFieldAdded := false

	//if strings.ToLower(def) == "now()" {
	//	def = "CURRENT_TIMESTAMP"
	//}
	//if version >= MySQLVersionLimit {
	//	if strings.ToLower(table) == "timestamp" {
	//		table = "datetime"
	//	}
	//	if strings.ToLower(def) == "0000-00-00 00:00:00" {
	//		def = "CURRENT_TIMESTAMP"
	//	}
	//
	//} else {
	//	if strings.ToLower(table) == "datetime" {
	//		table = "timestamp"
	//	}
	//	if strings.ToUpper(def) == "CURRENT_TIMESTAMP" {
	//		if useCURRENTTIMESTAMP {
	//			def = "0000-00-00 00:00:00"
	//		}
	//		useCURRENTTIMESTAMP = true
	//	}
	//}

	for _, f := range table.Fields {
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
	for _, i := range table.Indexes {
		sqlQuery += ", \n"
		if strings.ToUpper(i.Type) == "PRIMARY" {
			//sqlQuery += "PRIMARY KEY  (`" + i.FieldsArr[0] + "`)"
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

	return sqlQuery
}
