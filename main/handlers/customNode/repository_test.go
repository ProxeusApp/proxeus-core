package customNode

import (
	"testing"
)

func TestList(t *testing.T) {
	if List(nil, "priceretriever") == nil {
		t.Fail()
	}
}
