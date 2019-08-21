package flow_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/sqllib"
	"github.com/pavlo67/constructor/components/processor/importer"
	"github.com/pavlo67/constructor/components/structura/content"

	"github.com/pavlo67/constructor/applications/flow"
)

const limitDefault = 200

var tableFlow = "flow"
var tableTags = "tags"
var tableSources = "sources"

var fieldsToList = []string{"id", "source_time", "source_url", "title", "summary", "tags"}
var fieldsToListStr = strings.Join(fieldsToRead, ", ")

var fieldsToRead = []string{"flow.source_id", "source_time", "source_url", "title", "summary", "details", "href", "embedded", "tags"}
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToSave = []string{"source_id", "source_time", "source_url", "title", "summary", "details", "href", "embedded", "tags", "source_key", "origin"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var fieldsToSaveTag = []string{"tag", "flow_id"}
var fieldsToSaveTagStr = strings.Join(fieldsToSaveTag, ", ")

var _ flow.Operator = &flowSQLite{}

type flowSQLite struct {
	limit int
	db    *sql.DB

	stmListAll, stmListByTag, stmListSourcesByTag, stmListByURL, stmRead, stmURLs, stmTags, stmHas, stmSave, stmSaveTag, stmRemove *sql.Stmt
	sqlListAll, sqlListByTag, sqlListSourcesByTag, sqlListByURL, sqlRead, sqlURLs, sqlTags, sqlHas, sqlSave, sqlSaveTag, sqlRemove string
}

const onNew = "on flowSQLite.New(): "

func New(db *sql.DB, limit int) (flow.Operator, error) {
	if db == nil {
		return nil, errors.New(onNew + "no db")
	}

	if limit <= 0 {
		limit = limitDefault
	}

	flowOp := flowSQLite{
		db:    db,
		limit: limit,

		sqlListAll:   "SELECT " + fieldsToListStr + " FROM " + tableFlow + " WHERE saved_at <= ? ORDER BY saved_at DESC LIMIT " + strconv.Itoa(limit),
		sqlListByURL: "SELECT " + fieldsToListStr + " FROM " + tableFlow + " WHERE source_url = ? AND saved_at <= ? ORDER BY saved_at DESC LIMIT " + strconv.Itoa(limit),
		sqlListByTag: "SELECT " + fieldsToListStr + " FROM " + tableTags + " JOIN " + tableFlow + " ON flow_id = flow." +
			"id WHERE tag = ? AND flow.saved_at <= ? ORDER BY flow.saved_at DESC LIMIT " + strconv.Itoa(limit),
		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + tableFlow + " WHERE id = ?",

		sqlURLs: "SELECT source_url, count(*) AS cnt FROM " + tableFlow + "                                               GROUP BY source_url ORDER BY cnt DESC",
		sqlTags: "SELECT tag,        count(*) AS cnt FROM " + tableTags + " JOIN " + tableFlow + " on flow_id   = flow.id GROUP BY tag        ORDER BY cnt DESC",

		sqlHas: "SELECT COUNT(*) FROM " + tableFlow + " WHERE source_id = ? AND source_key = ?",

		sqlSave:    "INSERT INTO " + tableFlow + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlSaveTag: "INSERT INTO " + tableTags + " (" + fieldsToSaveTagStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSaveTag))[1:] + ")",

		sqlRemove: "DELETE FROM " + tableFlow + " where ID = ?",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&flowOp.stmListAll, flowOp.sqlListAll},
		{&flowOp.stmListByURL, flowOp.sqlListByURL},
		{&flowOp.stmListByTag, flowOp.sqlListByTag},
		{&flowOp.stmRead, flowOp.sqlRead},
		{&flowOp.stmURLs, flowOp.sqlURLs},
		{&flowOp.stmTags, flowOp.sqlTags},
		{&flowOp.stmHas, flowOp.sqlHas},
		{&flowOp.stmSave, flowOp.sqlSave},
		{&flowOp.stmSaveTag, flowOp.sqlSaveTag},
		{&flowOp.stmRemove, flowOp.sqlRemove},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &flowOp, nil
}

func (flowOp *flowSQLite) List(errTitle, sqlQuery string, stm *sql.Stmt, values []interface{}, before *time.Time, options *content.GetOptions) ([]content.Brief, error) {
	if before == nil {
		values = append(values, time.Now())
	} else {
		values = append(values, *before)
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, errTitle+sqllib.CantQuery, sqlQuery, values)
	}
	defer rows.Close()

	var briefs []content.Brief

	for rows.Next() {
		brief := content.Brief{Info: common.Info{}}

		var sourceTime *time.Time
		var sourceURL, tags string

		err = rows.Scan(&brief.ID, &sourceTime, &sourceURL, &brief.Title, &brief.Summary, &tags)
		if err != nil {
			return briefs, errors.Wrapf(err, errTitle+sqllib.CantScanQueryRow, sqlQuery, values)
		}

		if sourceTime != nil {
			brief.Info["source_time"] = *sourceTime
		}
		brief.Info["source_url"] = sourceURL
		brief.Info["tags"] = strings.Split(tags, "\n")

		briefs = append(briefs, brief)
	}
	err = rows.Err()
	if err != nil {
		return briefs, errors.Wrapf(err, errTitle+": "+sqllib.RowsError, sqlQuery, values)
	}

	return briefs, nil
}

func (flowOp *flowSQLite) ListAll(before *time.Time, options *content.GetOptions) ([]content.Brief, error) {
	return flowOp.List("on flowSQLite.ListAll(): ", flowOp.sqlListAll, flowOp.stmListAll, nil, before, options)
}

func (flowOp *flowSQLite) ListByURL(sourceURL string, before *time.Time, options *content.GetOptions) ([]content.Brief, error) {
	return flowOp.List("on flowSQLite.ListByURL(): ", flowOp.sqlListByURL, flowOp.stmListByURL, []interface{}{sourceURL}, before, options)
}

func (flowOp *flowSQLite) ListByTag(tag string, before *time.Time, options *content.GetOptions) ([]content.Brief, error) {
	// TODO: list sources also

	return flowOp.List("on flowSQLite.ListByTag(): ", flowOp.sqlListByTag, flowOp.stmListByTag, []interface{}{tag}, before, options)
}

const onRead = "on flowSQLite.Read(): "

func (flowOp *flowSQLite) Read(idStr common.ID, options *content.GetOptions) (*importer.Item, error) {
	if len(idStr) < 1 {
		return nil, errors.New(onRead + "empty ID")
	}

	id, err := strconv.ParseUint(string(idStr), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong ID (%s)", idStr)
	}

	var item importer.Item
	var embedded, tags string

	err = flowOp.stmRead.QueryRow(id).Scan(&item.SourceID, &item.SourceTime, &item.SourceURL, &item.Title, &item.Summary, &item.Details, &item.Href, &embedded, &tags)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, flowOp.sqlRead, id)
	}

	item.Tags = strings.Split(tags, "\n")
	err = json.Unmarshal([]byte(embedded), &item.Embedded)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
	}

	return &item, nil
}

const onSources = "on flowSQLite.Sources(): "

func (flowOp *flowSQLite) URLs(_ *content.GetOptions) ([]flow.Part, error) {
	rows, err := flowOp.stmURLs.Query()
	if err != nil {
		return nil, errors.Errorf(onSources+sqllib.CantQuery, flowOp.sqlURLs, nil)
	}
	defer rows.Close()

	var parts []flow.Part

	for rows.Next() {
		var url string
		var cnt uint64
		err = rows.Scan(&url, &cnt)
		if err != nil {
			return nil, errors.Errorf(onSources+sqllib.CantScanQueryRow, flowOp.sqlURLs, nil)
		}

		parts = append(parts, flow.Part{url, cnt})
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Errorf(onSources+sqllib.CantScanQueryRow, flowOp.sqlURLs, nil)

	}

	return parts, nil
}

const onTags = "on flowSQLite.Tags(): "

func (flowOp *flowSQLite) Tags(*content.GetOptions) ([]flow.Part, error) {
	rows, err := flowOp.stmTags.Query()
	if err != nil {
		return nil, errors.Errorf(onTags+sqllib.CantQuery, flowOp.sqlTags, nil)
	}
	defer rows.Close()

	var parts []flow.Part

	for rows.Next() {
		var tag string
		var cnt uint64
		err = rows.Scan(&tag, &cnt)
		if err != nil {
			return nil, errors.Errorf(onTags+sqllib.CantScanQueryRow, flowOp.sqlTags, nil)
		}
		parts = append(parts, flow.Part{tag, cnt})
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Errorf(onTags+sqllib.CantScanQueryRow, flowOp.sqlTags, nil)

	}

	return parts, nil
}

const onHas = "on flowSQLite.Has(): "

func (flowOp *flowSQLite) Has(originKey importer.OriginKey) (bool, error) {
	if len(originKey.SourceKey) < 1 { // || len(originKey.SourceID) < 1
		return false, errors.New(onHas + "empty ID")
	}

	values := []interface{}{originKey.SourceID, originKey.SourceKey}
	var cnt uint64

	err := flowOp.stmHas.QueryRow(values...).Scan(&cnt)
	if err != nil {
		return false, errors.Wrapf(err, onHas+sqllib.CantScanQueryRow, flowOp.sqlHas, values)
	}

	return cnt > 0, nil
}

type tagItem struct {
	tag    string
	flowID *common.ID
}

const onSaveTags = "on flowSQLite.saveTags(): "

func (flowOp *flowSQLite) saveTags(tagItems []tagItem) error {

	var errs common.Errors

	for _, tagItem := range tagItems {
		values := []interface{}{tagItem.tag, tagItem.flowID}

		_, err := flowOp.stmSaveTag.Exec(values...)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onSaveTags+sqllib.CantExec, flowOp.sqlSaveTag, values))
		}

	}

	return errs.Err()
}

const onSave = "on flowSQLite.Save(): "

func (flowOp *flowSQLite) Save(items []importer.Item, options *content.SaveOptions) ([]common.ID, error) {
	var ids []common.ID
	var errs common.Errors

	for _, item := range items {
		embedded, err := json.Marshal(item.Embedded)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+"can't .Marshal: %s", item.Embedded)).Err()
		}

		values := []interface{}{item.SourceID, item.SourceTime, item.SourceURL, item.Title, item.Summary, item.Details, item.Href, embedded, strings.Join(item.Tags, "\n"), item.SourceKey, item.Origin}

		res, err := flowOp.stmSave.Exec(values...)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+sqllib.CantExec, flowOp.sqlSave, values)).Err()
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, flowOp.sqlSave, values)).Err()
		}
		id := common.ID(strconv.FormatInt(idSQLite, 10))

		var tagItems []tagItem
		for _, tag := range item.Tags {
			tagItems = append(tagItems, tagItem{tag, &id})
		}

		ids = append(ids, id)
		errs = errs.Append(flowOp.saveTags(tagItems))
	}

	return ids, errs.Err()
}

func (flowOp *flowSQLite) Close() error {
	return errors.Wrap(flowOp.db.Close(), "on flowSQLite.Close()")
}

func (flowOp *flowSQLite) Clean() error {
	_, err1 := flowOp.db.Exec("TRUNCATE " + tableFlow)
	_, err2 := flowOp.db.Exec("TRUNCATE " + tableTags)
	_, err3 := flowOp.db.Exec("TRUNCATE " + tableSources)

	return common.MultiError(err1, err2, err3).Err()
}

//const onRemove = "on flowSQLite.Remove()"
//
//func (flowOp *flowSQLite) Remove(sourceIDs []basis.ID, before *time.Time, options *structura.RemoveOptions) error {
//		var err error
//		var values []interface{}
//		var orderAndLimit, condition, conditionCompleted string
//
//		if options != nil {
//			condition, values, err = selectors.Mysql("", options.Selector)
//			if err != nil {
//				return crud.Result{}, errors.Wrapf(err, onDelete+"bad selector ('%#v')", options.Selector)
//			}
//
//			conditionCompleted = condition
//			if strings.TrimSpace(conditionCompleted) != "" {
//				conditionCompleted = " where " + conditionCompleted
//			}
//
//			orderAndLimit = mysqllib.OrderAndLimit(options.SortBy, options.Limits)
//		}
//
//		if strings.TrimSpace(condition) != "" {
//			condition = "where " + condition
//		}
//
//		sqlQuery := dsOp.sqlDelete + " " + condition + " " + orderAndLimit
//		res, err := dsOp.db.Exec(sqlQuery, values...)
//		if err != nil {
//			return crud.Result{}, errors.Wrapf(err, onDelete+"can't exec SQL: %s, %s", sqlQuery, values)
//		}
//		cnt, err := res.RowsAffected()
//		if err != nil {
//			return crud.Result{}, errors.Wrapf(err, onDelete+"can't get RowsAffected(): %s, %s", sqlQuery, values)
//		}
//		return crud.Result{cnt}, nil
//	}
//}

//const onLastKey = "on datastoreMySQL.LastKey()"
//
//func (dsOp *datastoreMySQL) LastKey(class flow.Type, options *crud.ReadOptions) (string, error) {
//
//	// TODO: use options!!!
//
//	values := []interface{}{string(class)}
//	rows, err := dsOp.stmLastKey.Query(values...)
//	if err == sql.ErrNoRows {
//		return "", nil
//	} else if err != nil {
//		return "", errors.Wrapf(err, onLastKey+"can't query (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//	defer rows.Close()
//	if rows.Next() {
//		var lastKey string
//		err = rows.Scan(&lastKey)
//		if err != nil {
//			return "", errors.Wrapf(err, onLastKey+"can't scan query row (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//		}
//		return lastKey, nil
//	}
//	err = rows.Err()
//	if err != nil {
//		return "", errors.Wrapf(err, onLastKey+"on rows.Err() (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//
//	return "", nil
//}
//
