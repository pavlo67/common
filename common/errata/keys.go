package errata

type Key string

const CantPerformKey Key = "cant_perform"

const NoCredsKey Key = "no_creds"
const InvalidCredsKey Key = "invalid_creds"
const NoUserKey Key = "no_user"
const DuplicateUserKey Key = "duplicate_user"
const NoRightsKey Key = "no_rights"

const NotUniqueEmailKey Key = "not_unique_email"
const WrongPathKey Key = "wrong_path"
const WrongBodyKey Key = "wrong_body"
const WrongIDKey Key = "wrong_id"
const WrongJSONKey Key = "wrong_json"

const NotFoundKey Key = "not_found"

// var NotFound = errors.New(string(NotFoundKey))

const NullItemKey Key = "null_item"

// var NullItem = errors.New(string(NullItemKey))

const NotImplementedKey Key = "not_implemented"

var NotImplemented = New(string(NotImplementedKey))
