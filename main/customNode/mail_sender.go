package customNode

import (
	"fmt"
	"log"

	"git.proxeus.com/core/central/sys/email"

	"git.proxeus.com/core/central/sys/workflow"
)

type mailSenderNode struct{}

func NewMailSender(n *workflow.Node) (workflow.NodeIF, error) {
	return &mailSenderNode{}, nil
}

func (me mailSenderNode) Execute(n *workflow.Node, dat interface{}) (proceed bool, err error) {
	log.Println("Node mail sender called")
	mailSender, err := email.NewSparkPostEmailSender("")
	if err != nil {
		log.Println("Can't initialize SparkPostEmailSender: " + err.Error())
		return false, err
	}

	mailSender.Send(&email.Email{
		From:    "info@proxeus.com",
		To:      []string{"silvio.rainoldi@blockfactory.com"},
		Subject: "Workflow example connector",
		Body: fmt.Sprintf(
			"Hey, this has been sent from the flow on workflow. CHF/XES: %s", dat.(map[string]interface{})["CHFXES"],
		),
	})

	log.Println("MailSender node's email sent")
	return true, nil
}

func (me mailSenderNode) Remove(n *workflow.Node) {}
func (me mailSenderNode) Close()                  {}
