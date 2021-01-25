package errors

import "errors"

type Key string

const CantPerformErr Key = "cant_perform"

const NoCredsErr Key = "no_creds"
const InvalidCredsErr Key = "invalid_creds"
const NoUserErr Key = "no_user"
const DuplicateUserErr Key = "duplicate_user"
const NoRightsErr Key = "no_rights"
const OverdueRightsErr Key = "overdue_rights"

const NotUniqueEmailErr Key = "not_unique_email"
const WrongPathErr Key = "wrong_path"
const WrongBodyErr Key = "wrong_body"
const WrongIDErr Key = "wrong_id"
const WrongJSONErr Key = "wrong_json"

const NotFoundErr Key = "not_found"

var NotFound = errors.New(string(NotFoundErr))

var NullItem = errors.New("item is null")

const NotImplementedErr Key = "not_implemented"

var NotImplemented = errors.New(string(NotImplementedErr))
