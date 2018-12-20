package basis

import (
	"encoding/json"
	"errors"
	"strings"
)

const ForUser = "info for user"

var ErrGlobalISNotFound = errors.New("record with this GlobalIS isn't found")

var ErrNotImplemented = errors.New("поки не імплементовано")
var ErrWrongDataType = errors.New("значення невідповідного типу")
var ErrNullItem = errors.New("порожнє значення")
var ErrDuplicateItem = errors.New("дублікат значення")
var ErrNotFound = errors.New("пошук не дав результатів")
var ErrCantCreate = errors.New("не вдається створити унікальний ключ")

var ErrAuthenticated = errors.New("необхідна авторизація")

var ErrBadGenus = errors.New("помилковий тип даних")
var ErrBadQuery = errors.New("помилка при запиті до даних")
var ErrEmptyQuery = errors.New("порожній пошуктовий запит")
var ErrJSONFormat = errors.New("неприйнятний формат даних, мав би бути JSON")

var ErrNoData = errors.New("відсутні дані")
var ErrCantDecodeData = errors.New("не вдається декодувати дані")
var ErrCantPerform = errors.New("не вдається виконати операцію")

var ErrTest = errors.New("помилка на тесті")

const CantPrepareQuery = "can't prepare (query='%s')"
const CantExecQuery = "can't exec (query='%s', values='%#v')"

// var ErrUserNotFound = errors.New("нема такого користувача")
// var ErrAuthorization = errors.New("помилковий логін/пароль")
// var ErrBadIdentityString = errors.New("bad confidenter string")

var ErrBadIdentity = errors.New("bad identity")

type Closer interface {
	Close() error
}

type Errors []error

// it's would be not good to use Errors as an error interface directly because of:

// func (errs Errors) Error() string {
//    return errs.String()
// }
// var errs basis.Errors; log.Println(errs == nil) // true
// var err  error;        log.Println(err  == nil) // true
// err = errs;            log.Println(err  == nil) // false

func MultiError(errors ...error) Errors {
	var errs Errors

	for _, err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (errs Errors) String() string {
	var errstrings []string
	for _, err := range errs {
		if err != nil {
			errstring := err.Error()
			if errstring == "" {
				errstrings = append(errstrings, "???")
			} else {
				errstrings = append(errstrings, errstring)
			}
		}
	}
	return strings.Join(errstrings, " / ")
}

func (errs Errors) Append(err error) Errors {
	if err != nil {
		return append(errs, err)
	}

	return errs
}

func (errs Errors) AppendErrs(errsToAppend Errors) Errors {
	if len(errs) == 0 {
		return errsToAppend
	}

	for _, err := range errsToAppend {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (errs Errors) Err() error {
	errstring := errs.String()
	if errstring != "" {
		return errors.New(errstring)
	}

	return nil
}

func (errs Errors) MarshalJSON() ([]byte, error) {
	messages := []string{}

	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	return json.Marshal(messages)
}

func (errs Errors) SetForUser() Errors {
	return append(Errors{errors.New(ForUser)}, errs...)
}

func (errs Errors) ForUser() Errors {
	if len(errs) > 0 && errs[0].Error() == ForUser {
		return errs[1:]
	}
	return nil
}
