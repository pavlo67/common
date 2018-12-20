package identity

import (
	"errors"
)

var ErrNoRights = errors.New("користувачу бракує прав :-(")

// Right ...
type Right string

// Managers ...
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
