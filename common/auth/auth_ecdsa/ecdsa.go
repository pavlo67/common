package auth_ecdsa

//import (
//	"crypto/ecdsa"
//	"crypto/elliptic"
//	"crypto/rand"
//	r "math/rand"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/btcsuite/btcutil/base58"
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/workshop/common/auth"
//	"github.com/pavlo67/workshop/common/libraries/encrlib"
//	"github.com/pavlo67/workshop/common/libraries/strlib"
//)
//
//const Cryptype encrlib.Cryptype = "ecdsa"
//const Proto = "ecdsa"
//
//var _ auth.Operator = &authECDSA{}
//
//var errWrongSignature = errors.New("wrong signature")
//var errWrongNumber = errors.New("wrong user's number")
//var errEmptyPublicKeyAddress = errors.New("empty public Key address")
//var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")
//
//var cnt uint32
//
//type Session struct {
//	IP        string
//	StartedAt time.Time
//}
//
//type authECDSA struct {
//	sessions map[uint64]Session
//	mutex    *sync.Mutex
//
//	maxSessionDuration time.Duration
//	numbersLimit       int
//
//	acceptableIDs []string
//}
//
//func New(numbersLimit int, maxSessionDuration time.Duration, acceptableIDs []string) (auth.Operator, error) {
//	r.Seed(time.Now().UnixNano())
//
//	is := &authECDSA{
//		sessions: map[uint64]Session{},
//		mutex:    &sync.Mutex{},
//
//		maxSessionDuration: maxSessionDuration,
//		numbersLimit:       numbersLimit,
//
//		acceptableIDs: acceptableIDs,
//	}
//
//	return is, nil
//}
//
//// 	SetCreds creates either session-generated key or new "BTC identity" and returns it
//func (is *authECDSA) SetCreds(userKey auth.Key, creds auth.Creds) (*auth.Creds, error) {
//	toSet := auth.CredsType(creds[auth.CredsToSet])
//
//	if toSet == auth.CredsKeyToSignature {
//		now := time.Now()
//
//		is.mutex.Lock() // Lock() -----------------------------------------------------
//
//		if is.numbersLimit > 0 && len(is.sessions) >= is.numbersLimit {
//			for n, s := range is.sessions {
//				if now.Sub(s.StartedAt) >= is.maxSessionDuration {
//					delete(is.sessions, n)
//				}
//			}
//		}
//
//		cnt++
//		numberToSend := uint64(cnt)<<32 + uint64(r.Uint32())
//
//		is.sessions[numberToSend] = Session{
//			IP:        creds[auth.CredsIP], // TODO??? check if IP isn't empty
//			StartedAt: now,
//		}
//
//		is.mutex.Unlock() // Unlock() -------------------------------------------------
//
//		return &auth.Creds{auth.CredsKeyToSignature: strconv.FormatUint(numberToSend, 10)}, nil
//	}
//
//	// TODO: modify acceptableIDs if it's necessary
//
//	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//	if err != nil {
//		return nil, err
//	} else if privKey == nil {
//		return nil, errEmptyPrivateKeyGenerated
//	}
//
//	privKeyBytes, err := encrlib.ECDSASerialize(*privKey)
//	if err != nil {
//		return nil, err
//	}
//
//	publKeyBase58 := base58.Encode(encrlib.ECDSAPublicKey(*privKey))
//	nickname := publKeyBase58
//	if creds[auth.CredsNickname] != "" {
//		nickname = creds[auth.CredsNickname]
//	}
//
//	credsNew := &auth.Creds{
//		auth.CredsNickname:          nickname,
//		auth.CredsPrivateKey:        string(privKeyBytes),
//		auth.CredsPublicKeyBase58:   publKeyBase58,
//		auth.CredsPublicKeyEncoding: Proto,
//	}
//
//	return credsNew, nil
//}
//
//const onAuthorize = "on authECDSA.Authenticate(): "
//
//func (is *authECDSA) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
//	if toAuth[auth.CredsPublicKeyEncoding] != Proto {
//		return nil, errors.Wrap(auth.ErrEncryptionType, onAuthorize)
//	}
//
//	publKeyBase58 := toAuth[auth.CredsPublicKeyBase58]
//	if len(publKeyBase58) < 1 {
//		return nil, errEmptyPublicKeyAddress
//	}
//	publKey := base58.Decode(publKeyBase58)
//
//	// TODO: use mutex is is.acceptableIDs can be modified using .SetCreds or somehow else
//	if is.acceptableIDs != nil && !strlib.In(is.acceptableIDs, publKeyBase58) {
//		return nil, nil
//	}
//
//	keyToSignature := toAuth[auth.CredsKeyToSignature]
//	numberToSend, err := strconv.ParseUint(keyToSignature, 10, 64)
//	if err != nil {
//		return nil, errors.Wrap(auth.ErrSignaturedKey, "not a number!")
//	}
//
//	is.mutex.Lock() // Lock() -----------------------------------------------------
//
//	session, ok := is.sessions[numberToSend]
//	if !ok {
//		return nil, errors.Wrap(auth.ErrSignaturedKey, "no appropriate session")
//		is.mutex.Unlock() // Unlock() ---------------------------------------------
//	}
//	delete(is.sessions, numberToSend)
//
//	is.mutex.Unlock() // Unlock() -------------------------------------------------
//
//	if time.Now().Sub(session.StartedAt) > is.maxSessionDuration {
//		return nil, errors.Wrap(auth.ErrAuthSession, "session is expired")
//	}
//
//	if session.IP != toAuth[auth.CredsIP] {
//		return nil, auth.ErrIP
//	}
//
//	signature := []byte(toAuth[auth.CredsSignature])
//
//	if !encrlib.ECDSAVerify(keyToSignature, publKey, signature) {
//		return nil, errWrongSignature
//	}
//
//	var nickname = publKeyBase58
//	if nicknameReceived := toAuth[auth.CredsNickname]; strings.TrimSpace(nicknameReceived) != "" {
//		nickname = nicknameReceived
//	}
//
//	return &auth.Identity{
//		Key:   auth.Key(Proto + "://" + publKeyBase58),
//		Creds: auth.Creds{auth.CredsNickname: nickname},
//	}, nil
//}
//
////
////func (*authECDSA) Accepts() ([]auth.CredsType, error) {
////	return []auth.CredsType{auth.CredsSignature}, nil
////}
