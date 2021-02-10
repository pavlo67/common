package auth

import (
	"errors"

	"github.com/pavlo67/common/common"
)

const InterfaceKey common.InterfaceKey = "auth"

const IntefaceKeyAuthenticate common.InterfaceKey = "auth_authenticate"
const IntefaceKeySetCreds common.InterfaceKey = "auth_set_creds"

var ErrAuthRequired = errors.New("authorization required")
var ErrPassword = errors.New("wrong password")
var ErrSignaturedKey = errors.New("wrong signatured key")
var ErrAuthSession = errors.New("wrong authorization session")
var ErrEncryptionType = errors.New("wrong encryption type")
var ErrIP = errors.New("wrong IP")
var ErrNoCreds = errors.New("no creds")
var ErrNoUser = errors.New("no user")

//var ErrBadIdentity = errors.New("bad identity")
