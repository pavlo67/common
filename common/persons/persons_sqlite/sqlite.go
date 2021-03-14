package persons_sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common/db"

	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/strlib"
)

var fieldsToInsert = []string{"urn", "nickname", "email", "roles", "creds", "data", "history"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToUpdate = append(fieldsToInsert, "updated_at")
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToRead = append(fieldsToUpdate, "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ persons.Operator = &personsSQLite{}

type personsSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlStat, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove, stmStat, stmClean *sql.Stmt
}

const onNew = "on personsSQLite.New(): "

func New(db *sql.DB, table string) (persons.Operator, db.Cleaner, error) {
	if table == "" {
		table = persons.CollectionDefault
	}

	personsOp := personsSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		sqlRemove: "DELETE FROM " + table + " where id = ?",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",
		sqlStat:   "SELECT COUNT(*) FROM " + table,

		sqlClean: "DELETE FROM " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&personsOp.stmInsert, personsOp.sqlInsert},
		{&personsOp.stmUpdate, personsOp.sqlUpdate},
		{&personsOp.stmRead, personsOp.sqlRead},
		{&personsOp.stmRemove, personsOp.sqlRemove},
		{&personsOp.stmStat, personsOp.sqlStat},
		{&personsOp.stmClean, personsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &personsOp, &personsOp, nil
}

const onSave = "on personsSQLite.Save(): "

func (personsOp *personsSQLite) Save(item persons.Item, identity *auth.Identity) (auth.ID, error) {
	if identity == nil || (item.ID != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item})
	}

	var err error
	rolesBytes := []byte{} // to satisfy NOT NULL constraint
	if len(item.Identity.Roles) > 0 {
		if rolesBytes, err = json.Marshal(item.Roles); err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal item.Identity.Roles (%#v)", item.Roles)
		}
	}

	creds := item.Creds()
	email := creds[auth.CredsEmail]
	delete(creds, auth.CredsEmail)

	credsBytes := []byte{} // to satisfy NOT NULL constraint
	if len(creds) > 0 {
		if credsBytes, err = json.Marshal(creds); err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal creds (%#v)", creds)
		}
	}

	dataBytes := []byte{} // to to satisfy NOT NULL constraint
	if len(item.Data) > 0 {
		if dataBytes, err = json.Marshal(item.Data); err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal data (%#v)", item.Data)
		}
	}

	// TODO!!! append to item.History

	historyBytes := []byte{} // to to satisfy NOT NULL constraint
	if len(item.History) > 0 {
		historyBytes, err = json.Marshal(item.History)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item.History)
		}
	}

	if item.ID != "" {
		itemOld, err := personsOp.read(item.Identity.ID)
		if err != nil || itemOld == nil {
			errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
			if identity.HasRole(rbac.RoleAdmin) {
				return "", errors.CommonError(common.WrongIDKey, common.Map{"on": onSave, "item": item, "reason": errorStr})
			} else {
				l.Error(errorStr)
				return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item, "requestedRole": rbac.RoleAdmin})
			}
		}
		// "issued_id", "nickname", "email", "roles", "creds", "data", "history", "updated_at"
		values := []interface{}{item.Identity.URN, item.Identity.Nickname, email, rolesBytes,
			credsBytes, dataBytes, historyBytes, time.Now(), item.ID}

		if _, err = personsOp.stmUpdate.Exec(values...); err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlUpdate, strlib.Stringify(values))
		}

	} else {
		// "issued_id", "nickname", "email", "roles", "creds", "data", "history"
		values := []interface{}{item.URN, item.Nickname, creds[auth.CredsEmail], rolesBytes, credsBytes, dataBytes, historyBytes}

		res, err := personsOp.stmInsert.Exec(values...)
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlInsert, strlib.Stringify(values))
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, personsOp.sqlInsert, strlib.Stringify(values))
		}

		item.ID = auth.ID(strconv.FormatInt(idSQLite, 10))
	}

	return item.ID, nil
}

const onRemove = "on personsSQLite.Remove()"

func (personsOp *personsSQLite) Remove(id auth.ID, identity *auth.Identity) error {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return errors.CommonError(common.NoRightsKey, common.Map{"on": onRemove, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return fmt.Errorf(onRemove+"wrong id (%s)", id)
	}

	if _, err = personsOp.stmRemove.Exec(idNum); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, personsOp.sqlRemove, idNum)
	}

	return nil
}

func (personsOp *personsSQLite) read(id auth.ID) (*persons.Item, error) {

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(onRead+"wrong id (%s)", id)
	}

	var item persons.Item
	var email string
	var rolesBytes, credsBytes, dataBytes, historyBytes []byte

	// "issued_id", "nickname", "email", "roles", "creds", "data", "history", "updated_at", "created_at"

	if err = personsOp.stmRead.QueryRow(idNum).Scan(
		&item.Identity.URN, &item.Identity.Nickname, &email, &rolesBytes, &credsBytes, &dataBytes,
		&historyBytes, &item.UpdatedAt, &item.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.CommonError(common.ErrNotFound, onRead)
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, personsOp.sqlRead, idNum)
	}

	if err := item.CompletePersonFromJSON(id, rolesBytes, credsBytes, dataBytes, historyBytes, email); err != nil {
		return nil, errors.CommonError(err, onRead)
	}

	return &item, nil
}

const onRead = "on personsSQLite.Read(): "

func (personsOp *personsSQLite) Read(id auth.ID, identity *auth.Identity) (*persons.Item, error) {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onRead, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	return personsOp.read(id)
}

const onList = "on personsSQLite.List()"

func (personsOp *personsSQLite) List(selector *selectors.Term, identity *auth.Identity) ([]persons.Item, error) {
	if !identity.HasRole(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
	}

	var condition string
	var values []interface{}

	if selector != nil {
		valuesStr, ok := selector.Values.([]string)
		if !ok {
			return nil, fmt.Errorf(onList+": wrong selector: %#v", selector)
		}

		switch selector.Key {
		case persons.HasEmail:
			if len(valuesStr) != 1 {
				return nil, fmt.Errorf(onList+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			condition = `email = ?`
			values = []interface{}{valuesStr[0]}

		case persons.HasNickname:
			if len(valuesStr) != 1 {
				return nil, fmt.Errorf(onList+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			condition = `nickname = ?`
			values = []interface{}{valuesStr[0]}

		default:
			return nil, fmt.Errorf(onList+": wrong selector.Key: %#v", selector)
		}
	}

	query := sqllib.SQLList(personsOp.table, fieldsToListStr, condition, identity)
	stm, err := personsOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	}

	//l.Infof("%s / %#v\n%s", condition, values, query)

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []persons.Item

	for rows.Next() {
		var idNum int64
		var item persons.Item

		var email string
		var rolesBytes, credsBytes, dataBytes, historyBytes []byte

		// "issued_id", "nickname", "email", "roles", "creds", "data", "history", "updated_at", "created_at"
		// "id"

		if err := rows.Scan(
			&item.Identity.URN, &item.Identity.Nickname, &email, &rolesBytes, &credsBytes, &dataBytes,
			&historyBytes, &item.UpdatedAt, &item.CreatedAt, &idNum); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, query, values)
		}

		if err := item.CompletePersonFromJSON(auth.ID(strconv.FormatInt(idNum, 10)), rolesBytes, credsBytes, dataBytes, historyBytes, email); err != nil {
			return nil, errors.CommonError(err, onList)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onStat = "on personsSQLite.Stat()"

func (personsOp *personsSQLite) Stat(*selectors.Term, *auth.Identity) (db.StatMap, error) {
	var num int
	if err := personsOp.stmStat.QueryRow().Scan(&num); err == sql.ErrNoRows {
		return nil, errors.CommonError(common.ErrNotFound, onStat)
	} else if err != nil {
		return nil, errors.Wrapf(err, onStat+sqllib.CantScanQueryRow, personsOp.sqlStat, nil)
	}

	return db.StatMap{"*": db.Stat{num}}, nil
}

func (personsOp *personsSQLite) Close() error {
	return errors.Wrap(personsOp.db.Close(), "on personsSQLite.Close()")
}
