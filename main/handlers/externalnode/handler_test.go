package externalnode

import (
	"testing"
)

func TestList(t *testing.T) {
	if List(nil, "mailsender") == nil {
		t.Fail()
	}
	if List(nil, "priceretriever") == nil {
		t.Fail()
	}
}
