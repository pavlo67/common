package groupsstub

import (
	"time"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/confidenter/rights"
	"github.com/pavlo67/punctum/starter/joiner"
)

var _ groups.Operator = &ControllerStubStarter{}

type ControllerStubStarter struct {
	userData     map[auth.ID][]group
	rViewDefault auth.ID
}

type group struct {
	id string
	is auth.ID
}

func New(groupIDs map[auth.ID][]string, rViewDefault auth.ID) (*ControllerStubStarter, error) {
	domain := joiner.SystemDomain()

	d := ControllerStubStarter{
		map[auth.ID][]group{},
		rViewDefault,
	}
	for is, userGroupIDs := range groupIDs {
		for _, groupID := range userGroupIDs {
			ident := auth.IDentity{domain, groups.IdentityPathDefault, groupID}
			d.userData[is] = append(
				d.userData[is],
				group{
					groupID,
					ident.String(),
				},
			)
		}
	}

	return &d, nil
}

func (ctrlOp *ControllerStubStarter) Create(userIS auth.ID, data groups.Item) (string, error) {
	return "", nil
}

func (ctrlOp *ControllerStubStarter) Read(userIS auth.ID, id string) (*groups.Item, error) {
	return nil, nil
}

func (ctrlOp *ControllerStubStarter) ReadList(userIS auth.ID, options *content.ListOptions, selector selectors.Selector) ([]groups.Item, uint64, error) {
	return nil, 0, nil
}

func (ctrlOp *ControllerStubStarter) Update(userIS auth.ID, data groups.Item) (result crud.Result, err error) {
	return crud.Result{}, nil
}

func (ctrlOp *ControllerStubStarter) Delete(userIS auth.ID, id string) (result crud.Result, err error) {
	return crud.Result{}, nil
}

func (ctrlOp *ControllerStubStarter) SetRights(userIS, groupIS, memberIS auth.ID, rightsToTime map[rights.Right]*time.Time) error {
	return nil
}

func (ctrlOp *ControllerStubStarter) BelongsTo(userIS, dataIS auth.ID) (bool, error) {
	if dataIS == basis.Anyone {
		return true, nil
	}

	for _, g := range ctrlOp.userData[userIS] {
		if g.is == dataIS {
			return true, nil
		}
	}
	return false, nil

}

func (ctrlOp *ControllerStubStarter) AllForUser(userIS auth.ID) ([]auth.IDentityNamed, error) {
	return nil, nil
}

func (ctrlOp *ControllerStubStarter) Close() error {
	return nil
}
