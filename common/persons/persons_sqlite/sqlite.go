package persons_sqlite

import (
	"database/sql"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/common/common/sqllib"
)

var fieldsToInsert = []string{"nickname", "email", "roles", "creds", "data", "tags", "issued_id", "history"}
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

const onSave = "on personsSQLite.Save(): "

func (personsOp *personsSQLite) Add(identity auth.Identity, creds auth.Creds, data common.Map, options *crud.Options) (auth.ID, error) {
	return "", common.ErrNotImplemented

	//if options == nil || options.Identity == nil {
	//	return nil, errors.CommonError(common.NoRightsKey)
	//}

	//// TODO!!! rbac check
	//
	//if item.ID == "" {
	//	// TODO!!!
	//	item.OwnerID = options.Identity.ID
	//}
	//
	//var err error
	//
	//embeddedBytes := []byte{} // to satisfy NOT NULL constraint
	//if len(item.Content.Embedded) > 0 {
	//	if embeddedBytes, err = json.Marshal(item.Content.Embedded); err != nil {
	//		return nil, errors.Wrapf(err, onSave+"can't marshal .Embedded(%#v)", item.Content.Embedded)
	//	}
	//}
	//
	//tagsBytes := []byte{} // to to satisfy NOT NULL constraint
	//if len(item.Content.Tags) > 0 {
	//	if tagsBytes, err = json.Marshal(item.Content.Tags); err != nil {
	//		return nil, errors.Wrapf(err, onSave+"can't marshal .Tags(%#v)", item.Content.Tags)
	//	}
	//}
	//
	//// TODO!!! append to .History
	//
	//historyBytes := []byte{} // to satisfy NOT NULL constraint
	//if len(item.History) > 0 {
	//	historyBytes, err = json.Marshal(item.History)
	//	if err != nil {
	//		return nil, errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item.History)
	//	}
	//}
	//
	//// "title", "summary", "type_key", "data", "embedded", "tags",
	//// "issued_id", "owner_id", "viewer_id", "history"
	//values := []interface{}{
	//	item.Content.Title, item.Content.Summary, item.Content.TypeKey, item.Content.Data, embeddedBytes, tagsBytes,
	//	item.IssuedID, item.OwnerID, item.ViewerID, historyBytes}
	//
	//if item.ID == "" {
	//	res, err := personsOp.stmInsert.Exec(values...)
	//	if err != nil {
	//		return nil, errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlInsert, strlib.Stringify(values))
	//	}
	//
	//	idSQLite, err := res.LastInsertId()
	//	if err != nil {
	//		return nil, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, personsOp.sqlInsert, strlib.Stringify(values))
	//	}
	//	item.ID = auth.ID(strconv.FormatInt(idSQLite, 10))
	//
	//} else {
	//	values = append(values, time.Now().Format(time.RFC3339), item.ID)
	//	if _, err := personsOp.stmUpdate.Exec(values...); err != nil {
	//		return nil, errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlUpdate, strlib.Stringify(values))
	//	}
	//
	//}
	//
	//return &item, nil
}

const onChange = "on personsSQLite.Change(): "

func (personsOp *personsSQLite) Change(persons.Item, *crud.Options) (*persons.Item, error) {
	return nil, common.ErrNotImplemented

}

const onRead = "on personsSQLite.Read(): "

func (personsOp *personsSQLite) Read(auth.ID, *crud.Options) (*persons.Item, error) {

	return nil, common.ErrNotImplemented

	//idNum, err := strconv.ParseUint(string(id), 10, 64)
	//if err != nil {
	//	return nil, fmt.Errorf(onRead+"wrong id (%s)", id)
	//}
	//
	//item := persons.Item{ID: id}
	//
	//var embeddedBytes, tagsBytes, historyBytes []byte
	//
	//// "title", "summary", "type_key", "data", "embedded", "tags",
	//// "issued_id", "owner_id", "viewer_id", "history", "updated_at", "created_at"
	//
	//if err = personsOp.stmRead.QueryRow(idNum).Scan(
	//	&item.Content.Title, &item.Content.Summary, &item.Content.TypeKey, &item.Content.Data, &embeddedBytes, &tagsBytes,
	//	&item.IssuedID, &item.OwnerID, &item.ViewerID, &historyBytes, &item.UpdatedAt, &item.CreatedAt); err == sql.ErrNoRows {
	//	return nil, common.ErrNotFound
	//} else if err != nil {
	//	return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, personsOp.sqlRead, idNum)
	//}
	//
	//if len(embeddedBytes) > 0 {
	//	if err = json.Unmarshal(embeddedBytes, &item.Content.Embedded); err != nil {
	//		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embeddedBytes)
	//	}
	//}
	//
	//if len(tagsBytes) > 0 {
	//	if err = json.Unmarshal(tagsBytes, &item.Content.Tags); err != nil {
	//		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tagsBytes)
	//	}
	//}
	//
	//if len(historyBytes) > 0 {
	//	if err = json.Unmarshal(historyBytes, &item.History); err != nil {
	//		return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", historyBytes)
	//	}
	//}
	//
	//return &item, nil
}

const onRemove = "on personsSQLite.Remove()"

func (personsOp *personsSQLite) Remove(auth.ID, *crud.Options) error {
	return common.ErrNotImplemented

	//// TODO!!! rbac check
	//
	//idNum, err := strconv.ParseUint(string(id), 10, 64)
	//if err != nil {
	//	return fmt.Errorf(onRemove+"wrong id (%s)", id)
	//}
	//
	//if _, err = personsOp.stmRemove.Exec(idNum); err != nil {
	//	return errors.Wrapf(err, onRemove+sqllib.CantExec, personsOp.sqlRemove, idNum)
	//}
	//
	//return nil
}

const onList = "on personsSQLite.List()"

func (personsOp *personsSQLite) List(options *crud.Options) ([]persons.Item, error) {
	return nil, common.ErrNotImplemented
	//
	//var termSQL selectors.TermSQL
	//
	//if selector := options.GetSelector(); selector != nil {
	//	var ok bool
	//	if termSQL, ok = selector.(selectors.TermSQL); !ok {
	//		return nil, fmt.Errorf(onList+": wrong selector: %#v", selector)
	//	}
	//}
	//
	//query := sqllib.SQLList(personsOp.table, fieldsToListStr, termSQL.Condition, options)
	//stm, err := personsOp.db.Prepare(query)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	//}
	//
	////l.Infof("%s / %#v\n%s", condition, values, query)
	//
	//rows, err := stm.Query(termSQL.Values...)
	//if err == sql.ErrNoRows {
	//	return nil, nil
	//} else if err != nil {
	//	return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, query, termSQL.Values)
	//}
	//defer rows.Close()
	//
	//var items []persons.Item
	//
	//for rows.Next() {
	//	var idNum int64
	//	var item persons.Item
	//	var embeddedBytes, tagsBytes, historyBytes []byte
	//
	//	// "title", "summary", "type_key", "data", "embedded", "tags",
	//	// "issued_id", "owner_id", "viewer_id", "history", "updated_at", "created_at",
	//	// "id"
	//
	//	if err := rows.Scan(
	//		&item.Content.Title, &item.Content.Summary, &item.Content.TypeKey, &item.Content.Data, &embeddedBytes, &tagsBytes,
	//		&item.IssuedID, &item.OwnerID, &item.ViewerID, &historyBytes, &item.UpdatedAt, &item.CreatedAt,
	//		&idNum); err != nil {
	//		return items, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, query, termSQL.Values)
	//	}
	//
	//	if len(embeddedBytes) > 0 {
	//		if err = json.Unmarshal(embeddedBytes, &item.Content.Embedded); err != nil {
	//			return items, errors.Wrapf(err, onList+": can't unmarshal .Embedded (%s)", embeddedBytes)
	//		}
	//	}
	//
	//	if len(tagsBytes) > 0 {
	//		if err = json.Unmarshal(tagsBytes, &item.Content.Tags); err != nil {
	//			return items, errors.Wrapf(err, onList+": can't unmarshal .Tags (%s)", tagsBytes)
	//		}
	//	}
	//
	//	if len(historyBytes) > 0 {
	//		if err = json.Unmarshal(historyBytes, &item.History); err != nil {
	//			return items, errors.Wrapf(err, onList+": can't unmarshal .History (%s)", historyBytes)
	//		}
	//	}
	//
	//	item.ID = auth.ID(strconv.FormatInt(idNum, 10))
	//	items = append(items, item)
	//}
	//
	//if err = rows.Err(); err != nil {
	//	return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, termSQL.Values)
	//}
	//
	//return items, nil
}

func (personsOp *personsSQLite) Close() error {
	return errors.Wrap(personsOp.db.Close(), "on personsSQLite.Close()")
}

func (personsOp *personsSQLite) HasEmail(email string) (selectors.Term, error) {
	return nil, common.ErrNotImplemented

}

func (personsOp *personsSQLite) HasNickname(nickname string) (selectors.Term, error) {
	return nil, common.ErrNotImplemented

}
