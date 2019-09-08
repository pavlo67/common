package auth_users

import (
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/connector"
	"github.com/pavlo67/partes/connector/sender"
	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/encrlib"
	"github.com/pavlo67/punctum/basis/strlib"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/users"
	"github.com/pavlo67/punctum/content"
	"github.com/pavlo67/punctum/starter/joiner"
)

const createMessageSubject = "Підтвердження реєстрації / Registration confirmation: "
const createMessageContent = "Щоб завершити реєстрацію, будь ласка, пройдіть за ланкою / To complete registration please follow the link: "

const sendCodeMessageSubject = "Підтвердження зміни паролю / Password updating confirmation: "
const sendCodeMessageContent = "Щоб змінити пароль, будь ласка, пройдіть за ланкою / To change your password please follow the link: "

// -----------------------------------------------------------------------------

var _ auth.Operator = &authusers{}

type authusers struct {
	domain      string
	usersOp     users.Operator
	senderOp    sender.Operator
	salt        string
	useMessages bool
	callbacks   map[auth.Callback]string
}

const onNew = "on authusers.New"

func New(usersOp users.Operator, senderOp sender.Operator, salt string, useMessages bool) (*authusers, error) {
	if usersOp == nil {
		return nil, errors.New(onNew + ": no users.Operator for authusers.New()")
	}

	if senderOp == nil {
		return nil, errors.New(onNew + ": no sender.Operator for authusers.New()")
	}

	if salt == "" {
		return nil, errors.New(onNew + ": no salt for authusers.New()")
	}

	return &authusers{
		domain:      joiner.SystemDomain(),
		usersOp:     usersOp,
		senderOp:    senderOp,
		salt:        salt,
		useMessages: useMessages,
		callbacks:   map[auth.Callback]string{},
	}, nil
}

const onCreate = "on authusers.Create"

func (authOp *authusers) Create(creds ...auth.Creds) ([]auth.Message, error) {
	token := strlib.RandomString(users.ControlTokenLength)

	userToCreate := auth.NewUserToCreate(creds)

	hash, err := encrlib.GetEncodedPassword(userToCreate.Password, authOp.salt, 0)

	user := users.User{
		Nickname:     userToCreate.Nickname,
		Hash:         hash,
		Email:        userToCreate.Email,
		ControlToken: token,
	}

	_, err = authOp.usersOp.Create("", user, true)
	if err != nil {
		return nil, errors.Wrap(err, onCreate+": can't create user")
	}

	return nil, authOp.sendCode(auth.Confirm, user.Email, token)
}

const onAuthenticateWithCreds = "on authusers.AuthenticateWithCreds"
const onAuthenticateWithCode = "on authusers.RegisterWithCode"
const onSendCode = "on authusers.SendCode"

func (authOp *authusers) Use(toUse, toAuth auth.Creds, toSet ...auth.Creds) (*auth.User, []auth.Message, error) {
	if toAuth.Type == auth.CredsSentCode {
		// TODO: use "rector's" userIS to prevent second read
		user, err := authOp.readByField("control_token", toAuth.FirstValue())
		if err == nil {
			err = authOp.usersOp.SetVerified(auth.IS(user.ID), user.ID)
		}
		if err != nil {
			return nil, nil, errors.Wrap(err, onAuthenticateWithCode)
		}

		return &auth.User{
			ID:       user.ID,
			Nickname: "",
			Accesses: nil,
		}, nil, nil
	}

	var user *users.User
	var err error

	if toUse.Type == auth.CredsNickname {
		user, err = authOp.readByField("nickname", toUse.FirstValue())
		if err != nil {
			return nil, nil, errors.Wrap(err, onAuthenticateWithCreds)
		}
	} else if toUse.Type == auth.CredsEmail {
		user, err = authOp.readByField("email", toUse.FirstValue())
		if err != nil {
			return nil, nil, errors.Wrap(err, onAuthenticateWithCreds)
		}
	} else {
		return nil, nil, basis.ErrNotImplemented
	}

	if user == nil {
		return nil, nil, basis.ErrNull
	}

	if toAuth.Type == auth.CredsPassword {
		if !user.Verified {
			return nil, nil, auth.ErrUserNotVerified
		} else if !auth.CheckPassword(authOp.salt, toAuth.FirstValue(), user.Cryptype, user.Passhash) {
			return nil, nil, auth.ErrBadPasshash
		}

		return &auth.User{
			ID:       user.ID,
			Nickname: "",
			Accesses: nil,
		}, nil, nil
	}

	for _, ts := range toSet {
		if ts.Type == auth.CredsSentCode {
			user.ControlToken = strlib.RandomString(users.ControlTokenLength)
			_, err = authOp.usersOp.Update(auth.IS(user.ID), *user)
			if err != nil {
				return nil, nil, errors.Wrap(err, onSendCode+": can't set control token")
			}

			if !user.Verified {
				return nil, nil, authOp.sendCode(auth.Confirm, user.Email, user.ControlToken)
			}

			return nil, nil, authOp.sendCode(auth.SendCode, user.Email, user.ControlToken)
		}
	}

	return nil, nil, basis.ErrNotImplemented
}

func (authOp *authusers) AddCallback(key auth.Callback, url string) {
	authOp.callbacks[key] = url
}

func (authOp *authusers) sendCode(cbKey auth.Callback, email, token string) error {
	cb, ok := authOp.callbacks[cbKey]
	if !ok {
		return errors.Errorf("wrong confirmation callback type: '%s'", cbKey)
	}

	message := connector.Message{
		To:   email,
		Body: token,
	}

	switch cbKey {
	case auth.Confirm:
		message.Subject = createMessageSubject
		if authOp.useMessages {
			message.Body = createMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
		}
	case auth.SendCode:
		message.Subject = sendCodeMessageSubject
		if authOp.useMessages {
			message.Body = sendCodeMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
		}
	}

	return authOp.senderOp.Send(message)
}

func (authOp *authusers) readByField(field, value string) (*users.User, error) {
	usrs, _, err := authOp.usersOp.ReadList(
		"",
		&content.ListOptions{Selector: selectors.FieldEqual(field, strings.TrimSpace(value))},
	)
	if err != nil || len(usrs) != 1 {
		return nil, errors.Errorf("error finding user with %s = '%s' (found %d): %s", field, value, len(usrs), err)
	}

	return &usrs[0], nil
}

//var reEmailToLogin1 = regexp.MustCompile(`@.*`)
//var reEmailToLogin2 = regexp.MustCompile(`(\.|-)`)

//type respFB struct {
//	TargetID    string `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
//
//const onQueryPartnerUser = "on authusers.QueryPartnerUser"
//
//func (authOp *authusers) QueryPartnerUser(partnerKey, partnerToken string) (partnerUser *auth.User, err error) {
//
//	// TODO: implement it as an interface call to "QueryPartnerUser" implementation according to partnerURL
//
//	partnerURL := "https://graph.facebook.com/me?fields=id,name,email&access_token="
//	req, err := http.NewRequest("GET", partnerURL+partnerToken, nil)
//
//	var resp *http.Response
//	if err == nil {
//		var client = &http.Client{}
//		resp, err = client.Do(req)
//		defer resp.Body.Close()
//	}
//	var data respFB
//	decoder := json.NewDecoder(resp.Body)
//	err = decoder.Decode(&data)
//	if err != nil {
//		return nil, errors.Wrapf(err, onQueryPartnerUser+": can't decode JSON from FB")
//	}
//
//	// log.Println("UserIS FB data:", data)
//
//	// TODO: use respFB.Label to control user if his email is changed
//
//	return &auth.User{IdentityNamed: auth.IDentityNamed{Nick: data.Name}, Contacts: data.Email}, nil
//}
//
//const onAuthenticateWithPartnerUser = "on authusers.AuthenticateWithPartnerUser"
//
//func (authOp *authusers) AuthenticateWithPartnerUser(partnerUser *auth.User) (*auth.User, error) {
//	user, err := authOp.readByLogin(partnerUser.Contacts)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthenticateWithCreds)
//	}
//
//	if user != nil {
//		user.Contacts = partnerUser.Contacts
//		return user, nil
//	}
//
//	partnerUser.TargetID, err = authOp.usersOp.Create(nil, partnerUser, true)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthenticateWithPartnerUser+": can't create user")
//	}
//
//	return partnerUser, nil
//}
//
//func (authOp *authusers) Clean(selector selectors.Selector) error {
//	return authOp.usersOp.Clean(selector)
//}
//
//func (authOp *authusers) ConsoleDo(action string, values map[string]string) error {
//
//	//// TODO:  refactor it
//	//
//	//var sqlQuery string
//	//var err error
//	//var stmt *sql.Stmt
//	//var res sql.Result
//	//var valuesAll []interface{}
//	//if action == "new_user" {
//	//	valuesAll = []interface{}{values["login"], values["password"], values["email"]}
//	//	sqlQuery = "insert into `" + u.userTable + "` (nickname, passhash, email, verified, contacts, history) values (?,?,?,1,'','')"
//	//} else if action == "new_password" {
//	//	valuesAll = []interface{}{values["password"], values["login"]}
//	//	sqlQuery = "update `" + u.userTable + "` set passhash=? where nickname=?"
//	//}
//	//
//	//if stmt, err = u.dbh.Init(sqlQuery); err != nil {
//	//	return nil.Wrapf(err, "can't prepare sql:%v", sqlQuery)
//	//}
//	//defer stmt.Close()
//	//
//	//if res, err = stmt.Exec(valuesAll...); err != nil {
//	//	return nil.Wrapf(err, "can't exec sql:%v, values=%v", sqlQuery, valuesAll)
//	//}
//	//iA, err := res.RowsAffected()
//	//
//	//if iA < 1 {
//	//	return nil.New("!!! RowsAffected() = ZERO !")
//	//}
//	//
//	return nil
//}

//
//const onUpdatePasswordWithCode = "on authusers.UpdatePasswordWithCode"
//
//func (authOp *authusers) UpdatePasswordWithCode(controlToken string, encodedNewPassword encrlib.Hash) error {
//
//	if controlToken == "" {
//		return errors.Wrap(confidenter.ErrTokenNotFound, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.Passhash == "" {
//		return errors.Wrap(auth.ErrBadPassword, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.Cryptype != encrlib.SHA256 {
//		return errors.Wrap(encrlib.ErrBadCryptype, onUpdatePasswordWithCode)
//	}
//
//	usrs, _, err := authOp.usersOp.ReadList(nil, &content.ListOptions{Selector: selectors.FieldEqual(users.ControlTokenField, controlToken)})
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error finding user")
//	}
//	if len(usrs) != 1 {
//		return fmt.Errorf(onUpdatePasswordWithCode+": can't find user (find %d)", len(usrs))
//	}
//
//	userIS := &basis.UserIS{joiner.SystemDomain(), "/user/", usrs[0].TargetID}
//
//	// TODO: use correct domain/path/, prevent duplicate reading
//	user, err := authOp.usersOp.Read(userIS, usrs[0].TargetID)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error reading user")
//	}
//
//	user.Cryptype = encodedNewPassword.Cryptype
//	user.Passhash = encodedNewPassword.Passhash
//	_, err = authOp.usersOp.Update(userIS, user)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": can't set new password")
//	}
//
//	return nil
//}
