package common

type ErrorKey = Key

const CantPerformErr ErrorKey = ""

const NoCredsErr ErrorKey = "no_creds"
const InvalidCredsErr ErrorKey = "invalid_creds"
const NoUserErr ErrorKey = "no_user"
const NoRightsErr ErrorKey = "no_rights"
const OverdueRightsErr ErrorKey = "overdue_rights"

const NotUniqueEmailErr ErrorKey = "not_unique_email"
const WrongPathErr ErrorKey = "wrong_path"
const WrongBodyErr ErrorKey = "wrong_body"
const NotFoundErr ErrorKey = "not_found"
const WrongIDErr ErrorKey = "wrong_id"
const WrongJSONErr ErrorKey = "wrong_json"
const NotImplementedErr ErrorKey = "not_implemented"
