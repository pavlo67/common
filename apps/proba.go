package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/libraries/sqllib"
	"github.com/pavlo67/common/common/libraries/sqllib/sqllib_pg"
)

func main() {
	db, err := sqllib_pg.Connect(config.Access{
		Host: "localhost",
		User: "msq",
		Pass: "msq_psw1",
		Path: "nb_prod",
	})
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT created_at FROM " + "storage WHERE created_at < $1;"

	stm, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("can't db.Prepare(%s): %s", query, err)
	}

	before := time.Date(2020, 01, 30, 0, 0, 0, 0, time.UTC)
	values := []interface{}{before}

	rows, err := stm.Query(values...)

	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Fatalf(sqllib.CantQuery+": %s", query, values, err)
	}
	defer rows.Close()

	for rows.Next() {
		var createdAt time.Time

		err := rows.Scan(&createdAt)

		if err != nil {
			log.Fatalf(sqllib.CantScanQueryRow+": %s", query, values, err)
		}

		log.Print(createdAt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf(sqllib.RowsError+": %s", query, values, err)
	}
}
