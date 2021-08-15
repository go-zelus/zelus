package email

import (
	"crypto/tls"

	"github.com/go-zelus/zelus/core/logx"

	"github.com/go-zelus/zelus/core/config"

	"gopkg.in/gomail.v2"
)

var instance *Email

type Email struct {
	*SMTPInfo
}

type SMTPInfo struct {
	Host     string   `mapstructure:"host"`
	Port     int      `mapstructure:"port"`
	IsSSL    bool     `mapstructure:"is_ssl"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	From     string   `mapstructure:"from"`
	To       []string `mapstructure:"to"`
}

func Init() {
	New()
}

func New(info ...*SMTPInfo) *Email {
	if instance == nil {
		if len(info) > 0 {
			instance = &Email{SMTPInfo: info[0]}
		} else {
			conf := &SMTPInfo{}
			err := config.UnmarshalKey("email", conf)
			if err != nil {
				logx.Panic("email 配置解析错误")
			}
			instance = &Email{SMTPInfo: conf}
		}
	}
	return instance
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
