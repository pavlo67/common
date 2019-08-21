package auth_stub

import (
	"strings"

	"github.com/GehirnInc/crypt"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/encrlib"
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

func (u *isentityLoginStub) Accepts() ([]auth.CredsType, error) {
	return []auth.CredsType{auth.CredsPassword}, nil
}

func (u *isentityLoginStub) SetCreds(userID *common.ID, toSet []auth.Creds) (*auth.User, []auth.Creds, error) {
	return nil, nil, common.ErrNotImplemented
}

func (u *isentityLoginStub) Authorize(toAuth []auth.Creds) (*auth.User, []auth.Creds, error) {
	var login, password string
	var cryptype encrlib.Cryptype

	for _, creds := range toAuth {
		switch creds.Type {
		case auth.CredsNickname, auth.CredsEmail:
			login = creds.Value
		case auth.CredsPassword:
			password = creds.Value
			cryptype = creds.Cryptype
		}
	}

	for _, user := range u.users {
		if user.Login == login {
			switch cryptype {
			case encrlib.SHA256:
				crypt := crypt.SHA256.New()
				passwordHash, _ := crypt.Generate([]byte(strings.TrimSpace(password)), []byte(u.salt))
				if password == passwordHash {
					return &auth.User{ID: user.ID, Nick: user.Login}, nil, nil
				}
			default:
				if password == user.Password {
					return &auth.User{ID: user.ID, Nick: user.Login}, nil, nil
				}
			}
		}
	}

	return nil, nil, auth.ErrBadPassword
}
