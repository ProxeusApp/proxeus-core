package email

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	sp "github.com/SparkPost/gosparkpost"
)

type SparkPostEmailSender struct {
	client *sp.Client
}

func NewSparkPostEmailSender(apiKey string) (EmailSender, error) {
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
	return &SparkPostEmailSender{client: client}, nil
}

func (me *SparkPostEmailSender) Send(e *Email) error {
	content := sp.Content{
		From:    e.From,
		Subject: e.Subject,
	}
	//TODO make this if efficient to check for html content
	if body := strings.TrimLeftFunc(e.Body, unicode.IsSpace); strings.HasPrefix(body, "<") {
		content.HTML = e.Body
	} else {
		content.Text = e.Body
	}

	for _, f := range e.Attachments {
		attBytes, err := f.ReadAll()
		if err != nil {
			fmt.Println(err)
			return err
		}
		attach := sp.Attachment{
			MIMEType: f.ContentType(),
			Filename: f.Name(),
			B64Data:  base64.StdEncoding.EncodeToString(attBytes),
		}
		content.Attachments = append(content.Attachments, attach)
	}

	tx := &sp.Transmission{
		Recipients: e.To,
		Content:    content,
	}
	_, _, err := me.client.Send(tx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
