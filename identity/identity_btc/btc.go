package identity_btc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/encryption"
	"github.com/pavlo67/punctum/identity"
)

const Proto basis.Proto = "btc://"

var _ identity.Operator = &identityBTC{}

var errWrongAddressProto = errors.New("wrong address proto")
var errWrongSignature = errors.New("wrong signature")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type identityBTC struct{}

func New() (identity.Operator, error) {
	return &identityBTC{}, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityBTC) SetCreds(*identity.ID, []identity.Creds, ...identity.Creds) (*identity.User, []identity.Creds, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	} else if privateKey == nil {
		return nil, nil, errEmptyPrivateKeyGenerated
	}

	privateKeyBytes, err := encryption.ECDSASerialize(*privateKey)
	if err != nil {
		return nil, nil, err
	}

	privateKeyCreds := []identity.Creds{{
		Type:  identity.CredsPrivateKey,
		Value: string(privateKeyBytes),
	}}

	publicKeyAddress := string(Proto) + string(append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...))

	return &identity.User{
		ID:       identity.ID(publicKeyAddress),
		Nickname: publicKeyAddress,
	}, privateKeyCreds, nil
}

// var reProto = regexp.MustCompile(`^\w+://\s*`)

func (*identityBTC) Authorize(toAuth ...identity.Creds) (*identity.User, []identity.Creds, error) {
	var publicKeyAddress, publicKeyEncoded string
	var contentToSignature, signature []byte

	for _, creds := range toAuth {
		switch creds.Type {
		case identity.CredsPublicKeyAddress:
			// publicKeyAddress = reProto.ReplaceAllString(strings.TrimSpace(creds.Value), "")
			publicKeyAddress = strings.TrimSpace(creds.Value)
			if publicKeyAddress[:len(Proto)] != string(Proto) {
				return nil, nil, errWrongAddressProto
			}
			publicKeyEncoded = publicKeyAddress[len(Proto):]

		case identity.CredsContentToSignature:
			contentToSignature = []byte(creds.Value)

		case identity.CredsSignature:
			signature = []byte(creds.Value)
		}
	}

	if len(publicKeyEncoded) < 1 {
		return nil, nil, errEmptyPublicKeyAddress
	}

	publicKey := base58.Decode(publicKeyEncoded)

	if !encryption.ECDSAVerify(publicKey, contentToSignature, signature) {
		return nil, nil, errWrongSignature
	}

	return &identity.User{
		ID:       identity.ID(publicKeyAddress),
		Nickname: publicKeyAddress,
	}, nil, nil
}

func (*identityBTC) Accepts() ([]identity.CredsType, error) {
	return []identity.CredsType{identity.CredsSignature}, nil
}
