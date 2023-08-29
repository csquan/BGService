package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/ethereum/api-in/types"
	"net/http"
	"time"
)

//const PubPem = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEyN1hMUTBxYXJJUkgxR1hHaUE3MQpiZFJBZ3M4b1BlQUJST1FhVE1OMlZtUjJkY1VpbklnYUp4WVc0V1U0NkxxVFpWYzlGb1hpMmNtV2pyRWovWGRoCjRKaXFEdlhuN1JNdjgveXEwczZmUHFTc3BHUkxldEVYV0tzbUc2WEFycFZaWVQ1bDNZZEJKRXpjZHQ5MElTWUgKc2tBWjY2M0kzTGJFYzBSZks2cW14bzJNbFZyQ1FNTW5UTURQbHlwQVkyQU9qSHNsTW0ybVdjRWg4WVlYRm52LwovQUxNUndRQ3AwNXkwRjB1YXFrOGxiLy9aVjdaanpVWWRDMkZ4WW16MEdXakx3eDFQWk1qY2loVEp6dWJ6OTFuCkVSakdXQWlOc21jYzJ1L1BPczZCZ3V0NlBDb2VPYWxybE54RFl3R3U2Z05jbjZZMVo1elg4VTh6ZmF1TDVKR0MKUTV3U2JjK1N6eVhvT05wQVU1Nm92eW00Umsxdk9JdTNTN2V4MEFvV2lmTmNaamlONkZ1YjNya1k3ckE1WWg1TQpaWTBiRUNoQm5mcnlTeGpITU4xWks5V05ud2FsUUNBeHhkNFh6SWkxNnRlU2FWZC9hMFVIa1BlejZuWExzeXB1CmUwN1R3U2FlNzVhWWFPci9sWDFudmZlWlpyWkdMOUwzK2Naa3Y0UUkvMUJ6cGtYZ2ZZelVBdWQzdEhaNDJyZkYKMmRnTDZ4eFpaVUtKVmtNMC9lY1RHanRmZXo2K3o4UnM0VVQ2dWUwa3cxNUJDZ0JhZzJrdHAvQjFOQ3lhdWRhVQpDZERXemdNbU9KcDBNSFVUK3FIT0thc1pPU1dKYUxQOUd5ODZTSWVEWm44bnRRdGFUWUVjU3ZBVDMyMWZraHR1ClB0SXFrUGZwOS9nWlJZN1R6dnVkQjMwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo="

func nonce() string {
	result := make([]byte, 12)
	_, _ = rand.Read(result)
	return hex.EncodeToString(result)
}

// PubPem就是ApiSecret
func enc(nonce string, PubPem string) string {
	pemByt, _ := base64.StdEncoding.DecodeString(PubPem)
	block, _ := pem.Decode(pemByt)
	if block == nil {
		panic("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("key type is not RSA")
	}

	nowStr := time.Now().UTC().Format(http.TimeFormat)

	ciphertext, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pubKey, []byte(nowStr+nonce), nil)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

//type ReqAccount struct {
//	Handshake
//	AccountId string `json:"AccountId"`
//}

func LiveHandshake(key string, secret string) types.Handshake {
	nonce := nonce()
	verified := enc(nonce, secret)
	return types.Handshake{
		ApiKey:   key,
		Nonce:    nonce,
		Verified: verified,
	}
}

//func post() {
//	handshake := LiveHandshake()
//	data := ReqAccount{
//		Handshake: handshake,
//		AccountId: "b1JuaPlPaImVOSD",
//	}
//	jsonValue, _ := json.Marshal(data)
//	r, err := http.Post("https://api.huiwang.io/api/v1/account/query", "application/json", bytes.NewBuffer(jsonValue))
//	if err != nil {
//		panic(err)
//	}
//	body, _ := io.ReadAll(r.Body)
//	fmt.Println("Post result:", string(body))
//}

//
//func get() {
//	handshake := LiveHandshake()
//	jsonValue, _ := json.Marshal(handshake)
//	get, err := http.NewRequest("GET", "https://api.huiwang.io/api/v1/account/query?accountId=b1JuaPlPaImVOSD", bytes.NewBuffer(jsonValue))
//	if err != nil {
//		panic(err)
//	}
//	r, err := http.DefaultClient.Do(get)
//	body, _ := io.ReadAll(r.Body)
//	fmt.Println("Get result:", string(body))
//}
