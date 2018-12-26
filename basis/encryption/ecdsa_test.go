package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

var testDataStr = "test data"

func TestECDSASign(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	signature, err := ECDSASign(*privKey, []byte(testDataStr))
	require.NoError(t, err)

	publKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	fmt.Print("ADDRESS: ", string(Base58Encode(publKey)), "\n\n")

	ok := ECDSAVerify(publKey, []byte(testDataStr), signature)
	require.True(t, ok)

	publKeyBad := append(privKey.PublicKey.Y.Bytes(), privKey.PublicKey.X.Bytes()...)
	ok = ECDSAVerify(publKeyBad, []byte(testDataStr), signature)
	require.False(t, ok)
	ok = ECDSAVerify([]byte{}, []byte(testDataStr), signature)
	require.False(t, ok)

	ok = ECDSAVerify(publKey, []byte(testDataStr+" "), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publKey, []byte(testDataStr[:len(testDataStr)-1]), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publKey, []byte(testDataStr[:len(testDataStr)-1]+" "), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publKey, []byte(""), signature)
	require.False(t, ok)

	ok = ECDSAVerify(publKey, []byte(testDataStr), []byte(string(signature)+" "))
	require.False(t, ok)
	ok = ECDSAVerify(publKey, []byte(testDataStr), signature[:len(testDataStr)-1])
	require.False(t, ok)
	signatureBad := make([]byte, len(signature))
	for i, s := range signature[:len(signature)-1] {
		signatureBad[i] = s
	}
	signatureBad[len(signatureBad)-1] = byte(uint8(signatureBad[len(signatureBad)-1]) + 1%256)

	ok = ECDSAVerify(publKey, []byte(testDataStr), signatureBad)
	require.False(t, ok)
}

func TestECDSASerialize(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	signature, err := ECDSASign(*privKey, []byte(testDataStr))
	require.NoError(t, err)
	require.NotEmpty(t, signature)

	log.Printf("signature 1: %s\n\n", string(signature))

	publKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	ok := ECDSAVerify(publKey, []byte(testDataStr), signature)
	require.True(t, ok)

	serialization, err := ECDSASerialize(*privKey)
	require.NoError(t, err)

	// log.Printf("private key serialization: %s", serialization)

	err = ECDSADeserialize(serialization, nil)
	require.Error(t, err)

	var privKeyRestored ecdsa.PrivateKey
	err = ECDSADeserialize(serialization, &privKeyRestored)
	require.NoError(t, err)

	signatureRestored, err := ECDSASign(privKeyRestored, []byte(testDataStr))
	require.NoError(t, err)

	log.Printf("signature 2: %s\n\n", string(signatureRestored))

	publKeyRestored := append(privKeyRestored.PublicKey.X.Bytes(), privKeyRestored.PublicKey.Y.Bytes()...)
	ok = ECDSAVerify(publKeyRestored, []byte(testDataStr), signature)
	require.True(t, ok)
}
