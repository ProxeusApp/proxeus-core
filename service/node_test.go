package service

import (
	"testing"
)

func TestList(t *testing.T) {

	nodeService := NewNodeService(nil, nil)

	if nodeService.List("mailsender") == nil {
		t.Fail()
	}
	if nodeService.List("priceretriever") == nil {
		t.Fail()
	}
}
