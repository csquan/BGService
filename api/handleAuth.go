package api

import (
	"fmt"
	"github.com/ethereum/api-in/types"
	"github.com/ethereum/api-in/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *ApiService) email(c *gin.Context) {
	email := c.Query("email")
	fmt.Println(email)
	// 构建电子邮件内容
	to := []string{email}
	subject := "BG verifyCode!"
	verifyCode := util.GenerateVerifyCode(6)
	body := fmt.Sprintf("verifyCode :%s", verifyCode)
	err := util.SendEmail(a.config, to, subject, body)
	if err != nil {
		return
	}
	logrus.Info("邮件发送成功", verifyCode)
	res := types.HttpRes{}
	res.Code = 1
	res.Message = "success"
	res.Data = ""
	c.SecureJSON(http.StatusOK, res)
	return
}
