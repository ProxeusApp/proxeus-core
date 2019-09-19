package www

import (
	"encoding/base64"
	"testing"
)

func TestExtractApiKey(t *testing.T) {

	tests := []struct {
		title    string
		value    string
		expected string
	}{
		{
			"No header",
			"",
			"",
		},
		{
			"Authorization header but wrong type",
			"Basic 1234",
			"",
		},
		{
			"Authorization header right type, wrong spacing",
			"Bearer   1234",
			"",
		},
		{
			"Authorization header right type, wrong spacing2",
			" Bearer 1234",
			"",
		},
		{
			"Authorization header right type, wrong spacing3",
			"Bearer 1234 ",
			"",
		},
		{
			"Good",
			"Bearer 1234",
			"1234",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			result := extractSessionToken(test.value)
			if result != test.expected {
				t.Errorf("Expected %s and got %s", test.expected, result)
			}
		})
	}

}

func TestExtractBasicAuth(t *testing.T) {
	tests := []struct {
		title    string
		value    string
		user     string
		password string
	}{
		{
			"No header",
			"",
			"",
			"",
		},
		{
			"Authorization header but wrong type",
			"Bearer 1234",
			"",
			"",
		},
		{
			"Authorization header right type, wrong spacing",
			"Basic   " + base64.StdEncoding.EncodeToString([]byte("foo:bar")),
			"",
			"",
		},
		{
			"Authorization header right type, wrong spacing2",
			" Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")),
			"",
			"",
		},
		{
			"authorization header right type, wrong spacing3",
			"Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")) + " ",
			"",
			"",
		},
		{
			"authorization header right type, wrong content",
			"Basic " + base64.StdEncoding.EncodeToString([]byte("foo:")),
			"foo",
			"",
		},
		{
			"authorization header right type, wrong content",
			"Basic " + base64.StdEncoding.EncodeToString([]byte(":bar")),
			"",
			"bar",
		},
		{
			"authorization header right type, wrong content",
			"Basic " + base64.StdEncoding.EncodeToString([]byte(":")),
			"",
			"",
		},
		{
			"authorization header right type, wrong content",
			"Basic " + base64.StdEncoding.EncodeToString([]byte("")),
			"",
			"",
		},
		{
			"Good",
			"Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")),
			"foo",
			"bar",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			u, p := extractBasicAuth(test.value)
			if u != test.user || p != test.password {
				t.Errorf("Expected %s:%s and got %s:%s", test.user, test.password, u, p)
			}
		})
	}
}
