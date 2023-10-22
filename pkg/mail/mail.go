package mail

import (
	"gopkg.in/gomail.v2"
)

type MailService struct {
	smtpKey     string
	senderEmail string
	smtpHost    string
	smtpPort    int
	appLink     string
}

type Mail interface {
	SendActivationEmail(to string, subject string, link string) error
}

func New(smtpKey string, sender string, host string, port int, appLink string) *MailService {
	return &MailService{
		smtpKey:     smtpKey,
		senderEmail: sender,
		smtpHost:    host,
		smtpPort:    port,
		appLink:     appLink,
	}
}

func (ms *MailService) sendEmail(to string, subject string, html string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", ms.senderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", html)

	dialer := gomail.NewDialer(ms.smtpHost, ms.smtpPort, ms.senderEmail, ms.smtpKey)
	return dialer.DialAndSend(message)
}

func (ms *MailService) SendActivationEmail(to string, subject string, link string) error {

	t := ms.CreateActivationTemplate(link, to)

	return ms.sendEmail(to, subject, t)
}
