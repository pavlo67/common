package auth_persons

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GehirnInc/crypt"

	"github.com/pavlo67/common/common"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/selectors"
)

var _ auth.Operator = &authPersons{}

type authPersons struct {
	personsOp             persons.Operator
	maxPersonsToAuthCheck int
}

const onNew = "on authPersons.New()"

func New(personsOp persons.Operator, maxPersonsToAuthCheck int) (auth.Operator, error) {
	if personsOp == nil {
		return nil, errors.New(onNew + ": no persons.Operator")
	}
	if maxPersonsToAuthCheck < 1 {
		maxPersonsToAuthCheck = 1
	}

	return &authPersons{
		personsOp:             personsOp,
		maxPersonsToAuthCheck: maxPersonsToAuthCheck,
	}, nil
}

func hashCreds(creds, oldCreds auth.Creds) (auth.Creds, error) {

	if creds == nil {
		creds = auth.Creds{}
	}

	password := strings.TrimSpace(creds.StringDefault(auth.CredsPassword, ""))
	if password == "" {
		if oldCreds == nil {
			return nil, errata.KeyableError(errata.NoCredsKey, common.Map{"creds": creds, "reason": "no '" + auth.CredsPassword + "' key"})
		} else if creds != nil {
			creds[auth.CredsPasshash] = oldCreds[auth.CredsPasshash]
		}
	}

	crypt := crypt.SHA256.New()
	var salt []byte // TODO: generate salt
	hash, err := crypt.Generate([]byte(password), salt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't crypt.Generate(%s, %s)", password, salt))
	}

	creds[auth.CredsPasshash] = hash
	delete(creds, auth.CredsPassword)

	return creds, nil
}

const onSetCreds = "on authPersons.SetCreds()"

func (authOp *authPersons) SetCreds(authID auth.ID, toSet auth.Creds) (*auth.Creds, error) {
	if authID == "" {
		// TODO: set .Allowed = false and verify email

		toSetHashed, err := hashCreds(toSet, nil)
		if err != nil {
			return nil, errata.CommonError(err, onSetCreds)
		}

		identity := auth.Identity{
			Nickname: toSet.StringDefault(auth.CredsNickname, ""),
			Creds:    toSetHashed,
		}

		// TODO!!! hash password

		_, err = authOp.personsOp.Add(identity, nil, crud.OptionsWithRoles(rbac.RoleAdmin))
		if err != nil {
			return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.Save(%#v, nil)", identity)
		}

		return &toSet, nil
	}

	//if item.Creds, err = hashCreds(item.Creds, itemOld.Creds); err != nil {
	//	return nil, errata.CommonError(err, onChange)
	//}

	return nil, errata.NotImplemented

	//credsTypeToSet := auth.CredsType(toSet[auth.CredsToSet])
	//delete(toSet, auth.CredsToSet)
	//
	//credsToSet, ok := toSet[credsTypeToSet]
	//if !ok {
	//	return nil, errors.Errorf(onSetCreds+"no creds to set in %#v", toSet)
	//}
	//
	//selector := selectors.Binary(selectors.Eq, persons.UserKeyFieldName, selectors.Value{string(personKey)})
	//
	//items, err := authOp.personsOp.List(selector, nil)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.List(selector = %#v, nil)", *selector)
	//}
	//if len(items) < 1 {
	//	return nil, errors.Errorf(onSetCreds+"no persons with key %s)", personKey)
	//} else if len(items) > 1 {
	//	return nil, errors.Errorf(onSetCreds+"too many persons with key %s)", personKey)
	//}
	//
	//if credsTypeToSet == auth.CredsEmail {
	//	// TODO: verify and/or another actions with some other creds types
	//}
	//
	//items[0].Creds[credsTypeToSet] = credsToSet
	//
	//_, err = authOp.personsOp.Save(items[0], nil)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.Save(%#v, nil)", items[0])
	//
	//}
	//
	//return &items[0].Creds, nil
}

const onAuthenticate = "on authPersons.Authenticate()"

var reEmail = regexp.MustCompile("@")

func (authOp *authPersons) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
	var selector *selectors.Term

	nickname := toAuth.StringDefault(auth.CredsNickname, "")
	if nickname != "" {
		// selector = selectors.Binary(selectors.Eq, persons.NicknameFieldName, selectors.Value{nickname})
	} else {
		return nil, errata.KeyableError(errata.NoCredsKey, common.Map{"no nickname in creds": toAuth})
	}

	password := toAuth.StringDefault(auth.CredsPassword, "")

	//if login := toAuth.StringDefault(auth.CredsLogin, ""); login != "" {
	//	if reEmail.MatchString(login) {
	//		selector = selectors.Binary(selectors.Eq, persons.EmailFieldName, selectors.Value{login})
	//	} else {
	//		selector = selectors.Binary(selectors.Eq, persons.NicknameFieldName, selectors.Value{login})
	//	}
	//} else if email, ok := toAuth[auth.CredsEmail]; ok {
	//	selector = selectors.Binary(selectors.Eq, persons.EmailFieldName, selectors.Value{email})
	//} else if nickname, ok := toAuth[auth.CredsNickname]; ok {
	//	selector = selectors.Binary(selectors.Eq, persons.NicknameFieldName, selectors.Value{nickname})
	//} else {
	//	return nil, nil
	//	// return nil, errors.New(onAuthorize + "no login to auth")
	//}
	//selector = logic.AND(
	//	selector,
	//	selectors.Binary(selectors.Gt, persons.VerifiedFieldName, selectors.Value{0}),
	//)

	items, err := authOp.personsOp.List(crud.OptionsWithRoles(rbac.RoleAdmin)) // crud.Options{}.WithSelector(selector)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't .personsOp.List(selector = %#v, nil)", selector)
	}

	crypt := crypt.SHA256.New()

	for _, item := range items {
		if item.Nickname == nickname {
			if err = crypt.Verify(item.Creds.StringDefault(auth.CredsPasshash, ""), []byte(password)); err == nil {
				return &item.Identity, nil
			} else {
				l.Infof("can't verify %s on %s", item.Creds.StringDefault(auth.CredsPasshash, ""), password)
			}
		}
	}

	return nil, errata.KeyableError(errata.NoCredsKey, common.Map{onAuthenticate + ": wrong passhash in creds": toAuth})

	//maxPersonsToAuthCheck := authOp.maxPersonsToAuthCheck
	//if len(items) < authOp.maxPersonsToAuthCheck {
	//	maxPersonsToAuthCheck = len(items)
	//}
	//for i := 0; i < maxPersonsToAuthCheck; i++ {
	//
	//	//// TODO: use selector.AND (commented at the moment)
	//	//if !items[i].Allowed {
	//	//	continue
	//	//}
	//
	//	item := items[i]
	//
	//	if authOp.personsOp.CheckPassword(toAuth[auth.CredsPassword], item.Creds[auth.CredsPasshash]) {
	//		person := item.User
	//		person.Creds = auth.Creds{
	//			auth.CredsNickname: item.Creds[auth.CredsNickname],
	//		}
	//
	//		return &person, nil
	//	}
	//}

}

//func CheckCode() {
//	if toAuth.Type == auth.CredsSentCode {
//		// TODO: use "rector's" personIS to prevent second read
//		person, err := authOp.readByField("control_token", toAuth.FirstValue())
//		if err == nil {
//			err = authOp.personsOp.SetVerified(auth.IS(person.ID), person.ID)
//		}
//		if err != nil {
//			return nil, nil, errors.Wrap(err, onAuthenticateWithCode)
//		}
//
//		return &auth.User{
//			ID:       person.ID,
//			Nickname: "",
//			Accesses: nil,
//		}, nil, nil
//	}
//for _, ts := range toSet {
//if ts.Type == auth.CredsSentCode {
//person.ControlToken = strlib.RandomString(persons.ControlTokenLength)
//_, err = authOp.personsOp.Update(auth.IS(person.ID), *person)
//if err != nil {
//return nil, nil, errors.Wrap(err, onSendCode+": can't set control token")
//}
//
//if !person.Verified {
//return nil, nil, authOp.sendCode(auth.Confirm, person.Email, person.ControlToken)
//}
//
//return nil, nil, authOp.sendCode(auth.SendCode, person.Email, person.ControlToken)
//}
//}
//}

//const onAuthenticateWithCode = "on authPersons.RegisterWithCode"
//const onSendCode = "on authPersons.SendCode"

//const createMessageSubject = "Підтвердження реєстрації / Registration confirmation: "
//const createMessageContent = "Щоб завершити реєстрацію, будь ласка, пройдіть за ланкою / To complete registration please follow the link: "
//
//const sendCodeMessageSubject = "Підтвердження зміни паролю / Password updating confirmation: "
//const sendCodeMessageContent = "Щоб змінити пароль, будь ласка, пройдіть за ланкою / To change your password please follow the link: "
//

//func (authOp *authPersons) sendCode(cbKey auth.Callback, email, token string) error {
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
//const onQueryPartnerUser = "on authPersons.QueryPartnerUser"
//
//func (authOp *authPersons) QueryPartnerUser(partnerKey, partnerToken string) (partnerUser *auth.User, err error) {
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
//	// TODO: use respFB.Label to control person if his email is changed
//
//	return &auth.User{IdentityNamed: auth.IDentityNamed{Nick: data.Name}, Contacts: data.Email}, nil
//}
//
//const onAuthenticateWithPartnerUser = "on authPersons.AuthenticateWithPartnerUser"
//
//func (authOp *authPersons) AuthenticateWithPartnerUser(partnerUser *auth.User) (*auth.User, error) {
//	person, err := authOp.readByLogin(partnerUser.Contacts)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthorize)
//	}
//
//	if person != nil {
//		person.Contacts = partnerUser.Contacts
//		return person, nil
//	}
//
//	partnerUser.TargetID, err = authOp.personsOp.Create(nil, partnerUser, true)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthenticateWithPartnerUser+": can't create person")
//	}
//
//	return partnerUser, nil
//}
//
//func (authOp *authPersons) Clean(selector selectors.Selector) error {
//	return authOp.personsOp.Clean(selector)
//}
//
//func (authOp *authPersons) ConsoleDo(action string, values map[string]string) error {
//
//	//// TODO:  refactor it
//	//
//	//var sqlQuery string
//	//var err error
//	//var stmt *sql.Stmt
//	//var res sql.Result
//	//var valuesAll []interface{}
//	//if action == "new_person" {
//	//	valuesAll = []interface{}{values["login"], values["password"], values["email"]}
//	//	sqlQuery = "insert into `" + u.personTable + "` (nickname, passhash, email, verified, contacts, history) values (?,?,?,1,'','')"
//	//} else if action == "new_password" {
//	//	valuesAll = []interface{}{values["password"], values["login"]}
//	//	sqlQuery = "update `" + u.personTable + "` set passhash=? where nickname=?"
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
//const onUpdatePasswordWithCode = "on authPersons.UpdatePasswordWithCode"
//
//func (authOp *authPersons) UpdatePasswordWithCode(controlToken string, encodedNewPassword encrlib.EncryptedPass) error {
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
//	usrs, _, err := authOp.personsOp.ReadList(nil, &content.ListOptions{Selector: selectors.FieldEqual(persons.ControlTokenField, selectors.Value{controlToken})})
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error finding person")
//	}
//	if len(usrs) != 1 {
//		return fmt.Errorf(onUpdatePasswordWithCode+": can't find person (find %d)", len(usrs))
//	}
//
//	personIS := &basis.UserIS{joiner.SystemDomain(), "/person/", usrs[0].TargetID}
//
//	// TODO: use correct domain/path/, prevent duplicate reading
//	person, err := authOp.personsOp.Read(personIS, usrs[0].TargetID)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error reading person")
//	}
//
//	person.Cryptype = encodedNewPassword.Cryptype
//	person.EncryptedPass = encodedNewPassword.EncryptedPass
//	_, err = authOp.personsOp.Update(personIS, person)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": can't set new password")
//	}
//
//	return nil
//}
