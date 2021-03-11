package persons_sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/strlib"
)

var fieldsToInsert = []string{"issued_id", "nickname", "email", "roles", "creds", "data", "history"}
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

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove, stmClean *sql.Stmt
}

const onNew = "on personsSQLite.New(): "

func New(db *sql.DB, table string) (persons.Operator, crud.Cleaner, error) {
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

		sqlClean: "DELETE FROM " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&personsOp.stmInsert, personsOp.sqlInsert},
		{&personsOp.stmUpdate, personsOp.sqlUpdate},
		{&personsOp.stmRead, personsOp.sqlRead},
		{&personsOp.stmRemove, personsOp.sqlRemove},
		{&personsOp.stmClean, personsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &personsOp, &personsOp, nil
}

const onAdd = "on personsSQLite.Add(): "

func (personsOp *personsSQLite) Add(identity auth.Identity, creds auth.Creds, data common.Map, options *crud.Options) (auth.ID, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onAdd, "identity": identity, "data": data, "requestedRole": rbac.RoleAdmin})
	}

	var err error

	rolesBytes := []byte{} // to satisfy NOT NULL constraint
	if len(identity.Roles) > 0 {
		if rolesBytes, err = json.Marshal(identity.Roles); err != nil {
			return "", errors.Wrapf(err, onAdd+"can't marshal identity.Roles (%#v)", identity.Roles)
		}
	}

	credsBytes := []byte{} // to satisfy NOT NULL constraint
	if len(creds) > 0 {
		if credsBytes, err = json.Marshal(creds); err != nil {
			return "", errors.Wrapf(err, onAdd+"can't marshal creds (%#v)", creds)
		}
	}

	dataBytes := []byte{} // to to satisfy NOT NULL constraint
	if len(data) > 0 {
		if dataBytes, err = json.Marshal(data); err != nil {
			return "", errors.Wrapf(err, onAdd+"can't marshal data (%#v)", data)
		}
	}

	historyBytes := []byte{} // to to satisfy NOT NULL constraint
	// TODO!!! append to .History

	// "issued_id", "nickname", "email", "roles", "creds", "data", "history"
	values := []interface{}{identity.IssuedID, identity.Nickname, creds[auth.CredsEmail], rolesBytes, credsBytes, dataBytes, historyBytes}

	res, err := personsOp.stmInsert.Exec(values...)
	if err != nil {
		return "", errors.Wrapf(err, onAdd+sqllib.CantExec, personsOp.sqlInsert, strlib.Stringify(values))
	}

	idSQLite, err := res.LastInsertId()
	if err != nil {
		return "", errors.Wrapf(err, onAdd+sqllib.CantGetLastInsertId, personsOp.sqlInsert, strlib.Stringify(values))
	}

	return auth.ID(strconv.FormatInt(idSQLite, 10)), nil
}

const onChange = "on personsSQLite.Change(): "

func (personsOp *personsSQLite) Change(item persons.Item, options *crud.Options) (*persons.Item, error) {
	if options == nil || options.Identity == nil {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	itemOld, err := personsOp.read(item.Identity.ID)
	if err != nil || itemOld == nil {
		errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
		if options.HasRole(rbac.RoleAdmin) {
			return nil, errors.CommonError(common.WrongIDKey, common.Map{"on": onChange, "item": item, "reason": errorStr})
		} else {
			l.Error(errorStr)
			return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onChange, "item": item, "requestedRole": rbac.RoleAdmin})
		}
	}

	if itemOld.Identity.ID != options.Identity.ID && !options.Identity.Roles.Has(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onChange, "item": item})
	}

	rolesBytes := []byte{} // to satisfy NOT NULL constraint
	if len(item.Identity.Roles) > 0 {
		if rolesBytes, err = json.Marshal(item.Identity.Roles); err != nil {
			return nil, errors.Wrapf(err, onChange+"can't marshal item.Identity.Roles (%#v)", item.Identity.Roles)
		}
	}

	creds := item.Creds()

	email := creds[auth.CredsEmail]
	delete(creds, auth.CredsEmail)

	credsBytes := []byte{} // to satisfy NOT NULL constraint
	if len(creds) > 0 {
		if credsBytes, err = json.Marshal(creds); err != nil {
			return nil, errors.Wrapf(err, onChange+"can't marshal creds (%#v)", creds)
		}
	}

	dataBytes := []byte{} // to to satisfy NOT NULL constraint
	if len(item.Data) > 0 {
		if dataBytes, err = json.Marshal(item.Data); err != nil {
			return nil, errors.Wrapf(err, onChange+"can't marshal data (%#v)", item.Data)
		}
	}

	// TODO!!! append to item.History

	historyBytes := []byte{} // to to satisfy NOT NULL constraint
	if len(item.History) > 0 {
		historyBytes, err = json.Marshal(item.History)
		if err != nil {
			return nil, errors.Wrapf(err, onChange+"can't marshal .History(%#v)", item.History)
		}
	}

	// "issued_id", "nickname", "email", "roles", "creds", "data", "history", "updated_at"
	values := []interface{}{item.Identity.IssuedID, item.Identity.Nickname, email, rolesBytes,
		credsBytes, dataBytes, historyBytes, time.Now(), item.ID}

	if _, err = personsOp.stmUpdate.Exec(values...); err != nil {
		return nil, errors.Wrapf(err, onChange+sqllib.CantExec, personsOp.sqlUpdate, strlib.Stringify(values))
	}

	// TODO??? re-read
	return &item, nil
}

const onRemove = "on personsSQLite.Remove()"

func (personsOp *personsSQLite) Remove(id auth.ID, options *crud.Options) error {
	if id != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
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
		&item.Identity.IssuedID, &item.Identity.Nickname, &email, &rolesBytes, &credsBytes, &dataBytes,
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

func (personsOp *personsSQLite) Read(id auth.ID, options *crud.Options) (*persons.Item, error) {
	if id != options.Identity.ID && !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onRead, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	return personsOp.read(id)
}

const onList = "on personsSQLite.List()"

func (personsOp *personsSQLite) List(options *crud.Options) ([]persons.Item, error) {
	if !options.HasRole(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
	}

	var condition string
	var values []interface{}

	if selector := options.GetSelector(); selector != nil {
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

	query := sqllib.SQLList(personsOp.table, fieldsToListStr, condition, options)
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
			&item.Identity.IssuedID, &item.Identity.Nickname, &email, &rolesBytes, &credsBytes, &dataBytes,
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

func (personsOp *personsSQLite) Close() error {
	return errors.Wrap(personsOp.db.Close(), "on personsSQLite.Close()")
}
