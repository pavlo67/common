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

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/libs/addrlib"
	"github.com/pavlo67/workshop/common/libs/encrlib"
	"github.com/pavlo67/workshop/components/auth"
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
func (*identityECDSA) SetCreds(auth.User, auth.Creds) (*auth.Creds, error) {
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

	creds := auth.Creds{
		Values: map[auth.CredsType]string{
			auth.CredsPrivateKey:       string(privKeyBytes),
			auth.CredsPublicKeyAddress: publKeyAddress,
		},
	}

	return &creds, nil
}

func (is *identityECDSA) Authorize(toAuth auth.Creds) (*auth.User, error) {
	publKeyAddress := strings.TrimSpace(toAuth.Values[auth.CredsPublicKeyAddress])

	var publKeyEncoded string
	if len(publKeyAddress) < len(string(Proto)) || publKeyAddress[:len(string(Proto))] != string(Proto) {
		return nil, errWrongAddressProto
	}
	publKeyEncoded = publKeyAddress[len(string(Proto)):]

	contentToSignature := []byte(toAuth.Values[auth.CredsContentToSignature])
	numberToSignature := []byte(toAuth.Values[auth.CredsNumberToSignature])
	signature := []byte(toAuth.Values[auth.CredsSignature])

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
		ID:       common.ID(publKeyAddress),
		Nickname: publKeyAddress,
	}, nil
}

func (*identityECDSA) Accepts() ([]auth.CredsType, error) {
	return []auth.CredsType{auth.CredsSignature}, nil
}
