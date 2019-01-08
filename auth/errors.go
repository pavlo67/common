package auth

import "errors"

var ErrAuthenticated = errors.New("необхідна авторизація")
var ErrBadIdentity = errors.New("bad identity")
