package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/common/addrlib"
	"github.com/pavlo67/workshop/basis/common/encrlib"
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
func (*identityECDSA) SetCreds(auth.User, ...auth.Creds) ([]auth.Creds, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	} else if privKey == nil {
		return nil, errEmptyPrivateKeyGenerated
	}

	privKeyBytes, err := encrlib.ECDSASerialize(*privKey)
	if err != nil {
		return nil, err
	}

	publKeyAddress := string(Proto) + string(append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...))

	creds := []auth.Creds{
		{Type: auth.CredsPrivateKey, Value: string(privKeyBytes)},
		{Type: auth.CredsPublicKeyAddress, Value: publKeyAddress},
	}

	return creds, nil
}

func (is *identityECDSA) Authorize(toAuth ...auth.Creds) (*auth.User, error) {
	var publKeyAddress, publKeyEncoded string
	var contentToSignature, numberToSignature, signature []byte

	for _, creds := range toAuth {
		switch creds.Type {
		case auth.CredsPublicKeyAddress:
			publKeyAddress = strings.TrimSpace(creds.Value)
			if len(publKeyAddress) < len(string(Proto)) || publKeyAddress[:len(string(Proto))] != string(Proto) {
				return nil, errWrongAddressProto
			}
			publKeyEncoded = publKeyAddress[len(string(Proto)):]

		case auth.CredsContentToSignature:
			contentToSignature = []byte(creds.Value)

		case auth.CredsNumberToSignature:
			numberToSignature = []byte(creds.Value)

		case auth.CredsSignature:
			signature = []byte(creds.Value)
		}
	}

	if len(publKeyEncoded) < 1 {
		return nil, errEmptyPublicKeyAddress
	}

	publKey := base58.Decode(publKeyEncoded)

	is.numberedMutex.Lock()
	if num, ok := is.numberedIDs[publKeyEncoded]; ok {
		numNew, _ := strconv.ParseUint(string(numberToSignature), 10, 64)
		if numNew <= num {
			is.numberedMutex.Unlock()
			return nil, errWrongNumber
		}
		is.numberedIDs[publKeyEncoded] = numNew
	}
	is.numberedMutex.Unlock()

	if !encrlib.ECDSAVerify(publKey, append(contentToSignature, numberToSignature...), signature) {
		return nil, errWrongSignature
	}

	return &auth.User{
		ID:   common.ID(publKeyAddress),
		Nick: publKeyAddress,
	}, nil
}

func (*identityECDSA) Accepts() ([]auth.CredsType, error) {
	return []auth.CredsType{auth.CredsSignature}, nil
}
