package email

import (
	"crypto/tls"
	"strconv"
	"strings"

	"github.com/go-gomail/gomail"
)

/******************************************************************************
 **函数名称: SendMail
 **功    能: 发送邮件
 **输入参数:
 **     from: 发送邮箱
 **     passwd: 邮箱密码
 **     to: 接收邮箱
 **     subject: 邮件主题
 **     body: 邮件内容
 **输出参数:
 **返    回: 错误描述
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.08.10 08:23:39 #
 ******************************************************************************/
func SendMail(from string, passwd string,
	host string, to string, subject string, body string) error {
	hp := strings.Split(host, ":")
	port, _ := strconv.ParseInt(hp[1], 10, 32)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, from)
	m.SetHeader("To", strings.Split(to, ",")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(hp[0], int(port), from, passwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
