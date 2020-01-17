package auth_stub

import (
	"strings"

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

func (_ *authStub) InitAuthSession(_ auth.Creds) (*auth.Creds, error) {
	return nil, nil
}

func (u *authStub) SetCreds(user *auth.User, toSet auth.Creds) (*auth.User, *auth.Creds, error) {
	if user == nil {
		return nil, nil, auth.ErrNoUser
	}

	userStub := UserStub{
		Key:      user.Key,
		Nickname: user.Nickname,
	}

	passwordToSet := strings.TrimSpace(toSet.Values[auth.CredsPassword])
	l.Infof("password to set  : %s", passwordToSet)

	var err error
	userStub.PasswordHash, err = u.crypter.Generate([]byte(passwordToSet), []byte(u.salt))
	if err != nil {
		return nil, nil, err
	}

	l.Infof("password hash set: %s", userStub.PasswordHash)
	l.Infof("salt: %s", u.salt)

	//passwordHashAgain, _ := u.crypter.Generate([]byte(passwordToSet), []byte(u.salt))
	//l.Infof("password hash ???: %s", passwordHashAgain)
	//l.Infof("salt: %s", u.salt)

	l.Infof("verify: %s", u.crypter.Verify(userStub.PasswordHash, []byte(passwordToSet)))

	if toSet.Values[auth.CredsNickname] != "" {
		userStub.Nickname = toSet.Values[auth.CredsNickname]
		user.Nickname = toSet.Values[auth.CredsNickname]
	}

	for i, us := range u.users {
		if us.Key == user.Key {
			u.users[i] = userStub
			return user, nil, nil
		}
	}

	u.users = append(u.users, userStub)

	return user, nil, nil
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
				return &auth.User{Key: us.Key, Nickname: us.Nickname}, nil
			}
			return nil, auth.ErrPassword
		}
	}

	return nil, errors.Wrap(auth.ErrPassword, "no user")
}
