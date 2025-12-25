package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Sender struct {
	host string
	port string
	user string
	pass string
	from string
	addr string
	auth smtp.Auth
}

func NewSMTPSender(
	host string,
	port string,
	user string,
	pass string,
	from string,
) *Sender {
	addr := fmt.Sprintf("%s:%s", host, port)

	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	return &Sender{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
		addr: addr,
		auth: auth,
	}
}

func (s *Sender) Send(to, subject, body string) error {
	headers := map[string]string{
		"From":         s.from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=UTF-8",
	}

	var msg strings.Builder

	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(
		s.addr,
		s.auth,
		s.from,
		[]string{to},
		[]byte(msg.String()),
	)
}
