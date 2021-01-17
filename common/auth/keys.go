package auth

import "github.com/pavlo67/workshop/common/joiner"

const InterfaceKey joiner.InterfaceKey = "auth"
const InterfaceOperatorKey joiner.InterfaceKey = "auth_operator"

// const InterfaceJWTKey joiner.InterfaceKey = "auth_jwt"
// const InterfaceJWTInternalKey joiner.InterfaceKey = "auth_jwt_internal"

const EPAuth = "authenticate_with_creds"

const AuthHandlerKey joiner.InterfaceKey = EPAuth

//const SetCredsHandlerKey joiner.InterfaceKey = "set_creds_handler"
//const GetCredsHandlerKey joiner.InterfaceKey = "get_creds_handler"
