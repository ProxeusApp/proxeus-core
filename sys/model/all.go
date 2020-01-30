package model

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ProxeusApp/proxeus-core/sys/model/compatability"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

const (
	PaymentStatusCreated   = "created"
	PaymentStatusPending   = "pending"
	PaymentStatusConfirmed = "confirmed"
	PaymentStatusCancelled = "cancelled"
	PaymentStatusRedeemed  = "redeemed"
	PaymentStatusDeleted   = "deleted"
	PaymentStatusTimeout   = "timeout"
)

type (
	FormItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`

		//Data contains the form source
		Data compatability.CarriedStringMap `json:"data"`
	}

	UserDataItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		WorkflowID string    `json:"workflowID"`
		Finished   bool      `json:"finished"`
		Updated    time.Time `json:"updated" storm:"index"`
		Created    time.Time `json:"created" storm:"index"`

		Data compatability.CarriedStringMap `json:"data"`
	}

	FormComponentItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`

		//Settings and Template are just passing through
		Settings compatability.CarriedJsonRaw `json:"settings"`
		Template string                       `json:"template"`
	}

	WorkflowItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`
		Price   uint64    `json:"price" storm:"index"`

		Data            *workflow.Workflow `json:"data"`
		OwnerEthAddress string             `json:"ownerEthAddress"` //only used in frontend
		Deactivated     bool               `json:"deactivated"`
	}

	TemplateItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`

		Data TemplateLangMap `json:"data"`
	}

	SignatureRequestItem struct {
		ID string `storm:"id"`
		// DocumentID, ie. 8b5e1460-e456-4aea-91dd-efd7f6ff9b40
		DocId string `json:"docid" storm:"index"`
		// DataPath, ie. docs[0]
		DocPath string `json:"docpath" storm:"index"`
		// Keccak256 of Filecontent
		Hash string `json:"hash"`
		// RequestorEthAddress
		Requestor string `json:"requestor" storm:"index"`
		// Timestamp for requested, revoked, rejected
		RequestedAt time.Time `json:"requestedAt"`
		RevokedAt   time.Time `json:"revokedAt"`
		RejectedAt  time.Time `json:"rejectedAt"`
		// SignatoryEthAddress
		Signatory string `json:"signatory" storm:"index"`
		Revoked   bool   `json:"revoked"`
		Rejected  bool   `json:"rejected"`
	}

	WorkflowPaymentItem struct {
		//save from who payment
		ID         string    `json:"id" storm:"id,unique"`
		TxHash     string    `json:"hash" storm:"index,unique"`
		WorkflowID string    `json:"workflowID" storm:"index"`
		From       string    `json:"from" storm:"index"`
		To         string    `json:"to"`
		Xes        uint64    `json:"xes"`
		Status     string    `json:"status" storm:"index"`
		CreatedAt  time.Time `json:"createdAt"`
	}

	TemplateLangMap map[string]*file.IO

	TokenRequest struct {
		Email  string    `json:"email" validate:"email=true,required=true"`
		Token  string    `json:"token"`
		UserID string    `json:"userID"`
		Role   Role      `json:"role"`
		Type   TokenType `json:"type"`
	}

	TokenType string

	ExternalNode struct {
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name"`
		Detail string `json:"detail"`
		Url    string `json:"url"`
		Secret string `json:"secret"`
	}

	ExternalNodeInstance struct {
		Permissions
		ID       string `json:"id" storm:"id"`
		NodeName string `json:"nodeName"`
	}

	ExternalQuery struct {
		*ExternalNode
		*ExternalNodeInstance
	}
)

const (
	TokenResetPassword TokenType = "reset_pw"
	TokenRegister      TokenType = "register_usr"
	TokenChangeEmail   TokenType = "change_email"
)

func (me *FormItem) GetVersion() int {
	return 0
}

func (me *FormItem) Clone() FormItem {
	form := *me
	form.ID = ""
	return form
}

func (me *FormComponentItem) GetVersion() int {
	return 0
}

func (me *WorkflowItem) GetVersion() int {
	return 0
}

func (me *WorkflowItem) Clone() WorkflowItem {
	workflow := *me
	workflow.ID = "" // without id the repository will create a new one
	return workflow
}

func (me *TemplateItem) Clone() TemplateItem {
	template := *me
	template.ID = ""
	return template
}

func (me *TemplateItem) GetVersion() int {
	return 1
}

func (me *TemplateItem) GetTemplate(lang string) (*file.IO, error) {
	if me.Data != nil {
		if tmpl, ok := me.Data[lang]; ok {
			if tmpl != nil {
				//exists..
				return tmpl, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

func (me *UserDataItem) GetVersion() int {
	return 0
}

func (me *SignatureRequestItem) GetVersion() int {
	return 0
}

func (me *WorkflowPaymentItem) GetVersion() int {
	return 0
}

func (me *UserDataItem) GetAllFileInfos() []*file.IO {
	return nil
}

func (me *WorkflowItem) LoopNodes(looper *workflow.Looper, cb func(l *workflow.Looper, node *workflow.Node) bool) {
	if me.Data != nil {
		if looper == nil {
			looper = &workflow.Looper{}
		}
		me.Data.Loop(looper, cb)
	}
}

func (e *ExternalNode) HealthUrl() string {
	return fmt.Sprintf("%s/health", e.Url)
}

func (e ExternalQuery) jwtToken() string {
	claims := jwt.StandardClaims{
		Id:        e.ExternalNodeInstance.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(e.Secret))
	return t
}

func (e ExternalQuery) nodeUrl(method string) string {
	return fmt.Sprintf("%s/node/%s/%s?auth=%s",
		e.Url,
		e.ExternalNodeInstance.ID,
		method,
		e.jwtToken(),
	)
}

func (e ExternalQuery) ConfigUrl() string {
	return e.nodeUrl("config")
}

func (e ExternalQuery) NextUrl() string {
	return e.nodeUrl("next")
}

func (e ExternalQuery) RemoveUrl() string {
	return e.nodeUrl("remove")
}

func (e ExternalQuery) CloseUrl() string {
	return e.nodeUrl("close")
}
