package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
)

func randStringBytesCrypto(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func CreateAccount() (address string, privateKey string, name string, err error) {
	name, err = randStringBytesCrypto(5)
	if err != nil {
		return "", "", "", err
	}

	acc := account.Creation{
		Name:       name,
		Passphrase: "",
	}

	if err := account.CreateNewLocalAccount(&acc); err != nil {
		return "", "", "", err
	}
	addr, _ := store.AddressFromAccountName(name)

	privateStr, err := ExportPrivateKey(addr, "")
	if err != nil {

	}
	return addr, privateStr, name, nil
}

// ExportPrivateKey from account
func ExportPrivateKey(address, passphrase string) (string, error) {
	ks := store.FromAddress(address)
	allAccounts := ks.Accounts()
	for _, account := range allAccounts {
		if account.Address.String() == address {
			_, key, err := ks.GetDecryptedKey(keystore.Account{Address: account.Address}, passphrase)
			if err != nil {
				return "", err
			}
			fmt.Printf("%064x\n", key.PrivateKey.D)
			str := fmt.Sprintf("%064x\n", key.PrivateKey.D)
			return str, nil
		}
	}
	return "", nil
}

func GetPrivate(address, passphrase string) (string, error) {
	return "", nil
}
