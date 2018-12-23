package encryption

import "github.com/btcsuite/btcutil/base58"

func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)

	return []byte(encoded)
}

// Base58Decode decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	decoded := base58.Decode(string(input))

	return decoded
}
