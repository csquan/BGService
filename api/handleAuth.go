package api

import (
	"fmt"
	"github.com/ethereum/api-in/db"
	"github.com/ethereum/api-in/types"
	"github.com/ethereum/api-in/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (a *ApiService) email(c *gin.Context) {
	email := c.Query("email")
	// 构建电子邮件内容
	to := []string{email}
	subject := "BG verifyCode!"
	verifyCode := util.GenerateCode(6)
	body := fmt.Sprintf("verifyCode :%s", verifyCode)
	err := util.SendEmail(a.config, to, subject, body)
	if err != nil {
		res := util.ResponseMsg(0, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err = a.RedisEngine.Set(c, email, verifyCode, 1*time.Minute).Err()
	if err != nil {
		logrus.Error("设置值失败:", err)
		res := util.ResponseMsg(0, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	msg := fmt.Sprintf("to: %s, send: %s", email, verifyCode)
	logrus.Info(msg)
	res := util.ResponseMsg(1, "success", "")
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) register(c *gin.Context) {
	var payload *types.UserInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		res := util.ResponseMsg(0, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验邮箱是否被注册
	user := db.QueryEmail(a.dbEngine, payload.Email)
	if user != nil {
		res := util.ResponseMsg(0, "fail", "Email has already been registered.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验验证码
	if !util.CheckVerifyCode(c, a, payload.Email, payload.VerifyCode) {
		res := util.ResponseMsg(0, "fail", "Wrong code!")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 删除验证码key
	a.RedisEngine.Del(c, payload.Email)
	// 生成14位随机邀请码
	inviteCode := util.GenerateInviteCode(8)
	for db.QueryInviteCode(a.dbEngine, inviteCode) != nil {
		inviteCode = util.GenerateInviteCode(8)
	}
	// uid校验，生成
	uid := util.GenerateCode(14)
	for db.QuerySecret(a.dbEngine, uid) != nil {
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
		user := db.QueryInviteCode(a.dbEngine, payload.InviteCode)
		if db.QuerySecret(a.dbEngine, user.Uid) != nil {
			if err := db.UpdateUser(a.dbEngine, user.Uid); err != nil {
				res := util.ResponseMsg(0, "fail", err)
				c.SecureJSON(http.StatusOK, res)
				return
			}
		} else {
			res := util.ResponseMsg(0, "fail", "Incorrect invitation code")
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
		res := util.ResponseMsg(0, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(1, "success", "")
	c.SecureJSON(http.StatusOK, res)
	return
}
