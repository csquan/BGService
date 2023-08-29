package ecc

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

type PrivateKey struct {
	PublicKey
	D []byte
}

func (key PrivateKey) String() string {
	return hex.EncodeToString(key.D)
}

func (key PublicKey) String() string {
	return hex.EncodeToString(elliptic.Marshal(key.Curve, key.X, key.Y))
}

func PublicKeyFromString(public string) (*PublicKey, error) {
	pubByte, err := hex.DecodeString(public)
	if err != nil {
		return nil, err
	}
	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, pubByte)
	if x == nil || y == nil {
		return nil, errors.New("invalid public key")
	}
	return &PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

func KeyFromString(private string) (*PrivateKey, error) {
	d, err := hex.DecodeString(private)
	if err != nil {
		return nil, err
	}
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(d)
	return &PrivateKey{
		PublicKey: PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: d,
	}, nil
}

func GenerateKey() (*PrivateKey, error) {
	curve := elliptic.P256()
	d, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		PublicKey: PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: d,
	}, nil
}

func Encrypt(key crypto.PublicKey, data []byte) (encrypted []byte, err error) {
	if len(data) < 1 {
		err = errors.New("empty data")
		return
	}
	public := key.(*PublicKey)
	if public == nil {
		err = errors.New("invalid public key")
		return
	}
	private, err := GenerateKey()
	if err != nil {
		return
	}
	ephemeral := elliptic.Marshal(private.Curve, private.X, private.Y)
	sym, _ := public.Curve.ScalarMult(public.X, public.Y, private.D)
	// Create buffer
	buf := bytes.Buffer{}
	_, err = buf.Write(sym.Bytes())
	if err != nil {
		return
	}
	_, err = buf.Write([]byte{0x00, 0x00, 0x00, 0x01})
	if err != nil {
		return
	}
	_, err = buf.Write(ephemeral)
	if err != nil {
		return
	}
	hashed := sha256.Sum256(buf.Bytes())
	buf.Reset()
	block, err := aes.NewCipher(hashed[0:16])
	if err != nil {
		return
	}
	ch, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return
	}
	_, err = buf.Write(ephemeral)
	if err != nil {
		return
	}
	_, err = buf.Write(ch.Seal(nil, hashed[16:], data, nil))
	if err != nil {
		return
	}
	encrypted = buf.Bytes()
	return
}

func Decrypt(key crypto.PrivateKey, data []byte) (decrypted []byte, err error) {
	if len(data) < 82 {
		err = errors.New("invalid data size")
		return
	}
	private := key.(*PrivateKey)
	if private == nil {
		err = errors.New("invalid private key")
		return
	}
	curve, buf := elliptic.P256(), bytes.Buffer{}
	x, y := elliptic.Unmarshal(curve, data[0:65])
	sym, _ := curve.ScalarMult(x, y, private.D)
	_, err = buf.Write(sym.Bytes())
	if err != nil {
		return
	}
	_, err = buf.Write([]byte{0x00, 0x00, 0x00, 0x01})
	if err != nil {
		return
	}
	_, err = buf.Write(data[0:65])
	if err != nil {
		return
	}
	hashed := sha256.Sum256(buf.Bytes())
	buf.Reset()

	block, err := aes.NewCipher(hashed[0:16])
	if err != nil {
		return
	}
	ch, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return
	}
	decrypted, err = ch.Open(nil, hashed[16:], data[65:], nil)
	return
}
