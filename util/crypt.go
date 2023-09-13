package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ethereum/BGService/types"
)

// 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES 加密
func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	return base64.RawURLEncoding.EncodeToString(cryted)
}

// AES 解密
func AesDecrypt(cryted string, key string) string {
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	crytedByte, _ := base64.RawURLEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(types.PublicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RSA 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(types.PrivateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	priv := privInterface.(*rsa.PrivateKey)
	//解析PKCS1格式的私钥
	//priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//if err != nil {
	//	return nil, err
	//}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	//privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	//if err != nil {
	//	return nil, err
	//}
	//// 解密
	//priv := privInterface.(*rsa.PrivateKey)
	//return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	//decrypted := make([]byte, 0, len(ciphertext))
	//for i := 0; i < len(ciphertext); i += 128 {
	//	if i+128 < len(ciphertext) {
	//		partial, err1 := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext[i:])
	//		if err1 != nil {
	//			return []byte(""), err1
	//		}
	//		decrypted = append(decrypted, partial...)
	//	} else {
	//		partial, err1 := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext[i:])
	//		if err1 != nil {
	//			return []byte(""), err1
	//		}
	//		decrypted = append(decrypted, partial...)
	//	}
	//}
	//return decrypted, nil
}
