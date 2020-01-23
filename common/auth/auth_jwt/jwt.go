package auth_jwt

import (
	"crypto/rsa"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/libraries/encrlib"
)

const Proto = "jwt"

var _ auth.Operator = &authJWT{}

//var errEmptyPublicKeyAddress = errors.New("empty public Key address")
//var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type authJWT struct {
	privKey rsa.PrivateKey
	builder jwt.Builder
}

// TODO!!! add expiration time

const onNew = "on auth_jwt.New()"

func New(pathToStore string) (auth.Operator, error) {
	privKey, err := encrlib.NewRSAPrivateKey(pathToStore)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	signerOpts := (&jose.SignerOptions{}).WithType("JWT") // signerOpts.WithType("JWT")
	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: privKey}
	rsaSigner, err := jose.NewSigner(signingKey, signerOpts)
	if err != nil {
		return nil, errors.Wrapf(err, onNew+": can't jose.NewSigner(%#v, %#v)", signingKey, signerOpts)
	}

	return &authJWT{privKey: *privKey, builder: jwt.Signed(rsaSigner)}, nil
}

type jwtCreds struct {
	*jwt.Claims
	Creds auth.Creds `json:"creds,omitempty"`
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (authOp *authJWT) SetCreds(userKey identity.Key, creds auth.Creds) (*auth.Creds, error) {

	jc := jwtCreds{
		Claims: &jwt.Claims{
			//Issuer:   "issuer1",
			//Subject:  "subject1",
			// Audience: jwt.Audience{"aud1", "aud2"},
			ID:       string(userKey),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Expiry:   jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 8, 0, 0, time.UTC)),
		},

		Creds: auth.Creds{
			auth.CredsNickname: creds[auth.CredsNickname],
			// TODO: add some other creds...
		},
	}

	// add claims to the Builder
	builder := authOp.builder.Claims(jc)

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		return nil, errors.Wrap(err, "on authJWT.SetCreds() with builder.CompactSerialize()")
	}

	delete(creds, auth.CredsToSet)

	creds[auth.CredsJWT] = rawJWT

	return &creds, nil
}

func (authOp *authJWT) Authorize(toAuth auth.Creds) (*auth.User, error) {
	credsJWT, ok := toAuth[auth.CredsJWT]
	if !ok {
		return nil, nil
	}

	parsedJWT, err := jwt.ParseSigned(credsJWT)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse JWT: %s", credsJWT)
	}

	res := jwtCreds{}
	err = parsedJWT.Claims(&authOp.privKey.PublicKey, &res)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claims: %#v", parsedJWT)
	}

	return &auth.User{
		Key:   identity.Key(res.ID),
		Creds: res.Creds,
	}, nil
}

//func (*authJWT) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsSignature}, nil
//}
