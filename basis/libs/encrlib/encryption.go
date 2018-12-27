package encrlib

import (
	"fmt"
	"strings"

	_ "github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/pkg/errors"
)

type Cryptype string

const (
	SHA256  Cryptype = "SHA256"
	Provos  Cryptype = "Provos"
	NoCrypt Cryptype = ""
)

const CryptypePreferred Cryptype = SHA256

var ErrBadCryptype = errors.New("bad cryptography type")

type Hash struct {
	Passhash string   `bson:"passhash"           json:"passhash"`
	Cryptype Cryptype `bson:"cryptype,omitempty" json:"cryptype,omitempty"`
}

func PasswordValidation(password string, minLength int) (string, error) {
	password = strings.TrimSpace(password)

	if minLength > 0 && len(password) < minLength {
		return "", fmt.Errorf("закороткий пароль, повинно бути не менше %d символів", minLength)
	}
	return password, nil
}
