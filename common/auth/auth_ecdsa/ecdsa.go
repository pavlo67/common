package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	r "math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pavlo67/workshop/common/identity"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/libraries/encrlib"
	"github.com/pavlo67/workshop/common/libraries/strlib"
)

const Cryptype encrlib.Cryptype = "ecdsa"
const Proto = "ecdsa"

var _ auth.Operator = &identityECDSA{}

var errWrongSignature = errors.New("wrong signature")
var errWrongNumber = errors.New("wrong user's number")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

var cnt uint32

type Session struct {
	IP        string
	StartedAt time.Time
}

type identityECDSA struct {
	sessions map[uint64]Session
	mutex    *sync.Mutex

	maxSessionDuration time.Duration
	numbersLimit       int

	acceptableIDs []string
}

func New(numbersLimit int, maxSessionDuration time.Duration, acceptableIDs []string) (auth.Operator, error) {
	r.Seed(time.Now().UnixNano())

	is := &identityECDSA{
		sessions: map[uint64]Session{},
		mutex:    &sync.Mutex{},

		maxSessionDuration: maxSessionDuration,
		numbersLimit:       numbersLimit,

		acceptableIDs: acceptableIDs,
	}

	return is, nil
}

func (is *identityECDSA) InitAuthSession(toInit auth.Creds) (*auth.Creds, error) {
	now := time.Now()

	is.mutex.Lock() // Lock() -----------------------------------------------------

	if is.numbersLimit > 0 && len(is.sessions) >= is.numbersLimit {
		for n, s := range is.sessions {
			if now.Sub(s.StartedAt) >= is.maxSessionDuration {
				delete(is.sessions, n)
			}
		}
	}

	cnt++
	numberToSend := uint64(cnt)<<32 + uint64(r.Uint32())

	is.sessions[numberToSend] = Session{
		IP:        toInit.Values[auth.CredsIP], // TODO??? check if IP isn't empty
		StartedAt: now,
	}

	is.mutex.Unlock() // Unlock() -------------------------------------------------

	return &auth.Creds{
		Cryptype: encrlib.NoCrypt,
		Values:   auth.Values{auth.CredsKeyToSignature: strconv.FormatUint(numberToSend, 10)},
	}, nil
}

func (is *identityECDSA) Authorize(toAuth auth.Creds) (*auth.User, error) {
	if toAuth.Values[auth.CredsPublicKeyEncoding] != Proto {
		return nil, auth.ErrEncryptionType
	}

	publKeyAddress := toAuth.Values[auth.CredsPublicKeyAddress]
	if len(publKeyAddress) < 1 {
		return nil, errEmptyPublicKeyAddress
	}

	// TODO: use mutex is is.acceptableIDs can be modified using .SetCreds or somehow else
	if is.acceptableIDs != nil && !strlib.In(is.acceptableIDs, publKeyAddress) {
		return nil, nil
	}

	contentToSignature := toAuth.Values[auth.CredsContentToSignature]
	numberToSend, err := strconv.ParseUint(contentToSignature, 10, 64)
	if err != nil {
		return nil, errors.Wrap(auth.ErrSignaturedKey, "not a number!")
	}

	is.mutex.Lock() // Lock() -----------------------------------------------------

	session, ok := is.sessions[numberToSend]
	if !ok {
		return nil, errors.Wrap(auth.ErrSignaturedKey, "no appropriate session")
		is.mutex.Unlock() // Unlock() ---------------------------------------------
	}
	delete(is.sessions, numberToSend)

	is.mutex.Unlock() // Unlock() -------------------------------------------------

	if time.Now().Sub(session.StartedAt) > is.maxSessionDuration {
		return nil, errors.Wrap(auth.ErrAuthSession, "session is expired")
	}

	if session.IP != toAuth.Values[auth.CredsIP] {
		return nil, auth.ErrIP
	}

	signature := []byte(toAuth.Values[auth.CredsSignature])
	publKey := base58.Decode(publKeyAddress)

	if !encrlib.ECDSAVerify(publKey, []byte(contentToSignature), signature) {
		return nil, errWrongSignature
	}

	var nickname = publKeyAddress
	if nicknameReceived := toAuth.Values[auth.CredsNickname]; strings.TrimSpace(nicknameReceived) != "" {
		nickname = nicknameReceived
	}

	return &auth.User{
		Key:      identity.Key(publKeyAddress),
		Nickname: nickname,
		// Creds
	}, nil
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (*identityECDSA) SetCreds(_ *auth.User, toSet auth.Creds) (*auth.Creds, error) {
	// TODO: modify acceptableIDs if it's necessary

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

	publKeyAddress := string(append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...))

	creds := auth.Creds{
		Values: map[auth.CredsType]string{
			auth.CredsPrivateKey:        string(privKeyBytes),
			auth.CredsPublicKeyAddress:  publKeyAddress,
			auth.CredsPublicKeyEncoding: Proto,
			auth.CredsNickname:          toSet.Values[auth.CredsNickname],
			auth.CredsIentityKey:        publKeyAddress,
		},
	}

	return &creds, nil
}

//
//func (*identityECDSA) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsSignature}, nil
//}
