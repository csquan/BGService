package services

import (
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
)

func CreateAccount() (address string, err error) {
	acc := account.Creation{
		Name:       "name1",
		Passphrase: "",
	}

	if err := account.CreateNewLocalAccount(&acc); err != nil {
		return "", err
	}
	addr, _ := store.AddressFromAccountName(acc.Name)
	//这里不返回私钥，直接存储db--todo：后期密文存储
	return addr, nil
}

func GetPrivate(address, passphrase string) (string, error) {
	return "", nil
}
