package auth_jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"time"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/common/addrlib"
	"github.com/pkg/errors"
)

const Proto addrlib.Proto = "jwt://"

var _ auth.Operator = &authJWT{}

//var errEmptyPublicKeyAddress = errors.New("empty public Key address")
//var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type authJWT struct {
	privKey rsa.PrivateKey
	builder jwt.Builder
}

func New() (auth.Operator, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating random key: %s", err)
	}

	key := jose.SigningKey{Algorithm: jose.RS256, Key: privKey}

	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")
	rsaSigner, err := jose.NewSigner(key, &signerOpts)
	if err != nil {
		log.Fatalf("failed to create signer:%+v", err)
	}

	return &authJWT{
		privKey: *privKey,
		builder: jwt.Signed(rsaSigner),
	}, nil
}

type jwtCreds struct {
	*jwt.Claims
	Creds []auth.Creds `json:"creds,omitempty"`
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (authOp *authJWT) SetCreds(user auth.User, creds ...auth.Creds) ([]auth.Creds, error) {
	jc := jwtCreds{
		Claims: &jwt.Claims{
			//Issuer:   "issuer1",
			//Subject:  "subject1",
			// Audience: jwt.Audience{"aud1", "aud2"},
			ID:       string(user.ID),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Expiry:   jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 8, 0, 0, time.UTC)),
		},
		Creds: append(user.Creds, append(creds, auth.Creds{Type: auth.CredsNickname, Value: user.Nick})...),
	}
	// add claims to the Builder
	builder := authOp.builder.Claims(jc)

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		return nil, errors.Wrap(err, "on authJWT.SetCreds() with builder.CompactSerialize()")
	}

	return []auth.Creds{{
		Type:  auth.CredsJWT,
		Value: rawJWT,
	}}, nil
}

func (authOp *authJWT) Authorize(toAuth ...auth.Creds) (*auth.User, error) {
	for _, creds := range toAuth {
		switch creds.Type {
		case auth.CredsJWT:
			parsedJWT, err := jwt.ParseSigned(creds.Value)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse JWT: %s", creds.Value)
			}

			res := jwtCreds{}
			err = parsedJWT.Claims(&authOp.privKey.PublicKey, &res)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get claims: %#v", parsedJWT)
			}

			var nick string
			var creds []auth.Creds

			for _, c := range res.Creds {
				if c.Type == auth.CredsNickname {
					nick = c.Value
				} else {
					creds = append(creds, c)
				}
			}

			return &auth.User{
				ID:    common.ID(res.ID),
				Nick:  nick,
				Creds: creds,
			}, nil
		}
	}

	return nil, nil
}

//func (*authJWT) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsSignature}, nil
//}
