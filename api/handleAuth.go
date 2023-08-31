package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-contrib/sessions"
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
	msg := fmt.Sprintf("to: %s, send: %s", email, verifyCode)
	logrus.Info(msg)
	res := util.ResponseMsg(0, "success", "")
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
			if err := db.UpdateUser(a.dbEngine, user.Uid); err != nil {
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
		Uid:                 uid,
		UserName:            username,
		Password:            payload.Password,
		InvitationCode:      inviteCode,
		InvitatedCode:       payload.InviteCode,
		MailBox:             payload.Email,
		ConcernCoinList:     "{}",
		CollectStragetyList: "{}",
		JoinStrageyList:     "{}",
	}
	if err := db.InsertUser(a.dbEngine, &newUser); err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(0, "success", "")
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
		session.Set(payload.Email, has.Uid)
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
	res := util.ResponseMsg(0, "success", "")
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
	res := util.ResponseMsg(0, "success", "")
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
	err = db.UpdateUserPass(a.dbEngine, uidFormatted, user)
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
	err = db.UpdateUserPass(a.dbEngine, uidFormatted, user)
	if err != nil {
		return
	}
	body["status"] = 0
	res := util.ResponseMsg(0, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
