package auth

import "errors"

var ErrAuthRequired = errors.New("authorization required")
var ErrPassword = errors.New("wrong password")
var ErrSignaturedKey = errors.New("wrong signatured key")
var ErrAuthSession = errors.New("wrong authorization session")
var ErrEncryptionType = errors.New("wrong encryption type")
var ErrIP = errors.New("wrong IP")
var ErrNoCreds = errors.New("no creds")
var ErrNoUser = errors.New("no user")

//var ErrBadIdentity = errors.New("bad identity")
