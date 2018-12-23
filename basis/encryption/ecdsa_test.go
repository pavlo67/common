package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDataStr = "test data"

func TestECDSASign(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privateKey)

	signature, err := ECDSASign(*privateKey, []byte(testDataStr))
	require.NoError(t, err)

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	ok := ECDSAVerify(publicKey, []byte(testDataStr), signature)
	require.True(t, ok)

	publicKeyBad := append(privateKey.PublicKey.Y.Bytes(), privateKey.PublicKey.X.Bytes()...)
	ok = ECDSAVerify(publicKeyBad, []byte(testDataStr), signature)
	require.False(t, ok)
	ok = ECDSAVerify([]byte{}, []byte(testDataStr), signature)
	require.False(t, ok)

	ok = ECDSAVerify(publicKey, []byte(testDataStr+" "), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publicKey, []byte(testDataStr[:len(testDataStr)-1]), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publicKey, []byte(testDataStr[:len(testDataStr)-1]+" "), signature)
	require.False(t, ok)
	ok = ECDSAVerify(publicKey, []byte(""), signature)
	require.False(t, ok)

	ok = ECDSAVerify(publicKey, []byte(testDataStr), []byte(string(signature)+" "))
	require.False(t, ok)
	ok = ECDSAVerify(publicKey, []byte(testDataStr), signature[:len(testDataStr)-1])
	require.False(t, ok)
	signatureBad := make([]byte, len(signature))
	for i, s := range signature[:len(signature)-1] {
		signatureBad[i] = s
	}
	signatureBad[len(signatureBad)-1] = byte(uint8(signatureBad[len(signatureBad)-1]) + 1%256)

	ok = ECDSAVerify(publicKey, []byte(testDataStr), signatureBad)
	require.False(t, ok)
}

func TestECDSASerialize(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privateKey)

	signature, err := ECDSASign(*privateKey, []byte(testDataStr))
	require.NoError(t, err)
	require.NotEmpty(t, signature)

	log.Printf("signature 1: %s\n\n", string(signature))

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	ok := ECDSAVerify(publicKey, []byte(testDataStr), signature)
	require.True(t, ok)

	serialization, err := ECDSASerialize(*privateKey)
	require.NoError(t, err)

	// log.Printf("private key serialization: %s", serialization)

	err = ECDSADeserialize(serialization, nil)
	require.Error(t, err)

	var privateKeyRestored ecdsa.PrivateKey
	err = ECDSADeserialize(serialization, &privateKeyRestored)
	require.NoError(t, err)

	signatureRestored, err := ECDSASign(privateKeyRestored, []byte(testDataStr))
	require.NoError(t, err)

	log.Printf("signature 2: %s\n\n", string(signatureRestored))

	publicKeyRestored := append(privateKeyRestored.PublicKey.X.Bytes(), privateKeyRestored.PublicKey.Y.Bytes()...)
	ok = ECDSAVerify(publicKeyRestored, []byte(testDataStr), signature)
	require.True(t, ok)
}
