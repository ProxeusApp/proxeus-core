package database

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterTestingT(t)
	// cipher key
	secretKey := "thisis32bitlongpassphraseimusing"

	// plaintext
	pt := "11165875f50c4b87a32a501afa79bf64"

	c, err := EncryptWithAES(secretKey, pt)
	Expect(err).To(BeNil())

	// decrypt
	decrepted, err := DecryptWithAES(secretKey, c)
	Expect(err).To(BeNil())

	Expect(decrepted).To(Equal(pt))
}
