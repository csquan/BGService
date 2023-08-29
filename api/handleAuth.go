package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/smtp"
	"strings"
)

func (a *ApiService) email(c *gin.Context) {
	email := c.Query("email")
	smtpServer := a.config.Email.SmtpServer
	smtpPort := a.config.Email.SmtpPort
	smtpUsername := a.config.Email.SmtpUsername
	smtpPassword := a.config.Email.SmtpPassword
	// 构建电子邮件内容
	to := []string{email}
	subject := "Hello from Go!"
	body := "This is a test email sent from Go."

	message := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// 连接到SMTP服务器
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, smtpUsername, to, message)
	if err != nil {
		fmt.Println("邮件发送失败:", err)
		return
	}

	fmt.Println("邮件发送成功")
}
