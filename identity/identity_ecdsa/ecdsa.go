package identity_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"regexp"
	"strings"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/encryption"
	"github.com/pavlo67/punctum/identity"
)

const Proto basis.Proto = "ecdsa://"

var _ identity.Operator = &identityECDSA{}

var errWrongAddressProto = errors.New("wrong address proto")
var errWrongSignature = errors.New("wrong signature")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type identityECDSA struct{}

func New() (identity.Operator, error) {
	return &identityECDSA{}, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityECDSA) SetCreds(*identity.ID, []identity.Creds, ...identity.Creds) (*identity.User, []identity.Creds, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	} else if privKey == nil {
		return nil, nil, errEmptyPrivateKeyGenerated
	}

	privKeyBytes, err := encryption.ECDSASerialize(*privKey)
	if err != nil {
		return nil, nil, err
	}

	privKeyCreds := []identity.Creds{{
		Type:  identity.CredsPrivateKey,
		Value: string(privKeyBytes),
	}}

	publKeyAddress := string(Proto) + string(append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...))

	return &identity.User{
		ID:       identity.ID(publKeyAddress),
		Nickname: publKeyAddress,
	}, privKeyCreds, nil
}

var reProto = regexp.MustCompile(`^ecdsa//\s*`) // `^(ecdsa|xxx|yyyy)//\s*`

func (*identityECDSA) Authorize(toAuth ...identity.Creds) (*identity.User, []identity.Creds, error) {
	var publKeyAddress, publKeyEncoded string
	var contentToSignature, signature []byte

	for _, creds := range toAuth {
		switch creds.Type {
		case identity.CredsPublicKeyAddress:
			publKeyAddress = strings.TrimSpace(creds.Value)
			publKeyEncoded = reProto.ReplaceAllString(publKeyAddress, "")
			if len(publKeyAddress) == len(publKeyEncoded) {
				return nil, nil, errWrongAddressProto
			}

		case identity.CredsContentToSignature:
			contentToSignature = []byte(creds.Value)

		case identity.CredsSignature:
			signature = []byte(creds.Value)
		}
	}

	if len(publKeyEncoded) < 1 {
		return nil, nil, errEmptyPublicKeyAddress
	}

	publKey := base58.Decode(publKeyEncoded)

	if !encryption.ECDSAVerify(publKey, contentToSignature, signature) {
		return nil, nil, errWrongSignature
	}

	return &identity.User{
		ID:       identity.ID(publKeyAddress),
		Nickname: publKeyAddress,
	}, nil, nil
}

func (*identityECDSA) Accepts() ([]identity.CredsType, error) {
	return []identity.CredsType{identity.CredsSignature}, nil
}
