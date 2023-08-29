package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/ethereum/api-in/pkg/conf"
	"github.com/ethereum/api-in/pkg/model"
	"github.com/ethereum/api-in/pkg/util"
	jwt "github.com/golang-jwt/jwt/v4"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, util.Err) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", util.ErrPrvKeyDecode
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", util.ErrPrvKeyDecode
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", util.ErrTokenGen
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (model.TokenUser, util.Err) {
	var tu model.TokenUser

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return tu, util.ErrPubKeyDecode
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return tu, util.ErrPubKeyDecode
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return tu, util.ErrTokenDec
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return tu, util.ErrTokenInvalid
	}

	jsonStr, _ := json.Marshal(claims["sub"])
	_ = json.Unmarshal(jsonStr, &tu)
	return tu, nil
}

func GenAToken(payload interface{}) (string, util.Err) {
	return CreateToken(conf.Conf.AccessTokenExpiresIn, payload, conf.Conf.AccessTokenPrivateKey)
}

func GenRefToken(payload interface{}) (string, util.Err) {
	return CreateToken(conf.Conf.RefreshTokenExpiresIn, payload, conf.Conf.RefreshTokenPrivateKey)
}
