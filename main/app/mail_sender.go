package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/ProxeusApp/proxeus-core/sys/email"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

type mailSenderNode struct {
	ctx  *DocumentFlowInstance
	ctx2 *ExecuteAtOnceContext
}

func (me *mailSenderNode) Execute(n *workflow.Node) (proceed bool, err error) {
	log.Println("Node mail sender called")
	mailSender, err := email.NewSparkPostEmailSender("24e5ecf598c569a33e9ead68174e679be1c302aa", "")
	if err != nil {
		log.Println("Can't initialize sparkPostEmailSender: " + err.Error())
		return false, err
	}
	var chfxes interface{}
	if me.ctx != nil {
		chfxes, err = me.ctx.readData("input.CHFXES")
		if err != nil {
			return false, err
		}
	} else {
		chfxes = me.ctx2.data["CHFXES"]
		if chfxes == nil {
			return false, errors.New("no data")
		}
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
