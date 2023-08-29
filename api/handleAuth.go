package api

import (
	"fmt"
	"github.com/ethereum/api-in/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *ApiService) email(c *gin.Context) {
	email := c.Query("email")
	// 构建电子邮件内容
	to := []string{email}
	subject := "BG verifyCode!"
	verifyCode := util.GenerateVerifyCode(6)
	body := fmt.Sprintf("verifyCode :%s", verifyCode)
	err := util.SendEmail(a.config, to, subject, body)
	if err != nil {
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
