package auth_stub

import (
	"strings"

	"github.com/GehirnInc/crypt"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"

	"github.com/pavlo67/workshop/libraries/encrlib"
)

var _ auth.Operator = &isentityLoginStub{}

type isentityLoginStub struct {
	users []UserStub
	salt  string
}

//const login = "йа"
//const password = "мій пароль"

func New(users []UserStub, salt string) (*isentityLoginStub, error) {
	return &isentityLoginStub{
		users: users,
		salt:  salt,
	}, nil
}

//func (u *isentityLoginStub) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsPassword}, nil
//}

func (_ *isentityLoginStub) GetSessionKeys() (common.Map, error) {
	return nil, nil
}

func (u *isentityLoginStub) SetCreds(user auth.User, toSet auth.Creds) (*auth.Creds, error) {
	return nil, common.ErrNotImplemented
}

func (u *isentityLoginStub) Authorize(toAuth auth.Creds) (*auth.User, error) {
	login := toAuth.Values[auth.CredsLogin]

	nickname := toAuth.Values[auth.CredsNickname]
	if nickname != "" {
		login = nickname
	}

	email := toAuth.Values[auth.CredsEmail]
	if email != "" {
		login = email
	}

	password := toAuth.Values[auth.CredsPassword]
	cryptype := toAuth.Cryptype

	for _, user := range u.users {
		// l.Infof("%#v: %s, %s", user, login, password)

		if user.Login == login {
			switch cryptype {
			case encrlib.SHA256:
				crypt := crypt.SHA256.New()
				passwordHash, _ := crypt.Generate([]byte(strings.TrimSpace(password)), []byte(u.salt))
				if password == passwordHash {
					return &auth.User{Key: user.ID, Nickname: user.Login}, nil
				}
			default:
				if password == user.Password {
					return &auth.User{Key: user.ID, Nickname: user.Login}, nil
				}
			}
		}
	}

	return nil, auth.ErrPassword
}
