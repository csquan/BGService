package services

import (
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
)

func CreateAccount() (address string, mnemonic string, err error) {
	acc := account.Creation{
		Name:       "name1",
		Passphrase: "",
	}

	if err := account.CreateNewLocalAccount(&acc); err != nil {
		return "", "", err
	}
	addr, _ := store.AddressFromAccountName(acc.Name)

	return addr, acc.Mnemonic, nil
}

func GetPrivate(address, passphrase string) (string, error) {
	return "", nil
}
