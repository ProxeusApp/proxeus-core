package email

type Email struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type EmailSender interface {
	Send(e *Email) error
}
