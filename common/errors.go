package common

import (
	"errors"
)

type ErrorKey string

const CantPerformKey ErrorKey = "cant_perform"

const NoCredsKey ErrorKey = "no_creds"
const InvalidCredsKey ErrorKey = "invalid_creds"
const NoUserKey ErrorKey = "no_user"
const DuplicateUserKey ErrorKey = "duplicate_user"
const NoRightsKey ErrorKey = "no_rights"

const NotUniqueEmailKey ErrorKey = "not_unique_email"
const WrongPathKey ErrorKey = "wrong_path"
const WrongBodyKey ErrorKey = "wrong_body"
const WrongIDKey ErrorKey = "wrong_id"
const WrongJSONKey ErrorKey = "wrong_json"

const NotFoundKey ErrorKey = "not_found"

var ErrNotFound = errors.New(string(NotFoundKey))

const NullItemKey ErrorKey = "null_item"

var ErrNullItem = errors.New(string(NullItemKey))

const NotImplementedKey ErrorKey = "not_implemented"

var ErrNotImplemented = errors.New(string(NotImplementedKey))

const NotSupportedKey ErrorKey = "not_supported"

var ErrNotSupported = errors.New(string(NotSupportedKey))
