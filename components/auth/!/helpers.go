package a

import (
	"strings"

	"log"

	"github.com/GehirnInc/crypt"
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis/encrlib"
	"github.com/pavlo67/punctum/starter/joiner"
)

func CheckPassword(salt, passhash string, cryptypeToMustBe encrlib.Cryptype, passhashToMustBe string) bool {
	crypt := crypt.SHA256.New()
	passhash, _ = crypt.Generate([]byte(strings.TrimSpace(passhash)), []byte(salt))
	return passhash == passhashToMustBe
}

type UserToCreate struct {
	Nickname string
	Password string
	Email    string
}

func NewUserToCreate(creds []Creds) UserToCreate {
	var user UserToCreate

	for _, cr := range creds {
		switch cr.Type {
		case CredsNickname:
			user.Nickname = cr.FirstValue()
		case CredsEmail:
			user.Email = cr.FirstValue()
		case CredsPassword:
			user.Password = cr.FirstValue()
		default:
			log.Printf("WARNING: unused creds for authusers.Create: %#v", cr)
		}
	}

	return user
}

func IS(id string) auth.ID {
	//return auth.IDentity{
	//	Domain: joiner.SystemDomain(),
	//	WithParams:   UserPath,
	//	TargetID:     id,
	//}.String()

	return auth.ID(joiner.SystemDomain() + "/" + UserPath + "/" + id)
}
