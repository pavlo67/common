package auth_stub

import (
	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/rbac"
)

var _ auth.Operator = &authstub{}

type User struct {

	// TODO! ba careful: we set authID == nickname for this stub case

	passhash string
	roles    rbac.Roles
}

type authstub struct {
	crypter crypt.Crypter
	users   map[auth.ID]User
}

const onNew = "on authstub.New()"

func New(defaultUser config.Access) (auth.Operator, error) {

	crypter := crypt.SHA256.New()

	defaultUserPasshash, err := crypter.Generate([]byte(defaultUser.Pass), nil)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	authOp := authstub{
		crypter: crypter,
		users: map[auth.ID]User{
			auth.ID(defaultUser.User): {
				passhash: defaultUserPasshash,
				roles:    rbac.Roles{rbac.RoleAdmin},
			},
		},
	}

	return &authOp, nil
}

const onSetCreds = "on authstub.SetCreds()"

func (authOp *authstub) SetCreds(authID auth.ID, toSet auth.Creds) (*auth.Creds, error) {
	nickname := toSet[auth.CredsNickname]

	if authID == "" {
		authID = auth.ID(nickname)
	}

	passhash, err := authOp.crypter.Generate([]byte(toSet[auth.CredsPassword]), nil)
	if err != nil {
		return nil, errors.Wrapf(err, onSetCreds+": can't hash password: %s", err)
	}

	var user *User
	if userOld, ok := authOp.users[authID]; ok {
		user = &userOld
	} else {
		user = &User{}
	}

	user.passhash = passhash
	if role, ok := toSet[auth.CredsRole]; ok {
		// multiple roles aren't supported here
		user.roles = rbac.Roles{rbac.Role(role)}
	}

	authOp.users[authID] = *user

	return &auth.Creds{auth.CredsNickname: nickname}, nil
}

const onAuthenticate = "on authstub.Authenticate()"

func (authOp *authstub) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
	nickname := toAuth[auth.CredsNickname]
	authID := auth.ID(nickname)

	if user, ok := authOp.users[authID]; ok {
		if err := authOp.crypter.Verify(user.passhash, []byte(toAuth[auth.CredsPassword])); err == nil {
			return &auth.Identity{
				ID:       authID,
				Nickname: nickname,
				Roles:    user.roles,
			}, nil
		}
	}

	return nil, auth.ErrNotAuthenticated
}
