package ecies

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// HexToECDSA parses a P256 private key.
func HexToECDSA(hexKey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexKey)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return ToECDSA(b)
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	pk.PublicKey.Curve = elliptic.P256()
	if strict && 8*len(d) != pk.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", pk.Params().BitSize)
	}
	pk.D = new(big.Int).SetBytes(d)

	// The pk.D must < N
	if pk.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The pk.D must not be zero or negative.
	if pk.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(d)
	if pk.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return pk, nil
}
