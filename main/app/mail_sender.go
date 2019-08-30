package app

import (
	"fmt"
	"log"

	"git.proxeus.com/core/central/sys/email"

	"git.proxeus.com/core/central/sys/workflow"
)

type mailSenderNode struct {
	ctx *DocumentFlowInstance
}

func (me *mailSenderNode) Execute(n *workflow.Node) (proceed bool, err error) {
	log.Println("Node mail sender called")
	mailSender, err := email.NewSparkPostEmailSender("24e5ecf598c569a33e9ead68174e679be1c302aa", "")
	if err != nil {
		log.Println("Can't initialize sparkPostEmailSender: " + err.Error())
		return false, err
	}
	chfxes, err := me.ctx.readData("input.CHFXES")
	if err != nil {
		return false, err
	}
	err = mailSender.Send(&email.Email{
		From:    "info@proxeus.com",
		To:      []string{"info@proxeus.com"},
		Subject: "Workflow example connector",
		Body: fmt.Sprintf(
			"Hey, this has been sent from the flow on workflow. CHF/XES: %s", chfxes,
		),
	})
	if err != nil {
		return false, err
	}
	log.Println("MailSender node's email sent")
	return true, nil
}

func (me *mailSenderNode) Remove(n *workflow.Node) {}
func (me *mailSenderNode) Close()                  {}
