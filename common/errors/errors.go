package errors

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/pavlo67/common/common"
)

type Error interface {
	error
	Key() Key
	Data() common.Map
	Append(interface{}) Error
}

func CommonError(any ...interface{}) Error {
	var err *commonError

	for _, anything := range any {
		err = err.append(anything)
	}

	return err
}

func KeyableError(err error, key Key, data common.Map) Error {
	return &commonError{
		errs: Errors{err},
		key:  key,
		data: data,
	}
}

// Errors ---------------------------------------------------------------------------------------------------------------

var _ Error = &commonError{}

type commonError struct {
	errs Errors
	key  Key
	data common.Map
}

func (ce *commonError) Error() string {
	if ce == nil {
		return ""
	}
	return ce.errs.String()
}

func (ce *commonError) Key() Key {
	if ce == nil {
		return ""
	}

	return ce.key
}

func (ce *commonError) Data() common.Map {
	if ce == nil {
		return nil
	}

	return ce.data
}

func (ce *commonError) append(anything interface{}) *commonError {
	if anything == nil {
		return ce
	}

	if ce == nil {
		switch v := anything.(type) {
		case commonError:
			return &v
		case *commonError:
			return v
		case Error:
			return &commonError{
				errs: Errors{errors.New(v.Error())},
				key:  v.Key(),
				data: v.Data(),
			}
		case error:
			return &commonError{errs: Errors{v}}
		case string:
			return &commonError{errs: Errors{errors.New(v)}}
		}
		ce = &commonError{}
	}

	var err error
	switch v := anything.(type) {
	case commonError:
		err = &v
	case *commonError:
		err = v
	case Error:
		err = v
	case error:
		err = v
	case string:
		err = errors.New(v)
	default:
	}

	ce.errs = append(ce.errs, err)
	return ce
}

func (ce *commonError) Append(anything interface{}) Error {
	return ce.append(anything)

}

// Errors ------------------------------------------------------------------------------------------------

// DEPRECATED
type Errors []error

// DEPRECATED
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

// DEPRECATED
func (errs Errors) Append(err error) Errors {
	if err != nil {
		return append(errs, err)
	}

	return errs
}

// DEPRECATED
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

// DEPRECATED
func (errs Errors) Err() error {

	// TODO!!! errs.Error() must keep Keyable interface

	errstring := errs.String()
	if errstring != "" {
		return errors.New(errstring)
	}

	return nil
}

// DEPRECATED
func (errs Errors) MarshalJSON() ([]byte, error) {
	messages := []string{}

	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	return json.Marshal(messages)
}
