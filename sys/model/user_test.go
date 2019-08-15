package model

import (
	"testing"
)

func TestUser_HideApiKeys(t *testing.T) {
	u := &User{ID: "123456789010"}
	_, err := u.NewApiKey("key1")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = u.NewApiKey("key2")
	if err != nil {
		t.Error(err)
		return
	}
	for _, a := range u.ApiKeys {
		a.HideKey()
	}
	for _, a := range u.ApiKeys {
		if len(a.Key) != 11 {
			t.Error("wrong length! key not hidden?")
		}
	}
}
