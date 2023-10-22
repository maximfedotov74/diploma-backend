package mail

import "gopkg.in/gomail.v2"

type MailService struct {
	smtpKey     string
	senderEmail string
	smtpHost    string
	smtpPort    int
}

func New(smtpKey string, sender string, host string, port int) *MailService {
	return &MailService{
		smtpKey:     smtpKey,
		senderEmail: sender,
		smtpHost:    host,
		smtpPort:    port,
	}
}

func (ms *MailService) SendEmail(to string, subject string, html string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", ms.senderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", html)

	dialer := gomail.NewDialer(ms.smtpHost, ms.smtpPort, ms.senderEmail, ms.smtpKey)
	return dialer.DialAndSend(message)
}
