package services

import (
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
)

func CreateAccount() (address string, privateKey string, err error) {
	acc := account.Creation{
		Name:       "test1",
		Passphrase: "",
	}

	if err := account.CreateNewLocalAccount(&acc); err != nil {
		return "", "", err
	}
	addr, _ := store.AddressFromAccountName(acc.Name)

	privateStr, err := ExportPrivateKey(addr, "")
	if err != nil {

	}
	return addr, privateStr, nil
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
