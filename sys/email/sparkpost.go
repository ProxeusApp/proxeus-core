package email

import (
	"errors"
	"strings"
	"unicode"

	sp "github.com/SparkPost/gosparkpost"
)

type sparkPostEmailSender struct {
	client           *sp.Client
	defaultEmailFrom string
}

func NewSparkPostEmailSender(apiKey, defaultEmailFrom string) (EmailSender, error) {
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}

	client := &sp.Client{}
	err := client.Init(cfg)
	if err != nil {
		return nil, err
	}
	return &sparkPostEmailSender{client: client, defaultEmailFrom: defaultEmailFrom}, nil
}

func (me *sparkPostEmailSender) Send(e *Email) error {
	emailFrom := e.From
	if len(emailFrom) == 0 {
		emailFrom = me.defaultEmailFrom
	}
	if len(emailFrom) == 0 {
		return errors.New("the From attribute has to be populated to send an email")
	}
	content := sp.Content{
		From:    emailFrom,
		Subject: e.Subject,
	}
	//TODO make this if efficient to check for html content
	if body := strings.TrimLeftFunc(e.Body, unicode.IsSpace); strings.HasPrefix(body, "<") {
		content.HTML = e.Body
	} else {
		content.Text = e.Body
	}

	tx := &sp.Transmission{
		Recipients: e.To,
		Content:    content,
	}
	_, _, err := me.client.Send(tx)
	if err != nil {
		return err
	}
	return err
}
