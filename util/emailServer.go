package util

import (
	"github.com/ethereum/BGService/config"
	"gopkg.in/gomail.v2"
	"log"
)

func SendEmail(config *config.Config, to []string, subject string, body string) error {
	// 发送电子邮件
	// 配置SMTP服务器信息
	mailTitle := subject //邮件标题
	mailBody := body     //邮件内容,可以是html

	//接收者邮箱列表
	mailTo := to

	m := gomail.NewMessage()
	m.SetHeader("From", config.Email.SmtpUsername) //发送者腾讯企业邮箱账号
	m.SetHeader("To", mailTo...)                   //接收者邮箱列表
	m.SetHeader("Subject", mailTitle)              //邮件标题
	m.SetBody("text/html", mailBody)               //邮件内容,可以是html

	//发送邮件服务器、端口、发件人账号、发件人密码
	//服务器地址和端口是腾讯的
	d := gomail.NewDialer(config.Email.SmtpServer, config.Email.SmtpPort, config.Email.SmtpUsername, config.Email.SmtpPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Println("send mail failed", err)
		return err
	}

	log.Println("success")

	return nil
}
