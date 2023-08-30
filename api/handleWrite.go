package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"github.com/ethereum/api-in/db"
	"github.com/ethereum/api-in/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//func (a *ApiService) init(c *gin.Context) {
//	buf := make([]byte, 2048)
//	n, _ := c.Request.Body.Read(buf)
//	data1 := string(buf[0:n])
//	res := types.HttpRes{}
//
//	isValid := gjson.Valid(data1)
//	if isValid == false {
//		logrus.Error("Not valid json")
//		res.Code = http.StatusBadRequest
//		res.Message = "Not valid json"
//		c.SecureJSON(http.StatusBadRequest, res)
//		return
//	}
//	name := gjson.Get(data1, "name")
//	apiKey := gjson.Get(data1, "apiKey")
//	apiSecret := gjson.Get(data1, "apiSecret")
//
//	mechanismData := types.Mechanism{
//		Name:      name.String(),
//		ApiKey:    apiKey.String(),
//		ApiSecret: apiSecret.String(),
//	}
//
//	err := a.db.CommitWithSession(a.db, func(s *xorm.Session) error {
//		if err := a.db.InsertMechanism(s, &mechanismData); err != nil {
//			logrus.Errorf("insert  InsertMechanism task error:%v tasks:[%v]", err, mechanismData)
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	res.Code = 0
//	res.Message = err.Error()
//	res.Data = ""
//
//	c.SecureJSON(http.StatusOK, res)
//}

func (a *ApiService) order(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	data1 := string(buf[0:n])
	res := types.HttpRes{}

	isValid := gjson.Valid(data1)
	if isValid == false {
		logrus.Error("Not valid json")
		res.Code = http.StatusBadRequest
		res.Message = "Not valid json"
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = "null"

	c.SecureJSON(http.StatusOK, res)
}

func (a *ApiService) enroll(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	data1 := string(buf[0:n])
	res := types.HttpRes{}

	isValid := gjson.Valid(data1)
	if isValid == false {
		logrus.Error("Not valid json")
		res.Code = http.StatusBadRequest
		res.Message = "Not valid json"
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	uid := gjson.Get(data1, "uid")
	password := gjson.Get(data1, "password")

	user := types.Users{
		Uid:      uid.String(),
		Password: password.String(),
	}

	db.InsertUser(a.dbEngine, &user)
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = "null"

	c.SecureJSON(http.StatusOK, res)
}

// 引导下载google时调用，产生secret，保存进db
func (a *ApiService) generateSecret(c *gin.Context) {
	secret := GetSecret()
	res := types.HttpRes{}
	user := types.Users{
		Secret: secret,
	}

	db.InsertUser(a.dbEngine, &user)
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = secret

	c.SecureJSON(http.StatusOK, res)
}

// 输入google验证码，确认后触发后端验证
func (a *ApiService) verifyCode(c *gin.Context) {
	uid := c.Param("uid")
	code := c.Param("code")
	res := types.HttpRes{}
	_, secret := db.QuerySecret(a.dbEngine, uid)

	codeint, err := strconv.ParseInt(code, 10, 64)

	if err != nil {

	}

	VerifyCode(secret.Secret, int32(codeint))
	res.Code = 0
	res.Message = "success"
	res.Data = secret

	c.SecureJSON(http.StatusOK, res)
}

// 为了考虑时间误差，判断前当前时间及前后30秒时间
func VerifyCode(secret string, code int32) bool {
	// 当前google值
	if getCode(secret, 0) == code {
		return true
	}

	// 前30秒google值
	if getCode(secret, -30) == code {
		return true
	}

	// 后30秒google值
	if getCode(secret, 30) == code {
		return true
	}

	return false
}

// 获取Google Code
func getCode(secret string, offset int64) int32 {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return int32(oneTimePassword(key, toBytes(epochSeconds/30)))
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func GetSecret() string {
	randomStr := randStr(16)
	return strings.ToUpper(randomStr)
}

func randStr(strSize int) string {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
