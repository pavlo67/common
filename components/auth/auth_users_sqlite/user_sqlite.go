package auth_users_sqlite

import (
	"strings"

	"github.com/GehirnInc/crypt"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/libs/encrlib"

	"github.com/pavlo67/workshop/components/auth"
)

var _ auth.Operator = &isentityLoginStub{}

type isentityLoginStub struct {
	users []UserSQLite
	salt  string
}

//const login = "йа"
//const password = "мій пароль"

func New(users []UserSQLite, salt string) (*isentityLoginStub, error) {
	return &isentityLoginStub{
		users: users,
		salt:  salt,
	}, nil
}

//func (u *isentityLoginStub) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsPassword}, nil
//}

func (u *isentityLoginStub) SetCreds(user auth.User, toSet auth.Creds) (*auth.Creds, error) {
	return nil, common.ErrNotImplemented
}

func (u *isentityLoginStub) Authorize(toAuth auth.Creds) (*auth.User, error) {
	login := toAuth.Values[auth.CredsNickname]
	login = toAuth.Values[auth.CredsEmail]

	password := toAuth.Values[auth.CredsPassword]
	cryptype := toAuth.Cryptype

	for _, user := range u.users {
		if user.Login == login {
			switch cryptype {
			case encrlib.SHA256:
				crypt := crypt.SHA256.New()
				passwordHash, _ := crypt.Generate([]byte(strings.TrimSpace(password)), []byte(u.salt))
				if password == passwordHash {
					return &auth.User{ID: user.ID, Nick: user.Login}, nil
				}
			default:
				if password == user.Password {
					return &auth.User{ID: user.ID, Nick: user.Login}, nil
				}
			}
		}
	}

	return nil, auth.ErrBadPassword
}
