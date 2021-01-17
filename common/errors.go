package common

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrNotImplemented = errors.New("isn't implemented yet")
var ErrNotFound = errors.New("not found")
var ErrNullItem = errors.New("item is null")

type Key string
type Error interface {
	error
	Key() Key
	Data() Map
	Errors() Errors
	Err() error
	//	Err() error
}

// CommonError ------------------------------------------------------------------------------------------------------

var _ Error = &CommonError{}

func WrappedError(err error, errs0 ...error) Error {
	var errs Errors
	var key Key

	if len(errs0) < 1 {
		// no errs
	} else if len(errs0) > 1 {
		errs = errs0
	} else if errs, _ = interface{}(errs0[0]).(Errors); errs != nil {
		// errs is errs0[0].(Errors)
	} else {
		errs = Errors{errs0[0]}
	}

	if keyable, _ := err.(Error); keyable != nil {
		key = keyable.Key()
	}

	return &CommonError{
		Err0:    err,
		Key0:    key,
		Message: errs.String(),
	}
}

type CommonError struct {
	Err0    error
	Key0    Key
	Message string
}

func (commonError *CommonError) Error() string {
	if commonError == nil {
		return ""
	}
	if commonError.Message != "" {
		return commonError.Err0.Error() + " / " + commonError.Message
	}

	return commonError.Err0.Error()
}

func (commonError *CommonError) Key() Key {
	if commonError == nil {
		return ""
	}

	return commonError.Key0

}

//func (wrappedError *CommonError) Err() error {
//	if wrappedError == nil {
//		return nil
//	}
//	return wrappedError.Error
//}

func (commonError *CommonError) Errors() Errors {
	if commonError == nil {
		return nil
	}

	if commonError.Message != "" {
		return append(Errors{commonError.Err0}, errors.New(commonError.Message))
	}
	return Errors{commonError.Err0}
}

func (commonError *CommonError) Err() error {
	if commonError == nil {
		return nil
	}

	if commonError.Message != "" {
		return append(Errors{commonError.Err0}, errors.New(commonError.Message)).Err()
	}
	return commonError.Err0
}

func (commonError *CommonError) Data() Map {
	return nil
}

// StructuredError ------------------------------------------------------------------------------------------------------

type StructuredError struct {
	key  Key
	data Map
	errs Errors
}

var _ Error = &StructuredError{}

func KeyableError(errorKey Key, data Map, err error) Error {
	var errs Errors
	if keyable, ok := err.(Error); ok {
		if errorKey == "" {
			errorKey = keyable.Key()
		}
		errs = keyable.Errors()
	} else if err != nil {
		errs = Errors{err}
	} else {
		errs = Errors{errors.New("undefined error")}
	}

	structuredError := StructuredError{
		key:  errorKey,
		data: data,
		errs: errs,
	}

	return &structuredError
}

func (structuredError *StructuredError) Key() Key {
	if structuredError == nil {
		return ""
	}

	return structuredError.key
}

func (structuredError *StructuredError) Error() string {
	if structuredError == nil {
		return ""
	}
	return append(Errors{errors.New("key: " + string(structuredError.key))}, structuredError.errs...).String()
}

func (structuredError *StructuredError) Err() error {
	if structuredError == nil {
		return nil
	}
	return append(Errors{errors.New("key: " + string(structuredError.key))}, structuredError.errs...).Err()
}

func (structuredError *StructuredError) Data() Map {
	if structuredError == nil {
		return nil
	}
	return structuredError.data
}

func (structuredError *StructuredError) Errors() Errors {
	if structuredError == nil {
		return nil
	}
	return structuredError.errs
}

// Errors ---------------------------------------------------------------------------------------------------------------

type Errors []error

// it's would be not good to use Errors as an error interface directly because of:

// func (errs Errors) ErrStr() string {
//    return errs.IDStr()
// }
// var errs basis.Errors; logPrintln(errs == nil) // true
// var err  error;        logPrintln(err  == nil) // true
// err = errs;            logPrintln(err  == nil) // false

//func MultiError(errors ...error) Errors {
//	var errs Errors
//
//	for _, err := range errors {
//		if err != nil {
//			errs = append(errs, err)
//		}
//	}
//
//	return errs
//}

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

	// TODO!!! errs.Err() must keep Keyable interface

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
