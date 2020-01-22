package auth_users

import (
	"regexp"

	"github.com/GehirnInc/crypt"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/users"
)

var _ auth.Operator = &authPassUsers{}

type authPassUsers struct {
	usersOp             users.Operator
	maxUsersToAuthCheck int

	crypter crypt.Crypter
	salt    string
}

const onNew = "on authPassUsers.New"

func New(usersOp users.Operator, maxUsersToAuthCheck int, salt string) (auth.Operator, error) {
	if usersOp == nil {
		return nil, errors.New(onNew + ": no users.Operator for authPassUsers.New()")
	}
	if maxUsersToAuthCheck < 1 {
		maxUsersToAuthCheck = 1
	}

	return &authPassUsers{
		usersOp:             usersOp,
		maxUsersToAuthCheck: maxUsersToAuthCheck,

		crypter: crypt.SHA256.New(),
		salt:    salt,
	}, nil
}

const onSetCreds = "on authUsers.SetCreds(): "

func (authOp *authPassUsers) SetCreds(userKey identity.Key, toSet auth.Creds) (*auth.Creds, error) {
	var err error

	if userKey == "" {
		// TODO: set .Allowed = false and verify email

		user := auth.User{Creds: toSet}

		user.Key, err = authOp.usersOp.Save(users.Item{User: user, Allowed: true}, nil)
		if err != nil {
			return nil, errors.Wrapf(err, onSetCreds+"can't .usersOp.Save(users.Item{User: %#v, Allowed: true}, nil)", user)
		}

		return &user.Creds, nil
	}

	credsTypeToSet := auth.CredsType(toSet[auth.CredsToSet])
	delete(toSet, auth.CredsToSet)

	credsToSet, ok := toSet[credsTypeToSet]
	if !ok {
		return nil, errors.Errorf(onSetCreds+"no creds to set in %#v", toSet)
	}

	selector := selectors.Binary(selectors.Eq, users.UserKeyFieldName, selectors.Value{string(userKey)})

	items, err := authOp.usersOp.List(selector, nil)
	if err != nil {
		return nil, errors.Wrapf(err, onSetCreds+"can't .usersOp.List(selector = %#v, nil)", *selector)
	}
	if len(items) < 1 {
		return nil, errors.Errorf(onSetCreds+"no users with key %s)", userKey)
	} else if len(items) > 1 {
		return nil, errors.Errorf(onSetCreds+"too many users with key %s)", userKey)
	}

	if credsTypeToSet == auth.CredsEmail {
		// TODO: verify and/or another actions with some other creds types
	}

	items[0].Creds[credsTypeToSet] = credsToSet

	_, err = authOp.usersOp.Save(items[0], nil)
	if err != nil {
		return nil, errors.Wrapf(err, onSetCreds+"can't .usersOp.Save(%#v, nil)", items[0])

	}

	return &items[0].Creds, nil
}

const onAuthorize = "on authUsers.Authorize(): "

var reEmail = regexp.MustCompile("@")

func (authOp *authPassUsers) Authorize(toAuth auth.Creds) (*auth.User, error) {
	var selector *selectors.Term

	if login, ok := toAuth[auth.CredsLogin]; ok {
		if reEmail.MatchString(login) {
			selector = selectors.Binary(selectors.Eq, users.EmailFieldName, selectors.Value{login})
		} else {
			selector = selectors.Binary(selectors.Eq, users.NicknameFieldName, selectors.Value{login})
		}
	} else if email, ok := toAuth[auth.CredsEmail]; ok {
		selector = selectors.Binary(selectors.Eq, users.EmailFieldName, selectors.Value{email})
	} else if nickname, ok := toAuth[auth.CredsNickname]; ok {
		selector = selectors.Binary(selectors.Eq, users.NicknameFieldName, selectors.Value{nickname})
	} else {
		return nil, nil
		// return nil, errors.New(onAuthorize + "no login to auth")
	}

	//selector = logic.AND(
	//	selector,
	//	selectors.Binary(selectors.Gt, users.VerifiedFieldName, selectors.Value{0}),
	//)

	items, err := authOp.usersOp.List(selector, nil)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthorize+"can't .usersOp.List(selector = %#v, nil)", *selector)
	}

	maxUsersToAuthCheck := authOp.maxUsersToAuthCheck
	if len(items) < authOp.maxUsersToAuthCheck {
		maxUsersToAuthCheck = len(items)
	}

	for i := 0; i < maxUsersToAuthCheck; i++ {

		// TODO: use selector.AND (commented at the moment)
		if !items[i].Allowed {
			continue
		}

		item := items[i]

		if authOp.crypter.Verify(item.Creds[auth.CredsPasshash], []byte(toAuth[auth.CredsPassword])) == nil {
			user := item.User
			user.Creds = auth.Creds{
				auth.CredsNickname: item.Creds[auth.CredsNickname],
			}

			return &user, nil
		}
	}

	return nil, auth.ErrPassword
}

//func CheckCode() {
//	if toAuth.Type == auth.CredsSentCode {
//		// TODO: use "rector's" userIS to prevent second read
//		user, err := authOp.readByField("control_token", toAuth.FirstValue())
//		if err == nil {
//			err = authOp.usersOp.SetVerified(auth.IS(user.ID), user.ID)
//		}
//		if err != nil {
//			return nil, nil, errors.Wrap(err, onAuthenticateWithCode)
//		}
//
//		return &auth.User{
//			ID:       user.ID,
//			Nickname: "",
//			Accesses: nil,
//		}, nil, nil
//	}
//for _, ts := range toSet {
//if ts.Type == auth.CredsSentCode {
//user.ControlToken = strlib.RandomString(users.ControlTokenLength)
//_, err = authOp.usersOp.Update(auth.IS(user.ID), *user)
//if err != nil {
//return nil, nil, errors.Wrap(err, onSendCode+": can't set control token")
//}
//
//if !user.Verified {
//return nil, nil, authOp.sendCode(auth.Confirm, user.Email, user.ControlToken)
//}
//
//return nil, nil, authOp.sendCode(auth.SendCode, user.Email, user.ControlToken)
//}
//}
//}

//const onAuthenticateWithCode = "on authPassUsers.RegisterWithCode"
//const onSendCode = "on authPassUsers.SendCode"

//const createMessageSubject = "Підтвердження реєстрації / Registration confirmation: "
//const createMessageContent = "Щоб завершити реєстрацію, будь ласка, пройдіть за ланкою / To complete registration please follow the link: "
//
//const sendCodeMessageSubject = "Підтвердження зміни паролю / Password updating confirmation: "
//const sendCodeMessageContent = "Щоб змінити пароль, будь ласка, пройдіть за ланкою / To change your password please follow the link: "
//

//func (authOp *authPassUsers) sendCode(cbKey auth.Callback, email, token string) error {
//	cb, ok := authOp.callbacks[cbKey]
//	if !ok {
//		return errors.Errorf("wrong confirmation callback type: '%s'", cbKey)
//	}
//
//	message := connector.Message{
//		To:   email,
//		Body: token,
//	}
//
//	switch cbKey {
//	case auth.Confirm:
//		message.Subject = createMessageSubject
//		if authOp.useMessages {
//			message.Body = createMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
//		}
//	case auth.SendCode:
//		message.Subject = sendCodeMessageSubject
//		if authOp.useMessages {
//			message.Body = sendCodeMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
//		}
//	}
//
//	return authOp.senderOp.Send(message)
//}

//var reEmailToLogin1 = regexp.MustCompile(`@.*`)
//var reEmailToLogin2 = regexp.MustCompile(`(\.|-)`)

//type respFB struct {
//	TargetID    string `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
//
//const onQueryPartnerUser = "on authPassUsers.QueryPartnerUser"
//
//func (authOp *authPassUsers) QueryPartnerUser(partnerKey, partnerToken string) (partnerUser *auth.User, err error) {
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
//const onAuthenticateWithPartnerUser = "on authPassUsers.AuthenticateWithPartnerUser"
//
//func (authOp *authPassUsers) AuthenticateWithPartnerUser(partnerUser *auth.User) (*auth.User, error) {
//	user, err := authOp.readByLogin(partnerUser.Contacts)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthorize)
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
//func (authOp *authPassUsers) Clean(selector selectors.Selector) error {
//	return authOp.usersOp.Clean(selector)
//}
//
//func (authOp *authPassUsers) ConsoleDo(action string, values map[string]string) error {
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
//const onUpdatePasswordWithCode = "on authPassUsers.UpdatePasswordWithCode"
//
//func (authOp *authPassUsers) UpdatePasswordWithCode(controlToken string, encodedNewPassword encrlib.EncryptedPass) error {
//
//	if controlToken == "" {
//		return errors.Wrap(confidenter.ErrTokenNotFound, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.EncryptedPass == "" {
//		return errors.Wrap(auth.ErrBadPassword, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.Cryptype != encrlib.SHA256 {
//		return errors.Wrap(encrlib.ErrBadCryptype, onUpdatePasswordWithCode)
//	}
//
//	usrs, _, err := authOp.usersOp.ReadList(nil, &content.ListOptions{Selector: selectors.FieldEqual(users.ControlTokenField, selectors.Value{controlToken})})
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
//	user.EncryptedPass = encodedNewPassword.EncryptedPass
//	_, err = authOp.usersOp.Update(userIS, user)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": can't set new password")
//	}
//
//	return nil
//}
