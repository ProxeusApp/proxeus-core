package email

import "git.proxeus.com/core/central/sys/file"

type Email struct {
	From        string
	To          []string
	Subject     string
	Body        string
	Attachments []*file.IO
}

type EmailSender interface {
	Send(e *Email) error
}
