package auth

import "errors"

var ErrAuthenticated = errors.New("authorization required")
var ErrBadPassword = errors.New("wrong password")
var ErrBadIdentity = errors.New("bad identity")
