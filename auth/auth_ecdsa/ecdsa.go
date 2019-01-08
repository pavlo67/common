package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis/addrlib"
	"github.com/pavlo67/punctum/basis/encrlib"
)

const Proto addrlib.Proto = "ecdsa://"

var _ auth.Operator = &identityECDSA{}

var errWrongAddressProto = errors.New("wrong address proto")
var errWrongSignature = errors.New("wrong signature")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type identityECDSA struct{}

func New() (auth.Operator, error) {
	return &identityECDSA{}, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityECDSA) SetCreds(*auth.ID, []auth.Creds, ...auth.Creds) (*auth.User, []auth.Creds, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	} else if privKey == nil {
		return nil, nil, errEmptyPrivateKeyGenerated
	}

	privKeyBytes, err := encrlib.ECDSASerialize(*privKey)
	if err != nil {
		return nil, nil, err
	}

	privKeyCreds := []auth.Creds{{
		Type:  auth.CredsPrivateKey,
		Value: string(privKeyBytes),
	}}

	publKeyAddress := string(Proto) + string(append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...))

	return &auth.User{
		ID:   auth.ID(publKeyAddress),
		Nick: publKeyAddress,
	}, privKeyCreds, nil
}

var reProto = regexp.MustCompile(`^ecdsa//\s*`) // `^(ecdsa|xxx|yyyy)//\s*`

func (*identityECDSA) Authorize(toAuth ...auth.Creds) (*auth.User, []auth.Creds, error) {
	var publKeyAddress, publKeyEncoded string
	var contentToSignature, signature []byte

	for _, creds := range toAuth {
		switch creds.Type {
		case auth.CredsPublicKeyAddress:
			publKeyAddress = strings.TrimSpace(creds.Value)
			publKeyEncoded = reProto.ReplaceAllString(publKeyAddress, "")
			if len(publKeyAddress) == len(publKeyEncoded) {
				return nil, nil, errWrongAddressProto
			}

		case auth.CredsContentToSignature:
			contentToSignature = []byte(creds.Value)

		case auth.CredsSignature:
			signature = []byte(creds.Value)
		}
	}

	if len(publKeyEncoded) < 1 {
		return nil, nil, errEmptyPublicKeyAddress
	}

	publKey := base58.Decode(publKeyEncoded)

	if !encrlib.ECDSAVerify(publKey, contentToSignature, signature) {
		return nil, nil, errWrongSignature
	}

	return &auth.User{
		ID:   auth.ID(publKeyAddress),
		Nick: publKeyAddress,
	}, nil, nil
}

func (*identityECDSA) Accepts() ([]auth.CredsType, error) {
	return []auth.CredsType{auth.CredsSignature}, nil
}
