package main

import (
	"fmt"
	"strconv"

	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pavlo67/constructor/components/basis/config"
	"github.com/pavlo67/constructor/components/basis/filelib"
	"github.com/pavlo67/constructor/components/basis/sqllib/sqllib_sqlite"
)

func main() {

	sqlOp := &sqllib_sqlite.SQLite{}

	err := sqlOp.Connect(config.ServerAccess{Path: filelib.CurrentPath() + "test.sqlite.db"})
	if err != nil {
		log.Fatal(err)
	}

	database, err := sqlOp.DB()
	if err != nil {
		log.Fatal(err)
	}

	//database, _ := sql.Open("sqlite3", "./nraboy.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Nic", "Raboy")
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}
