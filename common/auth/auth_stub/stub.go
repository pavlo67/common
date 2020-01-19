package auth_stub

import (
	"strconv"
	"strings"
	"time"

	"github.com/GehirnInc/crypt"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/libraries/encrlib"
)

var _ auth.Operator = &authStub{}

type UserStub struct {
	Key          identity.Key
	Nickname     string
	PasswordHash string
}

type authStub struct {
	users   []UserStub
	crypter crypt.Crypter
	salt    string
}

func New(users []UserStub, salt string) (*authStub, error) {
	return &authStub{
		users:   users,
		crypter: crypt.SHA256.New(),
		salt:    salt,
	}, nil
}

func (u *authStub) SetCreds(userKey identity.Key, creds auth.Creds, _ auth.CredsType) (identity.Key, *auth.Creds, error) {
	passwordToSet := strings.TrimSpace(creds.Values[auth.CredsPassword])
	l.Infof("password to set  : %s", passwordToSet)

	passwordHash, err := u.crypter.Generate([]byte(passwordToSet), []byte(u.salt))
	if err != nil {
		return "", nil, err
	}

	l.Infof("password hash set: %s", passwordHash)
	l.Infof("salt: %s", u.salt)
	l.Infof("verify: %s", u.crypter.Verify(passwordHash, []byte(passwordToSet)))

	userStub := UserStub{
		Nickname:     creds.Values[auth.CredsNickname],
		PasswordHash: passwordHash,
	}

	if userKey == "" {
		userKey = identity.Key(strconv.FormatInt(time.Now().UnixNano(), 10))
		userStub.Key = userKey

		u.users = append(u.users, userStub)

		return userKey, &creds, nil
	}

	for i, us := range u.users {
		if us.Key == userKey {
			u.users[i] = userStub
			return userKey, &creds, nil
		}
	}

	return "", nil, auth.ErrNoUser
}

func (u *authStub) Authorize(toAuth auth.Creds) (*auth.User, error) {

	login := toAuth.Values[auth.CredsLogin]
	if login == "" {
		email := toAuth.Values[auth.CredsEmail]
		if email != "" {
			login = email
		}
		nickname := toAuth.Values[auth.CredsNickname]
		if nickname != "" {
			login = nickname
		}
	}

	password := strings.TrimSpace(toAuth.Values[auth.CredsPassword])

	// l.Infof("password to check: %s", password)

	// var passwordHash string

	switch toAuth.Cryptype {
	//case encrlib.SHA256:
	//	passwordHash = strings.TrimSpace(password)
	case encrlib.NoCrypt:
		//passwordHash, _ = u.crypter.Generate([]byte(password), []byte(u.salt))
	default:
		return nil, auth.ErrEncryptionType
	}

	// l.Infof("password hash    : %s", passwordHash)

	for _, us := range u.users {
		// l.Infof("%#v: %s, %s", us, login, password)

		if us.Nickname == login {
			// l.Info("++")
			if u.crypter.Verify(us.PasswordHash, []byte(password)) == nil {
				return &auth.User{Key: us.Key, Creds: auth.Creds{
					Cryptype: encrlib.NoCrypt,
					Values: auth.Values{
						auth.CredsNickname: us.Nickname,
					},
				}}, nil
			}
			return nil, auth.ErrPassword
		}
	}

	return nil, errors.Wrap(auth.ErrPassword, "no user")
}
