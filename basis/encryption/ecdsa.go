package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"io"
	"math/big"
)

func ECDSASign(privKey ecdsa.PrivateKey, data []byte) ([]byte, error) {
	h := md5.New()
	io.WriteString(h, string(data))

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

func ECDSAVerify(publicKey, dataRaw, signature []byte) bool {
	h := md5.New()
	io.WriteString(h, string(dataRaw))
	data := h.Sum(nil)

	// build key and verify data
	sigLen := len(signature)

	s := big.Int{}
	s.SetBytes([]byte(signature)[(sigLen / 2):])

	r := big.Int{}
	r.SetBytes([]byte(signature)[:(sigLen / 2)])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKey)
	x.SetBytes(publicKey[:(keyLen / 2)])
	y.SetBytes(publicKey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, data, &r, &s)
}

func ECDSASerialize(privKey ecdsa.PrivateKey) ([]byte, error) {

	privKey.Curve = nil

	return json.Marshal(privKey)

	//var encoded bytes.Buffer
	// gob.Register(privKey.Curve)
	//enc := gob.NewEncoder(&encoded)
	//err := enc.Encode(privKey)
	//if err != nil {
	//	return nil, err
	//}
	//return encoded.Bytes(), nil
}

func ECDSADeserialize(data []byte, privKey *ecdsa.PrivateKey) error {
	//decoder := gob.NewDecoder(bytes.NewReader(data))
	//err := decoder.Decode(privKey)
	//if err != nil {
	//	return err
	//}

	err := json.Unmarshal(data, privKey)
	if err != nil {
		return err
	}

	// it's not necessary if privKey.Curve would be serialized instead "nilled"
	privKey.Curve = elliptic.P256()

	return nil
}
