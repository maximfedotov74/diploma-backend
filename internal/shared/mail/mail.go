package mail

import (
	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	SmtpKey     string
	SenderEmail string
	SmtpHost    string
	SmtpPort    int
	AppLink     string
}

type MailService struct {
	config MailConfig
}

func NewMailService(config MailConfig) *MailService {
	return &MailService{
		config: config,
	}
}

func (ms *MailService) sendEmail(to string, subject string, html string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", ms.config.SenderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", html)

	dialer := gomail.NewDialer(ms.config.SmtpHost, ms.config.SmtpPort, ms.config.SenderEmail, ms.config.SmtpKey)
	return dialer.DialAndSend(message)
}

func (ms *MailService) SendActivationEmail(to string, subject string, link string) error {

	t := ms.createActivationTemplate(link, to)

	return ms.sendEmail(to, subject, t)
}

func (ms *MailService) SendChangePasswordEmail(to string, subject string, code string) error {
	t := ms.createChangePasswordCodeTemplate(code, to)

	return ms.sendEmail(to, subject, t)
}

func (ms *MailService) SendOrderActivationEmail(to string, subject string, link string) error {
	t := ms.createOrderActivationTemplate(link, to)

	return ms.sendEmail(to, subject, t)
}
