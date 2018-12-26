package identity

import (
	"errors"

	"github.com/pavlo67/punctum/basis"
)

var ErrNoRights = errors.New("користувачу бракує прав :-(")

type Right string
type Managers map[Right]ID

const (
	Owner           Right = "owner"
	Member          Right = "member"
	Manager         Right = "manager"
	Delete          Right = "delete"          // delete object
	Create          Right = "create"          // create object
	View            Right = "view"            // view the object (except rights to it)
	ViewDefault     Right = "view_default"    // view the object by default
	ViewRights      Right = "view_rights"     // view who has any rights to the object
	Change          Right = "change"          // change the object (except rights to it)
	ChangeRights    Right = "change_rights"   // change any rights to the object
	AddTo           Right = "add_to"          // adding to the object (except rights to it)
	RemoveFrom      Right = "remove_from"     // remove from the object (except rights to it)
	SendTo          Right = "send_to"         // send message to the object
	RViewMembers    Right = "view_members"    //
	RAddMember      Right = "add_member"      //
	RRemoveMember   Right = "remove_member"   //
	RRestrictMember Right = "restrict_member" //
)

const AllowedForAll ID = "*"
const AllowedForAllAuthorized ID = "!"

func HasRights(user *User, identOpsMap map[CredsType][]Operator, allowedIDs []ID) (bool, error) {
	if allowedIDs == nil {
		return true, nil
	} else if len(allowedIDs) < 1 {
		return false, nil
	}

	if user == nil {
		for _, allowedID := range allowedIDs {
			if allowedID == AllowedForAll {
				return true, nil
			}
		}
		return false, nil
	}

	for _, allowedID := range allowedIDs {
		if allowedID == AllowedForAll || allowedID == AllowedForAllAuthorized || allowedID == user.ID {
			return true, nil
		}
	}

	if len(user.Accesses) > 0 {
		for _, access := range user.Accesses {
			if access.ID == user.ID {
				return true, nil
			}
		}
	}

	for _, identOp := range identOpsMap[CredsAllowedID] {
		if identOp == nil {
			continue
		}
		var errs basis.Errors
		for _, allowedID := range allowedIDs {
			user, _, err := identOp.Authorize(Creds{Type: CredsID, Value: string(user.ID)}, Creds{Type: CredsAllowedID, Value: string(allowedID)})
			if err != nil {
				errs = append(errs, err)
			}
			if user != nil {
				return true, errs.Err()
			}
		}

		return false, errs.Err()
	}

	return false, nil
}
