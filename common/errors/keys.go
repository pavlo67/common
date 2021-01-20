package errors

import "errors"

type Key string

const CantPerformErr Key = ""

const NoCredsErr Key = "no_creds"
const InvalidCredsErr Key = "invalid_creds"
const NoUserErr Key = "no_user"
const NoRightsErr Key = "no_rights"
const OverdueRightsErr Key = "overdue_rights"

const NotUniqueEmailErr Key = "not_unique_email"
const WrongPathErr Key = "wrong_path"
const WrongBodyErr Key = "wrong_body"
const NotFoundErr Key = "not_found"
const WrongIDErr Key = "wrong_id"
const WrongJSONErr Key = "wrong_json"
const NotImplementedErr Key = "not_implemented"

var NotImplemented = errors.New("isn't implemented yet")
var NotFound = errors.New("not found")
var NullItem = errors.New("item is null")
