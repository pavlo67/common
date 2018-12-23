package identity_btc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"io"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/identity"
)

var _ identity.Operator = &identityBTC{}

var ErrWrongSignature = errors.New("wrong signature")
var ErrEmptyPublicKeyAddress = errors.New("empty public Key address")

type identityBTC struct{}

func New() (identity.Operator, error) {
	return &identityBTC{}, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityBTC) SetCreds(*basis.ID, []identity.Creds, ...identity.Creds) (*identity.User, []identity.Creds, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	privateCreds := []identity.Creds{{
		Type:  identity.CredsPrivateKey,
		Value: string(privateKey),
	}}

	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	pk := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.Curve{},
			X:     &big.Int{},
			Y:     &big.Int{},
		},
		D: &big.Int{},
	}

	return &identity.User{
		ID:       "",
		Nickname: "",
		Accesses: nil,
	}, nil, basis.ErrNotImplemented
}

func (*identityBTC) Authorize(toAuth ...identity.Creds) (*identity.User, []identity.Creds, error) {
	var publicKeyAddress, contentToSignature, signature string

	for _, creds := range toAuth {
		switch creds.Type {
		case identity.CredsPublicKeyAddress:
			publicKeyAddress = strings.TrimSpace(creds.Value)
		case identity.CredsContentToSignature:
			contentToSignature = creds.Value
		case identity.CredsSignature:
			signature = creds.Value
		}
	}

	if publicKeyAddress == "" {
		return nil, nil, ErrEmptyPublicKeyAddress
	}
	publicKey := base58.Decode(publicKeyAddress)

	h := md5.New()
	io.WriteString(h, contentToSignature)
	data := h.Sum(nil)

	// build key and verify data
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes([]byte(signature)[:(sigLen / 2)])
	s.SetBytes([]byte(signature)[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKey)
	x.SetBytes(publicKey[:(keyLen / 2)])
	y.SetBytes(publicKey[(keyLen / 2):])

	curve := elliptic.P256()

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	if !ecdsa.Verify(&rawPubKey, data, &r, &s) {
		return nil, nil, ErrWrongSignature
	}

	return &identity.User{
		ID:       basis.ID("btc://" + publicKeyAddress),
		Nickname: publicKeyAddress,
	}, nil, nil
}

func (*identityBTC) Accepts() ([]identity.CredsType, error) {
	return []identity.CredsType{identity.CredsSignature}, nil
}
