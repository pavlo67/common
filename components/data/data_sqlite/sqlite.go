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
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

var fieldsToUpdate = []string{"url", "type", "title", "summary", "embedded", "tags", "details"}
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToInsert = append(fieldsToUpdate, "source", "source_key", "source_time", "source_data", "export_id")
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToRead = append(fieldsToInsert, "created_at", "updated_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ data.Operator = &dataSQLite{}
var _ crud.Cleaner = &dataSQLite{}

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
		sqlRemove: "DELETE FROM " + table + " where ID = ?",

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

//func sqlList(table, condition string, options *crud.GetOptions) string {
//	if strings.TrimSpace(condition) != "" {
//		condition = " WHERE " + condition
//	}
//
//	var limit string
//
//	order := "created_at DESC"
//	if options != nil {
//		if len(options.OrderBy) > 0 {
//			order = strings.Join(options.OrderBy, ", ")
//		}
//
//		if options.Limit0+options.Limit1 > 0 {
//			limit = " LIMIT " + strconv.FormatUint(options.Limit0, 10)
//			if options.Limit1 > 0 {
//				limit += ", " + strconv.FormatUint(options.Limit1, 10)
//			}
//		}
//	}
//
//	return "SELECT " + fieldsToListStr + " FROM " + table + condition + " ORDER BY " + order + limit
//}
//
//func sqlCount(table, condition string, _ *crud.GetOptions) string {
//	query := "SELECT COUNT(*) FROM " + table
//
//	if strings.TrimSpace(condition) != "" {
//		return query + " WHERE " + condition
//	}
//
//	return query
//}

const onSave = "on dataSQLite.Save(): "

func (dataOp *dataSQLite) Save(items []data.Item, _ *crud.SaveOptions) ([]common.ID, error) {
	var ids []common.ID

	for _, item := range items {

		//l.Info(item.CreatedAt.Format(time.RFC3339))

		var embedded, tags, details string

		if item.Embedded != nil {
			embeddedBytes, err := json.Marshal(item.Embedded)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Embedded)
			}
			embedded = string(embeddedBytes)
		}

		if item.Tags != nil {
			tagsBytes, err := json.Marshal(item.Tags)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Tags)
			}
			tags = string(tagsBytes)
		}

		if item.Details != nil {
			detailsBytes, err := json.Marshal(item.Details)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Details)
			}
			details = string(detailsBytes)
		}

		if item.ID != "" {
			values := []interface{}{
				item.URL, item.TypeKey, item.Title, item.Summary, embedded, tags, details, item.ID,
			}

			_, err := dataOp.stmUpdate.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlUpdate, values)
			}

			if dataOp.taggerOp != nil {
				err = dataOp.taggerOp.ReplaceTags(dataOp.interfaceKey, item.ID, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .ReplaceTags(%#v)", item.Tags)
				}
			}

			ids = append(ids, item.ID)

		} else {
			values := []interface{}{
				item.URL, item.TypeKey, item.Title, item.Summary, embedded, tags, details,
				item.Origin.Source, item.Origin.Key, item.Origin.Time, item.Origin.Data, item.ExportID,
			}

			res, err := dataOp.stmInsert.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, values)
			}

			idSQLite, err := res.LastInsertId()
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlInsert, values)
			}
			id := common.ID(strconv.FormatInt(idSQLite, 10))

			if dataOp.taggerOp != nil && len(item.Tags) > 0 {
				err = dataOp.taggerOp.AddTags(dataOp.interfaceKey, id, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .AddTags(%#v)", item.Tags)
				}
			}

			ids = append(ids, id)
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
	var embedded, tags, createdAt string
	var sourceTimePtr, updatedAtPtr *string

	err = dataOp.stmRead.QueryRow(idNum).Scan(
		&item.URL, &item.TypeKey, &item.Title, &item.Summary, &embedded, &tags, &item.DetailsRaw,
		&item.Origin.Source, &item.Origin.Key, &sourceTimePtr, &item.Origin.Data, &item.ExportID,
		&createdAt, &updatedAtPtr,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, idNum)
	}

	item.Status.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAt)
	}

	//l.Info(createdAt)
	//l.Info(item.CreatedAt.Format(time.RFC3339))

	if updatedAtPtr != nil {
		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
		}
		item.Status.UpdatedAt = &updatedAt
	}

	if sourceTimePtr != nil {
		sourceTime, err := time.Parse(time.RFC3339, *sourceTimePtr)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't parse .SourceTime (%s)", *sourceTimePtr)
		}
		item.Origin.Time = &sourceTime
	}

	if len(tags) > 0 {
		err = json.Unmarshal([]byte(tags), &item.Tags)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
		}
	}

	if len(embedded) > 0 {
		err = json.Unmarshal([]byte(embedded), &item.Embedded)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
		}
	}

	return &item, nil
}

const onDetails = "on dataSQLite.Details(): "

func (dataOp *dataSQLite) SetDetails(item *data.Item) error {
	if item == nil {
		return errors.New(onDetails + "nil item")
	}

	if len(item.DetailsRaw) < 1 {
		item.Details = nil
		return nil
	}

	switch item.TypeKey {
	case data.TypeKeyString:
		item.Details = string(item.DetailsRaw)

	case data.TypeKeyTest:
		item.Details = &data.Test{}
		err := json.Unmarshal(item.DetailsRaw, item.Details)
		if err != nil {
			return errors.Wrapf(err, onDetails+"can't .Unmarshal(%#v)", item.DetailsRaw)
		}

	default:

		// TODO: remove this kostyl
		item.Details = string(item.DetailsRaw)

		// return errors.Errorf(onDetails+"unknown item.TypeKey(%s) for item.DetailsRaw(%s)", item.TypeKey, item.DetailsRaw)

	}

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
			return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter", afterIDStr, err)
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

	//termUpd := selectors.Binary(selectors.Eq, "export_id", selectors.Value{""})
	//if term != nil {
	//	termUpd = logic.AND(term, termUpd)
	//}
	//
	//condition, values, err := selectors_sql.Use(termUpd)
	//if err != nil {
	//	return nil, errors.Errorf(onExport+"wrong term to update export_id's (%#v): %s", termUpd, err)
	//}
	//condition = " WHERE " + condition
	//
	//query := "UPDATE " + dataOp.table + " SET export_id = id " + condition
	//dataOp.db.Exec(query, values...)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onExport+sqllib.CantExec, query, values)
	//}
	//
	//termEx := selectors.Binary(selectors.Ne, "export_id", selectors.Value{""})
	//if term == nil {
	//	term = termEx
	//} else {
	//	term = logic.AND(term, termEx)
	//}

	return dataOp.List(term, options)
}

const onList = "on dataSQLite.List()"

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
		var embedded, tags, createdAt string
		var sourceTimePtr, updatedAtPtr *string

		err := rows.Scan(
			&idNum, &item.URL, &item.TypeKey, &item.Title, &item.Summary, &embedded, &tags, &item.DetailsRaw,
			&item.Origin.Source, &item.Origin.Key, &sourceTimePtr, &item.Origin.Data, &item.ExportID,
			&createdAt, &updatedAtPtr,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if item.Status.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return items, errors.Wrapf(err, onList+"can't parse .CreatedAt (%s)", createdAt)
		}

		if updatedAtPtr != nil {
			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
			}
			item.Status.UpdatedAt = &updatedAt
		}

		if sourceTimePtr != nil {
			sourceTime, err := time.Parse(time.RFC3339, *sourceTimePtr)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't parse .SourceTime (%s)", *sourceTimePtr)
			}
			item.Origin.Time = &sourceTime
		}

		if len(tags) > 0 {
			if err = json.Unmarshal([]byte(tags), &item.Tags); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
			}
		}

		if len(embedded) > 0 {
			if err = json.Unmarshal([]byte(embedded), &item.Embedded); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Embedded (%s)", embedded)
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

const onIDs = "on dataSQLite.IDs()"

func (dataOp *dataSQLite) ids(condition string, values []interface{}) ([]interface{}, error) {
	if strings.TrimSpace(condition) != "" {
		condition = " WHERE " + condition
	}

	query := "SELECT id FROM " + dataOp.table + condition
	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onIDs+": can't db.Prepare(%s)", query)
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onIDs+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var ids []interface{}

	for rows.Next() {
		var id common.ID

		err := rows.Scan(&id)
		if err != nil {
			return ids, errors.Wrapf(err, onIDs+sqllib.CantScanQueryRow, query, values)
		}

		ids = append(ids, id)
	}
	err = rows.Err()
	if err != nil {
		return ids, errors.Wrapf(err, onIDs+": "+sqllib.RowsError, query, values)
	}

	return ids, nil
}

const onClean = "on dataSQLite.Clean(): "

func (dataOp *dataSQLite) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	var termTags *selectors.Term

	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return errors.Errorf(onClean+"wrong selector (%#v): %s", term, err)
	}

	query := dataOp.sqlClean

	if strings.TrimSpace(condition) != "" {
		ids, err := dataOp.ids(condition, values)
		if err != nil {
			return errors.Wrap(err, onClean+"can't dataOp.ids(condition, values)")
		}
		termTags = logic.AND(selectors.In("key", dataOp.interfaceKey), selectors.In("id", ids...))

		query += " WHERE " + condition

	} else {
		termTags = selectors.In("key", dataOp.interfaceKey) // TODO!!! correct field key

	}

	_, err = dataOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	if dataOp.taggerCleaner != nil {
		err = dataOp.taggerCleaner.Clean(termTags, nil)
		if err != nil {
			return errors.Wrap(err, onClean)
		}
	}

	return err
}
