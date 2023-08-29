package rsa0

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	pk, _ := rsa.GenerateKey(rand.Reader, 4096)
	return pk, &pk.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(pk *rsa.PrivateKey) string {
	pkb := x509.MarshalPKCS1PrivateKey(pk)
	pkPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pkb,
		},
	)
	return string(pkPem)
}

func ParseRsaPrivateKeyFromPemStr(pkPem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pkPem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pub, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pub,
		},
	)

	return string(pubPem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func Encrypt(key *rsa.PublicKey, secretMessage string) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, key, []byte(secretMessage), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(privateKey *rsa.PrivateKey, cipherText string) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)

	plaintext, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, privateKey, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func sha3b(msg string) ([]byte, error) {
	byt := []byte(msg)

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha512.New()
	_, err := msgHash.Write(byt)
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}

func Sign(msg string, privateKey *rsa.PrivateKey) ([]byte, error) {
	msgHashSum, err := sha3b(msg)
	if err != nil {
		return nil, err
	}

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)

}

func Verify(msg string, signature []byte, publicKey *rsa.PublicKey) error {
	msgHashSum, err := sha3b(msg)
	if err != nil {
		return err
	}
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(publicKey, crypto.SHA512, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return err
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
	return nil
}
