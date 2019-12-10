package customNode

import (
	"testing"
)

func TestList(t *testing.T) {
	if List("mailsender") == nil {
		t.Fail()
	}
	if List("priceretriever") == nil {
		t.Fail()
	}
}
