package groups

//type Right string
//type Managers map[Right]common.ID
//
//const (
//	Owner   Right = "owner"
//	Manager Right = "manager"
//	Member  Right = "member"
//
//	View   Right = "view"   // view the object
//	Change Right = "change" // change the object
//	Use    Right = "use"    // use the object
//
//	Add    Right = "add"    // add to list
//	Remove Right = "remove" // remove from list
//)
//
//const AllowedForAll common.ID = "*"
//const AllowedForAllAuthorized common.ID = "!"

//func QueryAccessible(grpsOp Operator, dbh *sql.DB, is auth.ID, sqlQuery, sqlCondition, sqlPostCondition string, valuesCondition []interface{}) (string, *sql.Rows, error) {
//	var err error
//	var groups []auth.IDentityNamed
//
//	if grpsOp != nil {
//		if groups, err = grpsOp.AllForUser(is); err != nil {
//			return "", nil, err
//		}
//	}
//
//	available := "r_view in (?,?"
//	values := append(valuesCondition, string(is), string(basis.Anyone))
//	for _, g := range groups {
//		available += ", ?"
//		values = append(values, string(g.String()))
//	}
//	available += ")"
//	if sqlCondition != "" {
//		available = " and (" + available + ")"
//	}
//	sqlQuery += " where " + sqlCondition + available + " " + sqlPostCondition
//
//	var stmt *sql.Stmt
//	if stmt, err = dbh.Prepare(sqlQuery); err != nil {
//		return sqlQuery, nil, errors.Wrapf(err, basis.CantPrepareQuery, sqlQuery)
//	}
//	defer stmt.Close()
//
//	var rows *sql.Rows
//	if rows, err = stmt.Query(values...); err != nil {
//		if err == sql.ErrNoRows {
//			return sqlQuery, nil, basis.ErrNotFound
//		}
//		return sqlQuery, nil, errors.Wrapf(err, basis.CantExecQuery, sqlQuery, values)
//	}
//
//	return sqlQuery, rows, nil
//}
//
//func OneOf(is auth.ID, grpsOp Operator, controlledISs ...auth.ID) bool {
//	for _, controlledIS := range controlledISs {
//		if controlledIS == "" {
//			continue
//		}
//
//		if is == controlledIS ||
//			controlledIS == basis.Anyone ||
//			(controlledIS == basis.AnyoneRegistered && is != "") {
//			return true
//		}
//
//		if grpsOp != nil {
//			ok, err := grpsOp.BelongsTo(is, controlledIS)
//			if err != nil {
//				log.Print("ERROR", err)
//				return false
//			}
//			if ok {
//				return true
//			}
//		}
//	}
//
//	return false
//}
//
//func OneOfErr(is auth.ID, grpsOp Operator, controlledISs ...auth.ID) error {
//	for _, controlledIS := range controlledISs {
//		if controlledIS == "" {
//			continue
//		}
//
//		if is == controlledIS ||
//			controlledIS == basis.Anyone ||
//			(controlledIS == basis.AnyoneRegistered && is != "") {
//			return nil
//		}
//
//		if grpsOp != nil {
//			ok, err := grpsOp.BelongsTo(is, controlledIS)
//			if err != nil {
//				return err
//			}
//			if ok {
//				return nil
//			}
//		}
//	}
//
//	return rights.ErrNoRights
//}
//
//func SetRights(is auth.ID, grpsOp Operator, rView, rOwner auth.ID, managers rights.Managers) (auth.ID, auth.ID, rights.Managers, error) {
//	if len(managers) > 0 {
//		managers[rights.View] = rView
//		managers[rights.Owner] = rOwner
//		if managers[rights.View] == "" {
//			managers[rights.View] = is
//		}
//		if managers[rights.Owner] == "" {
//			managers[rights.Owner] = is
//		}
//	} else {
//		managers = rights.Managers{rights.View: rView, rights.Owner: rOwner}
//	}
//
//	for r, is := range managers {
//		if is == "" {
//			continue
//		}
//
//		// check if userIS belongs to the selected "group" to set the rights for it
//		if err := OneOfErr(is, grpsOp, managers[r]); err != nil {
//			return "", "", nil, err
//		}
//	}
//
//	rView = managers[rights.View]
//	rOwner = managers[rights.Owner]
//	delete(managers, rights.View)
//	delete(managers, rights.Owner)
//
//	return rView, rOwner, managers, nil
//}

//// CanOperate is universal method for all Census interfaces
//func CanOperate(confidenter *auth.IDentity, gr Operator, groupIS basis.UserIS, domain string, rr ...rights.Right) (rights.Managers, error) {
//
//	data, err := gr.Read(confidenter, groupIS)
//	if err != nil {
//		return nil, nil.Wrapf(err, "can't get object description to operate (confidenter, object): %v, %v", confidenter, groupIS)
//	}
//
//	for _, r := range rr {
//		err = groups.IsManager(confidenter, gr, data.Managers, r)
//		if err == nil {
//			return data.Managers, nil
//		}
//	}
//	return nil, nil.Wrapf(err, "can't confirm rights to operate group (confidenter, group, rights): %v, %v, %v", confidenter, groupIS, rr)
//}

//func GetAllAccessible(confidenter auth.IDentity, ctrl groups.Operator) (string, []string, error){
//
//	var valuesAll = []string{}
//	available := ""
//	userGroupsIS, err := ctrl.AllAccessible(confidenter)
//	if err == nil && len(userGroupsIS) > 0 {
//		available = "r_view in ("
//		for i, r := range userGroupsIS {
//			if i > 0 {
//				available += ", "
//			}
//			available += "?"
//			valuesAll = append(valuesAll, string(r))
//		}
//		available += ")"
//	}
//	return available, valuesAll, err
//}

//func HasRights(user *auth.User, identOpsMap map[auth.CredsType][]Operator, allowedIDs []common.ID) (bool, error) {
//	if allowedIDs == nil {
//		return true, nil
//	} else if len(allowedIDs) < 1 {
//		return false, nil
//	}
//
//	for _, allowedID := range allowedIDs {
//		if allowedID == AllowedForAll {
//			return true, nil
//		}
//	}
//
//	if user == nil {
//		return false, nil
//	}
//
//	for _, allowedID := range allowedIDs {
//		if allowedID == AllowedForAllAuthorized || allowedID == user.ID {
//			return true, nil
//		}
//	}
//
//	// TODO: check if user is in some of AllowedIDs... groups
//	// for _, identOp := range identOpsMap[CredsAllowedID] {
//	//	if identOp == nil {
//	//		continue
//	//	}
//	//	var errs basis.Errors
//	//	for _, allowedID := range allowedIDs {
//	//		user, _, err := identOp.Authorize(Creds{Type: CredsID, Value: string(user.ID)}, Creds{Type: CredsAllowedID, Value: string(allowedID)})
//	//		if err != nil {
//	//			errs = append(errs, err)
//	//		}
//	//		if user != nil {
//	//			return true, errs.Err()
//	//		}
//	//	}
//	//
//	//	return false, errs.Err()
//	// }
//
//	return false, nil
//}
