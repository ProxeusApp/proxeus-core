package service

import (
	"github.com/ProxeusApp/proxeus-core/sys/email"
)

type (
	EmailService interface {
		Send(emailTo, subject, body string) error
		SendFrom(emailFrom, emailTo, subject, body string) error
	}
	DefaultEmailService struct {
		emailSender email.EmailSender
	}
)

func NewEmailService(emailS email.EmailSender) EmailService {
	return &DefaultEmailService{emailSender: emailS}
}

//sends an email. The body can contain html tags
func (me *DefaultEmailService) Send(emailTo, subject, body string) error {
	settings, err := settingsDB().Get()
	if err != nil {
		return err
	}

	return me.SendFrom(settings.EmailFrom, emailTo, subject, body)
}

//sends an email. The body can contain html tags
// Additionally its possible to specify the sender of the email with the emailFrom parameter
func (me *DefaultEmailService) SendFrom(emailFrom, emailTo, subject, body string) error {
	mail := &email.Email{From: emailFrom, To: []string{emailTo}, Subject: subject, Body: body}

	return me.emailSender.Send(mail)
}
