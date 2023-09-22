package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/services"
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
	subject := "BG 認証コード"
	verifyCode := util.GenerateCode(6)
	sendBody := fmt.Sprintf("ユーザー %s 様<br/> こんにちは！<br/> BG をご利用いただきありがとうございます。認証コードを入力して認証を完了してください。<br/> 認証コード：%s<br/> 有効期限は 3 分間です。公開しないでください。", email, verifyCode)
	err := util.SendEmail(a.config, to, subject, sendBody)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err = a.RedisEngine.Set(c, email, verifyCode, 3*time.Minute).Err()
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
	if !util.CheckVerifyCode(c, a.RedisEngine, payload.Email, payload.VerifyCode) {
		res := util.ResponseMsg(-1, "fail", "Wrong verifyCode!")
		c.SecureJSON(http.StatusOK, res)
		return
	}
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
	// 用户填写了邀请码，加入邀请等级表
	if payload.InviteCode != "" {
		// 根据邀请码查邀请人用户信息
		err, user := db.QueryInviteCode(a.dbEngine, payload.InviteCode)
		if err != nil {
			// 处理错误
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if user == nil {
			res := util.ResponseMsg(-1, "fail", "Incorrect invitation code")
			c.SecureJSON(http.StatusOK, res)
			return
		}
		// 根据邀请人的信息查邀请人是否有上级邀请(处理二级邀请)
		err, Invite := db.QueryInvite(a.dbEngine, user.Uid)
		if err != nil {
			// 处理错误
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if Invite != nil {
			newSecondInvitation := types.Invitation{
				Uid:    Invite.Uid,
				SonUid: uid,
				Level:  "2",
			}
			err = db.InsertInvitation(a.dbEngine, &newSecondInvitation)
			if err != nil {
				res := util.ResponseMsg(-1, "fail", err)
				c.SecureJSON(http.StatusOK, res)
				return
			}
		}
		newInvitation := types.Invitation{
			Uid:    user.Uid,
			SonUid: uid,
			Level:  "1",
		}
		err = db.InsertInvitation(a.dbEngine, &newInvitation)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}
	newUser := types.Users{
		Uid:                 uid,
		UserName:            username,
		Password:            payload.Password,
		InvitationCode:      inviteCode,
		MailBox:             payload.Email,
		ConcernCoinList:     "{}",
		CollectStragetyList: "{}",
	}
	//这个下面得用事务 1.插入用户表 2.插入链上地址表
	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}
	if _, err := session.Insert(newUser); err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	//下面私钥加密存储
	addr, privateKey, name, err := services.CreateAccount()
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	priEncrypt := util.AesEncrypt(privateKey, types.AesKey)

	userKey := types.UserKey{
		Addr:       addr,
		Name:       name,
		PrivateKey: priEncrypt,
	}
	_, err = session.Table("userKey").Insert(userKey)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	userAddr := types.UserAddr{
		Uid:     uid,
		Network: "TRX",
		Addr:    addr,
	}

	_, err = session.Table("userAddr").Insert(userAddr)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}
	err = session.Commit()
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
		logrus.Error(err)
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
	// 用户未登录状态
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
	err, user := db.QueryEmail(a.dbEngine, payload.Email)
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
	body := make(map[string]interface{})
	// 是否谷歌验证
	if user.IsBindGoogle {
		body["status"] = 1
	} else {
		body["status"] = 0
		// 修改密码
		user.Password = payload.Password
		err = db.UpdateUser(a.dbEngine, user.Uid, user)
		if err != nil {
			return
		}
	}
	body["uid"] = user.Uid
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
		logrus.Info(err)
		res = util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("can not find user:", uid)

		res = util.ResponseMsg(-1, "can not find user", uid)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//产生secret
	user.Secret = GetSecret()

	err = db.UpdateUser(a.dbEngine, uid, user)
	if err != nil {
		res = util.ResponseMsg(-1, "update secret err", err)
		c.SecureJSON(http.StatusOK, res)
		return

	}
	//下面将信息存入db
	res = util.ResponseMsg(0, "success", user.Secret)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 输入google验证码，确认后触发后端验证
func (a *ApiService) verifyCode(c *gin.Context) {
	res := types.HttpRes{}

	var userCode types.UserCodeInfos

	err := c.BindJSON(&userCode)
	if err != nil {
		res = util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	uid := userCode.Uid
	code := userCode.Code

	_, secret := db.QuerySecret(a.dbEngine, uid)

	codeint, err := strconv.ParseInt(code, 10, 64)

	if err != nil {
		res = util.ResponseMsg(-1, "fail", err)
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
