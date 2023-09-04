package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (a *ApiService) email(c *gin.Context) {
	email := c.Query("email")
	// 构建电子邮件内容
	to := []string{email}
	subject := "BG verifyCode!"
	verifyCode := util.GenerateCode(6)
	sendBody := fmt.Sprintf("verifyCode :%s", verifyCode)
	err := util.SendEmail(a.config, to, subject, sendBody)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err = a.RedisEngine.Set(c, email, verifyCode, 1*time.Minute).Err()
	if err != nil {
		logrus.Error("设置值失败:", err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	msg := fmt.Sprintf("to: %s, send: %s", email, verifyCode)
	logrus.Info(msg)
	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) register(c *gin.Context) {
	var payload *types.UserInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验邮箱是否被注册
	err, has := db.QueryEmail(a.dbEngine, payload.Email)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if has != nil {
		res := util.ResponseMsg(-1, "fail", "Email has already been registered.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验验证码
	//if !util.CheckVerifyCode(c, a.RedisEngine, payload.Email, payload.VerifyCode) {
	//	res := util.ResponseMsg(-1, "fail", "Wrong verifyCode!")
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	// 删除验证码key
	a.RedisEngine.Del(c, payload.Email)
	// 生成8位随机邀请码
	inviteCode := util.GenerateInviteCode(8)
	for {
		err, user := db.QueryInviteCode(a.dbEngine, inviteCode)
		if err != nil {
			// 处理错误
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if user == nil {
			// 邀请码不存在，退出循环
			break
		}
		// 邀请码已存在，生成新的邀请码
		inviteCode = util.GenerateInviteCode(8)
	}
	// uid校验，生成
	uid := util.GenerateCode(14)
	for {
		err, user := db.QuerySecret(a.dbEngine, uid)
		if err != nil {
			// 处理错误
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if user == nil {
			break
		}
		uid = util.GenerateCode(14)
	}

	var username string
	if payload.UserName == "" {
		username = payload.Email
	} else {
		username = payload.UserName
	}
	// 用户填写了邀请码，给邀请码的用户邀请好友数量加1
	if payload.InviteCode != "" {
		err, user := db.QueryInviteCode(a.dbEngine, payload.InviteCode)
		if err != nil {
			// 处理错误
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if user != nil {
			if err := db.AddUserInvite(a.dbEngine, user.Uid); err != nil {
				res := util.ResponseMsg(-1, "fail", err)
				c.SecureJSON(http.StatusOK, res)
				return
			}
		} else {
			res := util.ResponseMsg(-1, "fail", "Incorrect invitation code")
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}
	newUser := types.Users{
		Uid:            uid,
		UserName:       username,
		Password:       payload.Password,
		InvitationCode: inviteCode,
		InvitatedCode:  payload.InviteCode,
		MailBox:        payload.Email,
	}
	if err := db.InsertUser(a.dbEngine, &newUser); err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) login(c *gin.Context) {
	var payload *types.LoginInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var passWord string
	// 获取数据库中的密码
	err, has := db.QueryEmail(a.dbEngine, payload.Email)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", "User does not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if has != nil {
		passWord = has.Password
	}
	if payload.Password == passWord {
		//set session
		session := sessions.Default(c)
		session.Set("Uid", has.Uid)
		session.Set("invitationCode", has.InvitationCode)
		err = session.Save()
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else {
		res := util.ResponseMsg(-1, "fail", "Incorrect password.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})

	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) logout(c *gin.Context) {
	//clear session
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})

	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) forgotPassword(c *gin.Context) {
	uid, _ := c.Get("Uid")
	var payload *types.ForgotPasswordInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验验证码
	if !util.CheckVerifyCode(c, a.RedisEngine, payload.Email, payload.VerifyCode) {
		res := util.ResponseMsg(-1, "fail", "Wrong verifyCode!")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	err, user := db.QuerySecret(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 验证新密码和老密码是否一致
	if user.Password == payload.Password {
		res := util.ResponseMsg(-1, "fail", "New password cannot be the same as the old password")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := map[string]int{
		"status": 1,
	}
	// 是否谷歌验证
	if user.IsBindGoogle {
		res := util.ResponseMsg(1, "fail", body)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 修改密码
	user.Password = payload.Password
	err = db.UpdateUser(a.dbEngine, uidFormatted, user)
	if err != nil {
		return
	}
	body["status"] = 0
	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) resetPassword(c *gin.Context) {
	uid, _ := c.Get("Uid")
	var payload *types.ForgotPasswordInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验验证码
	//if !util.CheckVerifyCode(c, a.RedisEngine, payload.Email, payload.VerifyCode) {
	//	res := util.ResponseMsg(-1, "fail", "Wrong verifyCode!")
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	err, user := db.QuerySecret(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 验证新密码和老密码是否一致
	if user.Password == payload.Password {
		res := util.ResponseMsg(-1, "fail", "New password cannot be the same as the old password")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := map[string]int{
		"status": 1,
	}
	// 修改密码
	user.Password = payload.Password
	err = db.UpdateUser(a.dbEngine, uidFormatted, user)
	if err != nil {
		return
	} else {
		body["status"] = 0
	}
	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 引导下载google时调用，产生secret，保存进db
func (a *ApiService) generateSecret(c *gin.Context) {
	uid := c.Query("uid")
	res := types.HttpRes{}

	//首先查询出这个用户
	user, err := db.GetUser(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("未找到用户记录", uid)

		res.Code = -1
		res.Message = "未找到用户记录"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//产生secret
	user.Secret = GetSecret()

	err = db.UpdateUser(a.dbEngine, uid, user)
	if err != nil {
		logrus.Info("update secret err:", err)

		res.Code = 0
		res.Message = "update secret err"
		res.Data = err

		c.SecureJSON(http.StatusOK, res)
		return
	}
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = user.Secret

	c.SecureJSON(http.StatusOK, res)
	return
}

// 输入google验证码，确认后触发后端验证
func (a *ApiService) verifyCode(c *gin.Context) {
	res := types.HttpRes{}

	var userCode types.UserCodeInfos

	err := c.BindJSON(&userCode)
	if err != nil {
		logrus.Info("传递的不是合法的json参数")

		res.Code = -1
		res.Message = "传递的不是合法的json参数"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	uid := userCode.Uid
	code := userCode.Code

	_, secret := db.QuerySecret(a.dbEngine, uid)

	codeint, err := strconv.ParseInt(code, 10, 64)

	if err != nil {
		logrus.Info("输入的动态码不是合法数字，请检查", code)

		res.Code = -1
		res.Message = "输入的动态码不是合法数字，请检查"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	isTrue := VerifyCode(secret.Secret, int32(codeint))

	res.Code = 0
	if isTrue {
		res.Message = "校验成功"
	} else {
		res.Message = "校验失败"
	}

	c.SecureJSON(http.StatusOK, res)
	return
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
