package rights

import (
	"errors"

	"github.com/pavlo67/workshop/common"
)

var ErrNoRights = errors.New("користувачу бракує прав :-(")

type Right string
type Managers map[Right]common.ID

const (
	Owner   Right = "owner"
	Manager Right = "manager"
	Member  Right = "member"

	View   Right = "view"   // view the object
	Change Right = "change" // change the object
	Use    Right = "use"    // use the object

	Add    Right = "add"    // add to list
	Remove Right = "remove" // remove from list
)

const AllowedForAll common.ID = "*"
const AllowedForAllAuthorized common.ID = "!"

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
