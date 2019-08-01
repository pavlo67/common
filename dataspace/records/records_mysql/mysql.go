package records_mysql

//import (
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"regexp"
//	"strconv"
//	"strings"
//
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/partes/libs/mysqllib"
//	"github.com/pavlo67/associatio/auth"
//	"github.com/pavlo67/associatio/basis"
//	"github.com/pavlo67/associatio/starter/config"
//	"github.com/pavlo67/associatio/starter/joiner"
//
//	"github.com/pavlo67/partes/crud"
//	"github.com/pavlo67/partes/crud/selectors"
//	"github.com/pavlo67/associatio/confidenter/groups"
//	"github.com/pavlo67/associatio/confidenter/rights"
//	"github.com/pavlo67/associatio/notebook/links"
//	"github.com/pavlo67/associatio/notebook/notes"
//)
//
//const MaxVarcharLen = 255
//
//// notesMySQL is a struct to implement Item interface using MySQL.
//type notesMySQL struct {
//	grpsOp       groups.Operator
//	dataManagers rights.Managers
//
//	linksOp    links.Operator
//	jointLinks bool
//
//	// generaOp genera.Operator
//
//	dbh   *sql.DB
//	table string
//
//	stmtCreate        *sql.Stmt
//	stmtRead          *sql.Stmt
//	stmtReadForUpdate *sql.Stmt
//	stmtUpdate        *sql.Stmt
//	stmtDelete        *sql.Stmt
//
//	stmtUpdateCountLinked *sql.Stmt
//	stmtUpdateLinks       *sql.Stmt
//
//	stmtImportTo    *sql.Stmt
//	stmtIsGlobal    *sql.Stmt
//	stmtSetGlobalIS *sql.Stmt
//
//	sqlCreate, sqlRead, sqlReadList, sqlReadForUpdate, sqlUpdate, sqlDelete, sqlDeleteAll, sqlUpdateLinks, sqlUpdateCountLinked, sqlIsGlobal, sqlImportTo, sqlSetGlobalIS string
//}
//
//var _ notes.Operator = &notesMySQL{}
//
//var fields = []string{
//	"global_is", notes.GenusFieldName,
//	"author", "name", "visibility",
//	"brief", "content",
//	"links", "tags", "count_linked",
//	"r_view", "r_owner", "managers",
//	"history", "status",
//}
//
//const onNew = "on objectsmysql.NewCRUDOperator()"
//
//func New(
//	mysqlConfig config.ServerAccess,
//	table string,
//	jointLinks bool,
//	grpsOp groups.Operator,
//	linksOp links.Operator,
//	// generaOp genera.Operator,
//	managers rights.Managers) (*notesMySQL, error) {
//
//	dbh, err := mysqllib.ConnectToMysql(mysqlConfig)
//	if err != nil {
//		return nil, errors.Wrapf(err, onNew+": credentials = %#v", mysqlConfig)
//	}
//
//	if grpsOp == nil {
//		l.Warn(onNew + ": no groups.Operator")
//		// return nil, errors.New(onNew + ": no grpsOp...")
//	}
//
//	if linksOp == nil {
//		l.Warn(onNew + ": no links.Operator")
//		// return nil, errors.New(onNew + ": no linksOp...")
//	}
//
//	//if generaOp == nil {
//	//	return nil, errors.New(onNew + ": no generaOp...")
//	//}
//
//	objOp := new(notesMySQL)
//	objOp.jointLinks = jointLinks
//	objOp.grpsOp = grpsOp
//	objOp.linksOp = linksOp
//	// objOp.generaOp = generaOp
//
//	objOp.dataManagers = managers
//	if len(objOp.dataManagers) < 1 {
//		objOp.dataManagers = rights.Managers{rights.Create: basis.Anyone}
//	}
//
//	objOp.dbh = dbh
//	objOp.table = table
//
//	fieldsToCreate := strings.Join(fields, ", ")
//	fieldsToRead := "id, " + strings.Join(fields, ", ") + ", created_at, updated_at"
//	fieldsToUpdate := strings.Join(fields, " = ?, ") + " = ?"
//
//	objOp.sqlCreate = "insert into `" + table + "` (" + fieldsToCreate + ") values (?,?, ?,?,?, ?,?, ?,?,?, ?,?,?, ?,?)"
//
//	objOp.sqlRead = "select " + fieldsToRead + " from `" + table + "` where id = ?"
//	objOp.sqlReadForUpdate = "select author, r_view, r_owner, managers, global_is, links, count_linked from `" + table + "` where id = ?"
//	objOp.sqlUpdate = "update `" + table + "` set " + fieldsToUpdate + " where id = ?"
//	objOp.sqlDelete = "delete from `" + table + "` where id = ?"
//
//	objOp.sqlUpdateLinks = "update `" + table + "` set links = ? where id = ?"
//	objOp.sqlUpdateCountLinked = "update `" + table + "` set count_linked = ? where id = ?"
//
//	objOp.sqlDeleteAll = "delete from `" + table + "` "
//	objOp.sqlReadList = "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "` "
//
//	objOp.sqlImportTo = "update `" + table + "` set `status` = ? where id = ? and r_owner = ?"
//
//	objOp.sqlIsGlobal = "select `id`, `links` from `" + table + "` where global_is = ?"
//	objOp.sqlSetGlobalIS = "update `" + table + "` set global_is = ?, links = ? where id = ? "
//
//	sqlStmts := []mysqllib.SqlStmt{
//		{&objOp.stmtCreate, objOp.sqlCreate},
//		{&objOp.stmtRead, objOp.sqlRead},
//		{&objOp.stmtReadForUpdate, objOp.sqlReadForUpdate},
//		{&objOp.stmtUpdate, objOp.sqlUpdate},
//		{&objOp.stmtDelete, objOp.sqlDelete},
//
//		{&objOp.stmtUpdateCountLinked, objOp.sqlUpdateCountLinked},
//		{&objOp.stmtUpdateLinks, objOp.sqlUpdateLinks},
//
//		{&objOp.stmtImportTo, objOp.sqlImportTo},
//		{&objOp.stmtIsGlobal, objOp.sqlIsGlobal},
//		{&objOp.stmtSetGlobalIS, objOp.sqlSetGlobalIS},
//	}
//
//	for _, sqlStmt := range sqlStmts {
//		if err = mysqllib.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
//			return nil, errors.Wrap(err, onNew)
//		}
//	}
//
//	return objOp, nil
//}
//
//// CRUD -------------------------------------------------------------------------------------
//
//const onCreate = "on notesMySQL.Create"
//
//func (objOp *notesMySQL) Create(userIS auth.ID, o notes.Item) (id string, err error) {
//	rView, rOwner, managers, err := groups.SetRights(userIS, objOp.grpsOp, o.RView, o.ROwner, o.Managers)
//	if err != nil {
//		return "", errors.Wrap(err, onCreate+": can't .SetRights()")
//	}
//
//	var managersStr []byte
//	if managers != nil {
//		if managersStr, err = json.Marshal(managers); err != nil {
//			return "", errors.Wrapf(err, onDelete+": can't json.Marshal($#v)", managers)
//		}
//	}
//
//	linksList := notes.PrepareLinks(userIS, objOp.grpsOp, o.ROwner, nil, o.Links, objOp.jointLinks)
//	linksListCopy := linksList
//	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
//	var jsonLinks []byte
//	if len(linksList) > 0 {
//		jsonLinks, err = json.Marshal(linksList)
//		if err != nil {
//			return "", errors.Wrapf(err, onCreate+": can't marshal .Links(%#v)", linksList)
//		}
//	}
//
//	if strings.TrimSpace(o.Genus) == "" {
//		o.Genus = notes.GenusDefault
//	}
//
//	if len([]rune(o.Author)) > MaxVarcharLen {
//		o.Author = string([]rune(o.Author)[:MaxVarcharLen])
//	}
//
//	values := []interface{}{
//		o.GlobalIS, o.Genus,
//		o.Author, strings.TrimSpace(o.Name), o.Visibility,
//		o.Brief, o.Content,
//		string(jsonLinks), o.Tags, 0,
//		string(rView), string(rOwner), managersStr,
//		o.History, o.Status,
//	}
//
//	res, err := objOp.stmtCreate.Exec(values...)
//	if err != nil {
//		return "", errors.Wrapf(err, onCreate+basis.CantExecQuery, objOp.sqlCreate, values)
//	}
//	idInt64, err := res.LastInsertId()
//	if err != nil {
//		return "", errors.Wrapf(err, onCreate+": can't get LastInsertId() "+basis.CantExecQuery, objOp.sqlCreate, values)
//	}
//
//	o.ID = strconv.FormatInt(idInt64, 10)
//
//	return o.ID, objOp.setLinks(userIS, o.ID, linksListCopy)
//}
//
//const onRead = "on notesMySQL.Read"
//
//func (objOp *notesMySQL) Read(userIS auth.ID, idStr string) (*notes.Item, error) {
//	if len(idStr) < 1 {
//		return nil, errors.Wrap(crud.ErrEmptySelector, onRead)
//	}
//
//	id, err := strconv.ParseUint(idStr, 10, 64)
//	if err != nil {
//		return nil, errors.Wrap(crud.ErrBadSelector, onRead+": "+idStr)
//	}
//
//	var o notes.Item
//	var jsonManagers, jsonLinks []byte
//	err = objOp.stmtRead.QueryRow(id).Scan(
//		&o.ID, &o.GlobalIS, &o.Genus,
//		&o.Author, &o.Name, &o.Visibility,
//		&o.Brief, &o.Content,
//		&jsonLinks, &o.Tags, &o.CountLinked,
//		&o.RView, &o.ROwner, &jsonManagers,
//		&o.History, &o.Status,
//		&o.CreatedAt, &o.UpdatedAt,
//	)
//	if err == sql.ErrNoRows {
//		return nil, errors.Wrap(basis.ErrNotFound, onRead+": "+idStr)
//	} else if err != nil {
//		return nil, errors.Wrapf(err, onRead+": "+mysqllib.CantScanQueryRow, objOp.sqlRead, idStr)
//	}
//
//	if !groups.OneOf(userIS, objOp.grpsOp, o.RView) {
//		return nil, fmt.Errorf(onRead+": %#v tries to get access to private object (%s) with r_view = %s and r_owner = %s", userIS, o.ID, o.ROwner, o.RView)
//	}
//
//	if len(jsonManagers) > 0 {
//		err = json.Unmarshal(jsonManagers, &o.Managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Managers:'%s'", string(jsonManagers))
//		}
//	}
//
//	if len(jsonLinks) > 0 {
//		err = json.Unmarshal(jsonLinks, &o.Links)
//		if err != nil {
//			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Links: '%s'", string(jsonLinks))
//		}
//
//		o.Links = links.Filter(userIS, objOp.grpsOp, o.Links)
//	}
//
//	return &o, nil
//}
//
//const onReadList = "on notesMySQL.ReadList"
//
//func (objOp *notesMySQL) ReadList(userIS auth.ID, options *content.ListOptions) ([]notes.Item, uint64, error) {
//	var err error
//	var values []interface{}
//	var orderAndLimit, condition, conditionCompleted string
//
//	var forAdmin, addGlobalIS bool
//	if options != nil {
//
//		condition, values, err = selectors.Mysql(userIS, options.Selector)
//		if err != nil {
//			return nil, 0, errors.Wrapf(err, ": bad selector ('%#v')", options.Selector)
//		}
//
//		conditionCompleted = condition
//		if strings.TrimSpace(conditionCompleted) != "" {
//			conditionCompleted = " where " + conditionCompleted
//		}
//
//		orderAndLimit = mysqllib.OrderAndLimit(options.SortBy, options.Limits)
//	}
//
//	sqlQuery := objOp.sqlReadList + conditionCompleted + orderAndLimit
//
//	var rows *sql.Rows
//	if forAdmin {
//		// TODO!!! add correct check if user is admin
//		rows, err = objOp.dbh.Query(sqlQuery, values...)
//	} else {
//		_, rows, err = groups.QueryAccessible(objOp.grpsOp, objOp.dbh, userIS, objOp.sqlReadList, condition, orderAndLimit, values)
//	}
//	if err != nil {
//		return nil, 0, err
//	}
//	defer rows.Close()
//
//	var oAll []notes.Item
//	var jsonLinks, jsonManagers []byte
//	for rows.Next() {
//		o := notes.Item{}
//		if err := rows.Scan(
//			&o.ID, &o.GlobalIS, &o.Genus,
//			&o.Author, &o.Name, &o.Visibility,
//			&o.Brief, &o.Content,
//			&jsonLinks, &o.Tags, &o.CountLinked,
//			&o.RView, &o.ROwner, &jsonManagers,
//			&o.History, &o.Status,
//			&o.CreatedAt, &o.UpdatedAt,
//		); err != nil {
//			return oAll, 0, errors.Wrapf(err, onReadList+": "+mysqllib.CantScanQueryRow, sqlQuery, values)
//		}
//		if len(jsonManagers) > 0 {
//			err = json.Unmarshal(jsonManagers, &o.Managers)
//			if err != nil {
//				return nil, 0, errors.Wrapf(err, onReadList+": can't unmarshal jsonManagers: %s", jsonManagers)
//			}
//		}
//		if len(jsonLinks) > 0 {
//			err = json.Unmarshal(jsonLinks, &o.Links)
//			if err != nil {
//				return nil, 0, errors.Wrapf(err, onReadList+": can't unmarshal jsonLinks: %s", jsonLinks)
//			}
//			o.Links = links.Filter(userIS, objOp.grpsOp, o.Links)
//		}
//		if addGlobalIS {
//			var needUpdate = false
//			if o.GlobalIS == "" {
//				needUpdate = true
//				o.GlobalIS = joiner.SystemDomain() + "/object/" + o.ID
//			}
//			var fillID = false
//			o.Links, fillID = notes.FillFilesIDs(o.Links)
//			if fillID {
//				needUpdate = true
//			}
//			if needUpdate {
//				_, err = objOp.Update(userIS, o)
//				if err != nil {
//					return nil, 0, errors.Wrapf(err, "can't set GlobalIS = '%s' for object.TargetID = '%s'", o.GlobalIS, o.ID)
//				}
//			}
//		}
//		oAll = append(oAll, o)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		return oAll, 0, errors.Wrapf(err, onReadList+": "+mysqllib.CantScanQueryRow, sqlQuery, values)
//	}
//
//	var allCount uint64
//	err = objOp.dbh.QueryRow("SELECT FOUND_ROWS()").Scan(&allCount)
//	if err != nil {
//		return nil, 0, errors.Wrapf(err, onReadList+": can't scan ('SELECT FOUND_ROWS()') for sql=", sqlQuery)
//	}
//
//	return oAll, allCount, nil
//}
//
//const onCanUpdate = "on notesMySQL.canUpdate()"
//
//func (objOp *notesMySQL) readForUpdate(id string) (author, rView, rOwner, jsonManagers, globalIS, jsonLinks []byte, countLinked uint, err error) {
//	err = objOp.stmtReadForUpdate.QueryRow(id).Scan(&author, &rView, &rOwner, &jsonManagers, &globalIS, &jsonLinks, &countLinked)
//	if err != nil {
//		return nil, nil, nil, nil, nil, nil, 0, errors.Wrapf(err, onCanUpdate+" (%s)", id)
//	}
//
//	return author, rView, rOwner, jsonManagers, globalIS, jsonLinks, countLinked, nil
//}
//
//const onUpdate = "on notesMySQL.Update()"
//
//// Update changes object's items.User data (accordingly to requester's rights).
//func (objOp *notesMySQL) Update(userIS auth.ID, o notes.Item) (crud.Result, error) {
//	author0, _, _, _, globalIS0, jsonLinks0, countLinked0, err := objOp.readForUpdate(o.ID)
//	if err != nil {
//		return crud.Result{}, errors.Wrap(err, onUpdate)
//	}
//
//	if globalIS0 == nil {
//		globalIS0 = []byte(o.GlobalIS)
//	}
//	if strings.TrimSpace(o.Author) == "" {
//		o.Author = string(author0)
//	} else if len([]rune(o.Author)) > MaxVarcharLen {
//		o.Author = string([]rune(o.Author)[:MaxVarcharLen])
//	}
//
//	rView, rOwner, managers, err := groups.SetRights(userIS, objOp.grpsOp, o.RView, o.ROwner, o.Managers)
//	if err != nil {
//		return crud.Result{}, errors.Wrap(err, onUpdate+": can't .CheckAndUpdateRights()")
//	}
//
//	var managersStr []byte
//	if managers != nil {
//		if managersStr, err = json.Marshal(managers); err != nil {
//			return crud.Result{}, errors.Wrapf(err, onUpdate+": can't json.Marshal($#v)", managers)
//		}
//	}
//
//	var linksListOld []links.Item
//	if len(jsonLinks0) > 0 {
//		err = json.Unmarshal(jsonLinks0, &linksListOld)
//		if err != nil {
//			return crud.Result{}, errors.Wrapf(err, onUpdate+": can't unmarshal linksList: '%s'", string(jsonLinks0))
//		}
//	}
//
//	linksList := notes.PrepareLinks(userIS, objOp.grpsOp, rOwner, linksListOld, o.Links, objOp.jointLinks)
//	linksListCopy := linksList
//	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
//	var jsonLinks []byte
//	if len(linksList) > 0 {
//		jsonLinks, err = json.Marshal(linksList)
//		if err != nil {
//			return crud.Result{}, errors.Wrapf(err, onUpdate+": can't marshal .Links(%#v)", linksList)
//		}
//	}
//
//	if strings.TrimSpace(o.Genus) == "" {
//		o.Genus = notes.GenusDefault
//	}
//
//	values := []interface{}{
//		globalIS0, o.Genus,
//		o.Author, strings.TrimSpace(o.Name), o.Visibility,
//		o.Brief, o.Content,
//		string(jsonLinks), o.Tags, countLinked0,
//		string(rView), string(rOwner), managersStr,
//		o.History, o.Status,
//		o.ID,
//	}
//
//	res, err := objOp.stmtUpdate.Exec(values...)
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onUpdate+": "+basis.CantExecQuery, objOp.sqlUpdate, values)
//	}
//
//	cnt, err := res.RowsAffected()
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onUpdate+": "+mysqllib.CantGetRowsAffected, objOp.sqlUpdate, values)
//	}
//
//	return crud.Result{cnt}, objOp.setLinks(userIS, o.ID, linksListCopy)
//}
//
//const onDelete = "on notesMySQL.DeleteList()"
//
//func (objOp *notesMySQL) Delete(userIS auth.ID, id string) (crud.Result, error) {
//	_, _, rOwner0, _, _, _, _, err := objOp.readForUpdate(id)
//	if err != nil {
//		return crud.Result{}, err
//	}
//
//	// TODO: check rights.Change
//	if !groups.OneOf(userIS, objOp.grpsOp, auth.ID(rOwner0)) {
//		return crud.Result{}, errors.New(onDelete + ": no rights")
//	}
//
//	res, err := objOp.stmtDelete.Exec(id)
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onDelete+": "+basis.CantExecQuery, objOp.sqlDelete, id)
//	}
//
//	cnt, err := res.RowsAffected()
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onDelete+": "+mysqllib.CantGetRowsAffected, objOp.sqlDelete, id)
//	}
//
//	return crud.Result{cnt}, objOp.setLinks(userIS, id, nil)
//}
//
//func (objOp *notesMySQL) QueryByPrefix(userIS auth.ID, rView auth.ID, genus *notes.Genus, prefix string) ([]notes.Asked, error) {
//	return nil, nil
//}
//
////const onClear = "on notesMySQL.Clean()"
////
////func (objOp *notesMySQL) Clean(selector selectors.Selector) error {
////	if !basis.CheckEnvs("ADMIN_MODE") {
////		return errors.NewCRUDOperator(onClear + ": insufficient environment")
////	}
////
////	// TODO: clear links
////
////	var err error
////
////	if selector == nil {
////		_, err = objOp.dbh.Exec("truncate `" + objOp.table + "`")
////	} else {
////		condition, values, err := selectors.Mysql("", selector)
////		if err != nil {
////			return errors.Wrap(err, onClear)
////		}
////		if strings.TrimSpace(condition) != "" {
////			condition = "where " + condition
////		}
////		_, err = objOp.dbh.Exec("delete from `"+objOp.table+"` "+condition, values...)
////	}
////
////	return err
////}
//
//// Close is native CRUD method.
//func (objOp *notesMySQL) Close() error {
//	return errors.Wrap(objOp.dbh.Close(), "on notesMySQL.dbh.Close()")
//}
//
//// Search -----------------------------------------------------------------------------------------
//
//var rePhrase = regexp.MustCompile(`^\s*".*"\s*$`)
//var reDelimiter = regexp.MustCompile(`[\.,\s\t;:\-\+\!\?\(\)\{\}\[\]\/'"\*]+`)
//
//func (objOp *notesMySQL) ReadListByWords(userIS auth.ID, options *content.ListOptions, searched string) (objects []notes.Item, allCnt uint64, err error) {
//	if !rePhrase.MatchString(searched) {
//		words := reDelimiter.Split(searched, -1)
//		searched = ""
//		for _, w := range words {
//			if len(w) > 2 {
//				searched += " +" + w
//			}
//		}
//	}
//
//	selectorSearched := selectors.Match("name,content,tags", searched, "IN BOOLEAN MODE")
//	if options == nil {
//		options = &content.ListOptions{Selector: selectorSearched}
//	} else if options.Selector == nil {
//		options.Selector = selectorSearched
//	} else {
//		options.Selector = selectors.And(options.Selector, selectorSearched)
//	}
//	return objOp.ReadList(userIS, options)
//}
//
//const onReadListByTag = "on notesMySQL.ReadListByTag"
//
//func (objOp *notesMySQL) ReadListByTag(userIS auth.ID, options *content.ListOptions, tag string) (linkedObjs []notes.Item, allCnt uint64, parentIDs []string, err error) {
//	var linked []links.Linked
//	if objOp.linksOp != nil {
//		linked, err = objOp.linksOp.QueryByTag(userIS, tag)
//		if err != nil {
//			return nil, 0, nil, errors.Wrap(err, onReadListByTag)
//		}
//	}
//
//	var linkedIDs []string
//	for _, l := range linked {
//		duplicatedID := false
//		for _, id := range linkedIDs {
//			if id == l.LinkedID {
//				duplicatedID = true
//				break
//			}
//		}
//		if !duplicatedID {
//			linkedIDs = append(linkedIDs, l.LinkedID)
//		}
//	}
//
//	if len(linkedIDs) <= 0 {
//		return nil, 0, nil, nil
//	}
//
//	// sort.Strings(linkedIDs)
//
//	selectorTagged := selectors.FieldStr("id", linkedIDs...)
//	if options == nil {
//		options = &content.ListOptions{Selector: selectorTagged}
//	} else if options.Selector == nil {
//		options.Selector = selectorTagged
//	} else {
//		options.Selector = selectors.And(options.Selector, selectorTagged)
//	}
//
//	for _, l := range linked {
//		id, _ := strconv.ParseUint(l.ObjectID, 10, 64)
//		if id > 0 {
//			duplicatedID := false
//			idStr := strings.TrimSpace(l.ObjectID)
//			for _, parentID := range parentIDs {
//				if idStr == parentID {
//					duplicatedID = true
//					continue
//				}
//			}
//			if !duplicatedID {
//				parentIDs = append(parentIDs, idStr)
//			}
//		}
//	}
//
//	linkedObjs, allCnt, err = objOp.ReadList(userIS, options)
//
//	return linkedObjs, allCnt, parentIDs, err
//}
//
//// links ------------------------------------------------------------------------------------
//
//// setLinks corrects object links without object itself
//
//func (objOp *notesMySQL) setLinks(userIS auth.ID, idStr string, linksListNew []links.Item) error {
//	var err error
//	var linkedInfo []links.LinkedInfo
//	if objOp.linksOp != nil {
//		linkedInfo, err = objOp.linksOp.SetLinks(userIS, "", idStr, linksListNew)
//		if err != nil {
//			return err
//		}
//	}
//
//	var errs basis.Errors
//	for _, l := range linkedInfo {
//		_, err := objOp.stmtUpdateCountLinked.Exec(l.CountLinked, l.ObjectID)
//		if err != nil {
//			errs = append(errs, err)
//		}
//	}
//
//	return errs.Err()
//}
//
//// UpdateLinks corrects object links within and without object itself
//
//const onUpdateLinks = "on notesMySQL.UpdateLinks"
//
//func (objOp *notesMySQL) UpdateLinks(userIS auth.ID, idStr string, linksListNew []links.Item, linkType string) error {
//	// TODO: lock object record for update (use history!!!)
//
//	o, err := objOp.Read(userIS, idStr)
//	if err != nil {
//		return errors.Wrap(err, onUpdateLinks)
//	}
//
//	linksList := notes.PrepareLinks(userIS, objOp.grpsOp, o.ROwner, o.Links, linksListNew, objOp.jointLinks, linkType)
//	linksListCopy := linksList
//	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
//	var jsonLinks []byte
//	if len(linksList) > 0 {
//		jsonLinks, err = json.Marshal(linksList)
//		if err != nil {
//			return errors.Wrapf(err, onUpdateLinks+": can't marshal .Links(%#v)", linksList)
//		}
//	}
//
//	values := []interface{}{jsonLinks, o.ID}
//	_, err = objOp.stmtUpdateLinks.Exec(values...)
//	if err != nil {
//		return errors.Wrapf(err, onUpdateLinks+": "+basis.CantExecQuery, objOp.sqlUpdateLinks, values)
//	}
//
//	return objOp.setLinks(userIS, o.ID, linksListCopy)
//}
//
////// export / import -----------------------------------------------------------------------------
////
////func (objOp *notesMySQL) GlobalIS(is string) (uint64, []notes.Item, error) {
////	var id uint64
////	var buf []byte
////	var links []notes.Item
////	err := objOp.stmtIsGlobal.QueryRow(is).Scan(&id, &buf)
////	if err != nil {
////		if err == sql.ErrNoRows {
////			return 0, nil, nil
////		}
////		return 0, nil, err
////	}
////	err = json.Unmarshal(buf, &links)
////	if err != nil {
////		log.Println("can't unmarshal links: ", string(buf), "for object.id:", id)
////	}
////	return id, links, nil
////}
//
//// fixtures -------------------------------------------------------------------------------------
//
////const ondeleteAll = "on ObjectsMysql.deleteAll"
////
////func (objOp *notesMySQL) deleteAll(userIS basis.UserIS, selector selectors.Selector) error {
////
////	condition, values, err := selectors.Mysql(userIS, selector)
////	if err != nil {
////		return nil.Wrapf(err, ": bad selector (%#v)", selector)
////	}
////
////	if strings.TrimSpace(condition) != "" {
////		condition = " where " + condition
////	}
////	sqlQuery := objOp.sqlDeleteAll + condition
////	_, err = objOp.dbh.Exec(sqlQuery, values...)
////	if err != nil {
////		return nil.Wrapf(err, ondeleteAll+": "+basis.CantExecQuery, sqlQuery, values)
////	}
////
////	return nil
////}
//
////const onloadFixture = "on notesMySQL.loadFixture"
////
////func (objOp *notesMySQL) loadFixture(userIS auth.ID, selector selectors.Selector, fixture fixturer.Fixture) error {
////	var numDeleted, numLoaded int
////
////	options := &content.ListOptions{Selector: selector}
////	objs, _, err := objOp.ReadList(userIS, options)
////	if err != nil {
////		log.Println(onloadFixture+"on read all data rows with fixture.Selector: ", err)
////	}
////	for _, o := range objs {
////		_, err = objOp.DeleteList(userIS, o.TargetID)
////		if err != nil {
////			log.Printf(onloadFixture+": on objOp.DeleteList(%v, %s): %s", userIS, o.TargetID, err)
////		} else {
////			numDeleted++
////		}
////	}
////
////	for _, row := range fixture.Data {
////		o, err := fixturer.Row(userIS, fixture.Fields, row, objOp.generaOp)
////		if err != nil {
////			log.Printf("data row (%#v) isn't loaded: %s", row, err)
////			continue
////		}
////		if o == nil {
////			log.Printf("empty data row (%#v) isn't loaded", row)
////			continue
////		}
////
////		_, err = objOp.Create(userIS, *o)
////		if err != nil {
////			log.Printf("data row (%#v) isn't loaded: %s", row, err)
////			continue
////		}
////
////		numLoaded++
////	}
////
////	log.Printf(onloadFixture+": deleted %d (of %d total) old data rows ", numDeleted, len(objs))
////	log.Printf(onloadFixture+": loaded %d (of %d total) fixture data rows ", numLoaded, len(fixture.Data))
////
////	return nil
////}
