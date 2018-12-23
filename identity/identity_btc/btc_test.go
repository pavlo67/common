package identity_btc

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/pavlo67/punctum/basis/encryption"
	"golang.org/x/crypto/ripemd160"
)

// Base58Encode encodes a byte array to Base58

func GetAddress(publicKey []byte) []byte {
	pubKeyHash, _ := HashPubKey(publicKey)

	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := encryption.Base58Encode(fullPayload)

	return address
}

// ValidateAddress check if address is valid, has valid format
func ValidateAddress(address string) bool {
	if len(address) == 0 {
		return false
	}

	pubKeyHash := encryption.Base58Decode([]byte(address))

	if len(pubKeyHash) <= AddressChecksumLen {
		return false
	}
	actualChecksum := pubKeyHash[len(pubKeyHash)-AddressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLen]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

const Version = byte(0x00)
const AddressChecksumLen = 4

// Converts hash of pubkey to address as a string
func PubKeyHashToAddres(pubKeyHash []byte) (string, error) {
	versionedPayload := append([]byte{Version}, pubKeyHash...)

	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := encryption.Base58Encode(fullPayload)

	return fmt.Sprintf("%s", address), nil
}

// Makes string adres from pub key
func PubKeyToAddres(pubKey []byte) (string, error) {
	pubKeyHash, err := HashPubKey(pubKey)

	if err != nil {
		return "", err
	}
	versionedPayload := append([]byte{Version}, pubKeyHash...)

	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := encryption.Base58Encode(fullPayload)

	return fmt.Sprintf("%s", address), nil
}

// Checksum generates a checksum for a public key
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, err
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
}
