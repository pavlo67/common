package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"strconv"
	"sync"

	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis/addrlib"
	"github.com/pavlo67/constructor/basis/encrlib"
)

const Proto addrlib.Proto = "ecdsa://"

var _ auth.Operator = &identityECDSA{}

var errWrongAddressProto = errors.New("wrong address proto")
var errWrongSignature = errors.New("wrong signature")
var errWrongNumber = errors.New("wrong user's number")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

type identityECDSA struct {
	numberedIDs   map[string]uint64
	numberedMutex *sync.Mutex
}

func New(numberedIDs []string) (auth.Operator, error) {
	is := &identityECDSA{
		numberedIDs:   map[string]uint64{},
		numberedMutex: &sync.Mutex{},
	}
	for _, id := range numberedIDs {
		is.numberedIDs[id] = 0
	}

	return is, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityECDSA) SetCreds(*auth.ID, []auth.Creds) (*auth.User, []auth.Creds, error) {
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

const proto = `ecdsa://`

func (is *identityECDSA) Authorize(toAuth []auth.Creds) (*auth.User, []auth.Creds, error) {
	var publKeyAddress, publKeyEncoded string
	var contentToSignature, numberToSignature, signature []byte

	for _, creds := range toAuth {
		switch creds.Type {
		case auth.CredsPublicKeyAddress:
			publKeyAddress = strings.TrimSpace(creds.Value)
			if len(publKeyAddress) < len(proto) || publKeyAddress[:len(proto)] != proto {
				return nil, nil, errWrongAddressProto
			}
			publKeyEncoded = publKeyAddress[len(proto):]

		case auth.CredsContentToSignature:
			contentToSignature = []byte(creds.Value)

		case auth.CredsNumberToSignature:
			numberToSignature = []byte(creds.Value)

		case auth.CredsSignature:
			signature = []byte(creds.Value)
		}
	}

	if len(publKeyEncoded) < 1 {
		return nil, nil, errEmptyPublicKeyAddress
	}

	publKey := base58.Decode(publKeyEncoded)

	is.numberedMutex.Lock()
	if num, ok := is.numberedIDs[publKeyEncoded]; ok {
		numNew, _ := strconv.ParseUint(string(numberToSignature), 10, 64)
		if numNew <= num {
			is.numberedMutex.Unlock()
			return nil, nil, errWrongNumber
		}
		is.numberedIDs[publKeyEncoded] = numNew
	}
	is.numberedMutex.Unlock()

	if !encrlib.ECDSAVerify(publKey, append(contentToSignature, numberToSignature...), signature) {
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
