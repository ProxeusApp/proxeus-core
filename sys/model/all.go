package model

import (
	"os"
	"time"

	"git.proxeus.com/core/central/sys/file"
	"git.proxeus.com/core/central/sys/workflow"
)

type (
	//Structure interface helps us to keep track of the inner structure of a struct between persisted and memory data
	Structure interface {
		GetVersion() int
	}

	FormItem struct {
		Permissions
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name" storm:"index"`
		Detail string `json:"detail"`
		//Permissions Permissions `json:"permissions"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`

		//Data contains the form source
		Data map[string]interface{} `json:"data"`
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

		Data map[string]interface{} `json:"data"`
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
		Settings interface{} `json:"settings"`
		Template interface{} `json:"template"`
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
		ID int `storm:"id,increment"`
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
		Hash       string `json:"hash" storm:"id"`
		WorkflowID string `json:"workflowID" storm:"index"`
		From       string `json:"From"`
		To         string `json:"To"`
		Xes        uint64 `json:"xes"`
	}

	TemplateLangMap map[string]*file.IO
)

func (me *FormItem) GetVersion() int {
	return 0
}

func (me *FormComponentItem) GetVersion() int {
	return 0
}

func (me *WorkflowItem) GetVersion() int {
	return 0
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
