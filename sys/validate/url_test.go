package validate

import "testing"

func TestUrlValidation(t *testing.T) {
	rules := Rules{"url": true}
	if len(Field("http://google.com", rules)) != 0 {
		t.Error("url validator error")
	}
	if len(Field("http//google.com", rules)) == 0 {
		t.Error("url validator error")
	}
	if len(Field("google.com", rules)) == 0 {
		t.Error("url validator error")
	}
	if len(Field("/foo/bar", rules)) == 0 {
		t.Error("url validator error")
	}
}
