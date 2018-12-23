package encryption

import (
	"github.com/GehirnInc/crypt"
	"github.com/pkg/errors"
)

func SHA256Hash(str, salt string) (string, error) {
	crypt := crypt.SHA256.New()
	hash, err := crypt.Generate([]byte(str), []byte(salt))

	return hash, errors.Wrap(err, "error hashing str")
}
