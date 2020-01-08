package data_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/common/libraries/strlib"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"
	"github.com/pavlo67/workshop/common/types"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

var fieldsToInsert = []string{"data_key", "url", "title", "summary", "embedded", "tags", "type_key", "details", "history"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToUpdate = append(fieldsToInsert, "updated_at")
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToRead = append(fieldsToUpdate, "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ data.Operator = &dataSQLite{}

type dataSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlList, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove, stmList           *sql.Stmt

	taggerOp      tagger.Operator
	interfaceKey  joiner.InterfaceKey
	taggerCleaner crud.Cleaner
}

const onNew = "on dataSQLite.New(): "

func New(access config.Access, table string, interfaceKey joiner.InterfaceKey, taggerOp tagger.Operator, taggerCleaner crud.Cleaner) (data.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = data.CollectionDefault
	}

	dataOp := dataSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		sqlRemove: "DELETE FROM " + table + " where id = ?",

		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",
		sqlList: sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at DESC"}}),

		sqlClean: "DELETE FROM " + table,

		taggerOp:      taggerOp,
		interfaceKey:  interfaceKey,
		taggerCleaner: taggerCleaner,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&dataOp.stmInsert, dataOp.sqlInsert},
		{&dataOp.stmUpdate, dataOp.sqlUpdate},
		{&dataOp.stmRemove, dataOp.sqlRemove},

		{&dataOp.stmRead, dataOp.sqlRead},
		{&dataOp.stmList, dataOp.sqlList},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &dataOp, &dataOp, nil
}

const onSave = "on dataSQLite.Save(): "

func (dataOp *dataSQLite) Save(items []data.Item, _ *crud.SaveOptions) ([]common.ID, error) {
	var ids []common.ID

	for _, item := range items {

		//l.Info(item.SentAt.Format(time.RFC3339))

		var err error

		var embedded, tags, details, history interface{}

		if len(item.Embedded) > 0 {
			embedded, err = json.Marshal(item.Embedded)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't marshal .Embedded(%#v)", item.Embedded)
			}
		} else {
			embedded = "" // to satisfy "NOT NULL" constraint
		}

		if len(item.Tags) > 0 {
			tags, err = json.Marshal(item.Tags)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't marshal .Tags(%#v)", item.Tags)
			}
		} else {
			tags = "" // to satisfy "NOT NULL" constraint
		}

		if item.Details != nil {
			details, err = json.Marshal(item.Details)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't marshal .Details(%#v)", item.Details)
			}
		} else {
			details = "" // to satisfy "NOT NULL" constraint
		}

		// TODO!!! append to .History

		if len(item.History) > 0 {
			history, err = json.Marshal(item.History)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item.History)
			}
		} else {
			history = "" // to satisfy "NOT NULL" constraint
		}

		if item.ID == "" {

			values := []interface{}{item.Key, item.URL, item.Title, item.Summary, embedded, tags, item.TypeKey, details, history}

			res, err := dataOp.stmInsert.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, strlib.Stringify(values))
			}

			idSQLite, err := res.LastInsertId()
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlInsert, strlib.Stringify(values))
			}
			id := common.ID(strconv.FormatInt(idSQLite, 10))

			if dataOp.taggerOp != nil && len(item.Tags) > 0 {
				err = dataOp.taggerOp.AddTags(dataOp.interfaceKey, id, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .AddTags(%#v)", item.Tags)
				}
			}

			ids = append(ids, id)

		} else {

			values := []interface{}{item.Key, item.URL, item.Title, item.Summary, embedded, tags, item.TypeKey, details, history, time.Now().Format(time.RFC3339), item.ID}

			_, err := dataOp.stmUpdate.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlUpdate, strlib.Stringify(values))
			}

			if dataOp.taggerOp != nil {
				err = dataOp.taggerOp.ReplaceTags(dataOp.interfaceKey, item.ID, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .ReplaceTags(%#v)", item.Tags)
				}
			}

			ids = append(ids, item.ID)
		}
	}

	return ids, nil
}

const onRead = "on dataSQLite.Read(): "

func (dataOp *dataSQLite) Read(id common.ID, _ *crud.GetOptions) (*data.Item, error) {
	if len(id) < 1 {
		return nil, errors.New(onRead + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong ID (%s)", id)
	}

	item := data.Item{ID: id}
	var embedded, tags, history, updatedAt, createdAt []byte

	err = dataOp.stmRead.QueryRow(idNum).Scan(
		&item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.TypeKey, &item.DetailsRaw, &history, &updatedAt, &createdAt,
	)

	// TODO: read updated_at & created_at

	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, idNum)
	}

	if len(tags) > 0 {
		err = json.Unmarshal(tags, &item.Tags)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
		}
	}

	if len(embedded) > 0 {
		err = json.Unmarshal(embedded, &item.Embedded)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
		}
	}

	if len(history) > 0 {
		err = json.Unmarshal(history, &item.History)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", history)
		}
	}

	return &item, nil
}

const onDetails = "on dataSQLite.Details(): "

func (dataOp *dataSQLite) SetDetails(item *data.Item) error {
	if item == nil {
		return errors.New(onDetails + "nil item")
	}

	// l.Infof("11111111111 %s %s %t", item.DetailsRaw, item.TypeKey, item.TypeKey == data.TypeKeyTest)

	if len(item.DetailsRaw) < 1 {
		item.Details = nil
		return nil
	}

	switch item.TypeKey {
	case types.KeyString:
		item.Details = string(item.DetailsRaw)

	case data.TypeKeyTest:
		item.Details = &data.Test{}
		err := json.Unmarshal(item.DetailsRaw, item.Details)
		if err != nil {
			return errors.Wrapf(err, onDetails+"can't .Unmarshal(%#v)", item.DetailsRaw)
		}

	default:

		// TODO: remove the kostyl
		item.Details = string(item.DetailsRaw)

		// return errors.Errorf(onDetails+"unknown item.TypeKey(%s) for item.DetailsRaw(%s)", item.TypeKey, item.DetailsRaw)

	}

	// l.Infof("11111111111 %#v", item.Details)

	return nil
}

const onRemove = "on dataSQLite.Remove()"

func (dataOp *dataSQLite) Remove(id common.ID, _ *crud.RemoveOptions) error {
	if len(id) < 1 {
		return errors.New(onRemove + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onRemove+"wrong ID (%s)", id)
	}

	_, err = dataOp.stmRemove.Exec(idNum)
	if err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, dataOp.sqlRemove, idNum)
	}

	if dataOp.taggerOp != nil {
		err = dataOp.taggerOp.ReplaceTags(dataOp.interfaceKey, id, nil, nil)
		if err != nil {
			return errors.Wrapf(err, onRemove+": can't .ReplaceTags(%#v)", nil)
		}
	}

	return nil
}

const onExport = "on dataSQLite.Export()"

func (dataOp *dataSQLite) Export(afterIDStr string, options *crud.GetOptions) ([]data.Item, error) {
	// TODO: remove limits
	// if options != nil {
	//	options.Limits = nil
	// }

	afterIDStr = strings.TrimSpace(afterIDStr)

	var term *selectors.Term

	var afterID int
	if afterIDStr != "" {
		var err error
		afterID, err = strconv.Atoi(afterIDStr)
		if err != nil {
			return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter: %s", afterIDStr, err)
		}

		// TODO!!! term with some item's autoincrement if original .ID isn't it (using .ID to find corresponding autoincrement value)
		term = selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
	}

	// TODO!!! order by some item's autoincrement if original .ID isn't it
	if options == nil {
		options = &crud.GetOptions{OrderBy: []string{"id"}}
	} else {
		options.OrderBy = []string{"id"}
	}

	return dataOp.List(term, options)
}

const onList = "on dataSQLite.ListTags()"

func (dataOp *dataSQLite) List(term *selectors.Term, options *crud.GetOptions) ([]data.Item, error) {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
	}

	query := dataOp.sqlList
	stm := dataOp.stmList

	if condition != "" || options != nil {
		query = sqllib.SQLList(dataOp.table, fieldsToListStr, condition, options)
		stm, err = dataOp.db.Prepare(query)
		if err != nil {
			return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
		}
	}

	//l.Infof("%s / %#v\n%s", condition, values, query)

	rows, err := stm.Query(values...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []data.Item

	for rows.Next() {
		var idNum int64
		var item data.Item
		var embedded, tags, history, updatedAt, createdAt []byte

		err := rows.Scan(
			&idNum, &item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.TypeKey, &item.DetailsRaw, &history, &updatedAt, &createdAt,
		)

		// TODO: read updated_at & created_at

		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if len(tags) > 0 {
			if err = json.Unmarshal(tags, &item.Tags); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
			}
		}

		if len(embedded) > 0 {
			if err = json.Unmarshal(embedded, &item.Embedded); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Embedded (%s)", embedded)
			}
		}

		if len(history) > 0 {
			err = json.Unmarshal(history, &item.History)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", history)
			}
		}

		item.ID = common.ID(strconv.FormatInt(idNum, 10))
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onCount = "on dataSQLite.Count(): "

func (dataOp *dataSQLite) Count(term *selectors.Term, options *crud.GetOptions) (uint64, error) {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		termStr, _ := json.Marshal(term)
		return 0, errors.Wrapf(err, onCount+": can't selectors_sql.Use(%s)", termStr)
	}

	query := sqllib.SQLCount(dataOp.table, condition, options)
	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return 0, errors.Wrapf(err, onCount+": can't db.Prepare(%s)", query)
	}

	var num uint64

	err = stm.QueryRow(values...).Scan(&num)
	if err != nil {
		return 0, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, query, values)
	}

	return num, nil
}

func (dataOp *dataSQLite) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on dataSQLite.Close()")
}
