package service

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"os"
)

var mailer *gomail.Dialer

func ConnectMailer() {
	mailer = gomail.NewDialer(
		os.Getenv("MAILER_HOST"),
		ParseToInt(os.Getenv("MAILER_PORT")),
		os.Getenv("MAILER_EMAIL"),
		os.Getenv("MAILER_PASSWORD"),
	)

	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}

func SendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAILER_EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return mailer.DialAndSend(m)
}
