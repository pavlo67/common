package basis

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrNotImplemented = errors.New("поки не імплементовано")

var ErrWrongDataType = errors.New("значення невідповідного типу")
var ErrNull = errors.New("відсутнє значення")
var ErrEmpty = errors.New("порожнє значення")
var ErrDuplicate = errors.New("дублікат значення")

var ErrEmptyQuery = errors.New("порожній запит")
var ErrBadQuery = errors.New("помилковий запит")
var ErrNotFound = errors.New("не знайдено")
var ErrCantPerform = errors.New("не вдається виконати операцію")
var ErrCantDecodeData = errors.New("не вдається декодувати дані")
var ErrNoData = errors.New("відсутні дані")
var ErrJSONFormat = errors.New("неприйнятний формат даних, мав би бути JSON")

const CantPrepareQuery = "can't prepare (query='%s')"
const CantExecQuery = "can't execute (query='%s', values='%#v')"

var ErrTest = errors.New("помилка на тесті")

// Errors ---------------------------------------------------------------------------------------------------------------

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

//const ForUser = "info for user"
//func (errs Errors) SetForUser() Errors {
//	return append(Errors{errors.New(ForUser)}, errs...)
//}
//
//func (errs Errors) ForUser() Errors {
//	if len(errs) > 0 && errs[0].Error() == ForUser {
//		return errs[1:]
//	}
//	return nil
//}
