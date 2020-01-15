package sqllib_pg

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"strconv"

	"regexp"

	"github.com/pavlo67/workshop/common/config"
)

func AddressPostgres(e config.Access) (string, error) {
	if e.Host == "" {
		e.Host = "localhost"
	}

	if e.Port == 0 {
		e.Port = 5432
	}

	return fmt.Sprintf(
		//"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"user=%s password=%s host=%s port=%d dbname=%s %s",
		e.User,
		e.Pass,
		e.Host,
		e.Port,
		e.Path,
		e.Options,
	), nil

}

var reQuestionMark = regexp.MustCompile("\\?")

func CorrectWildcards(query string) string {
	n := 1

LOOP:
	loc := reQuestionMark.FindStringIndex(query)
	if len(loc) > 0 {
		query = query[:loc[0]] + "$" + strconv.Itoa(n) + query[loc[1]:]
		n++

		goto LOOP
	}

	return query
}

func WildcardsForUpdate(fields []string) string {
	var wildcards []string

	for n, f := range fields {
		wildcards = append(wildcards, f+"=$"+strconv.Itoa(n+1))
	}

	return strings.Join(wildcards, ",")
}

func WildcardsForInsert(fields []string) string {
	var wildcards []string
	for n := 1; n <= len(fields); n++ {
		wildcards = append(wildcards, "$"+strconv.Itoa(n))
	}

	return strings.Join(wildcards, ",")
}

func Connect(access config.Access) (*sql.DB, error) {
	if strings.TrimSpace(access.Path) == "" {
		return nil, errors.New("no path to Postgres database is defined")
	}

	address, err := AddressPostgres(access)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, errors.Wrapf(err, "wrong db connect (access = %#v)", access)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "wrong .Ping on db connect (access = %#v)", access)
	}

	return db, nil
}
