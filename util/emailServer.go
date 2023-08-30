package util

import (
	"fmt"
	"github.com/ethereum/BGService/config"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(config *config.Config, to []string, subject string, body string) error {
	// 发送电子邮件
	// 配置SMTP服务器信息
	smtpServer := config.Email.SmtpServer
	smtpPort := config.Email.SmtpPort
	smtpUsername := config.Email.SmtpUsername
	smtpPassword := config.Email.SmtpPassword

	e := email.NewEmail()
	e.From = smtpUsername
	e.To = to
	e.Subject = subject
	e.Text = []byte(body)

	// 连接到SMTP服务器
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := e.Send(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth)
	return err
}
