package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI18nResolver_Resolve(t *testing.T) {
	i18nResolver := I18nResolver{}

	t.Run("with one attribute", func(t *testing.T) {
		response := i18nResolver.Resolve("Hi, my job is to {0}", "translate")

		assert.Equal(t, "Hi, my job is to translate", response)
	})

	t.Run("with multiple attributes and lines", func(t *testing.T) {
		response := i18nResolver.Resolve("Hello {0}, welcome to our "+
			"\nwebsite, {1}", "Silvio", "Proxeus.com")

		assert.Equal(t, "Hello Silvio, welcome to our \nwebsite, Proxeus.com", response)
	})

	t.Run("with non existing attribute", func(t *testing.T) {
		response := i18nResolver.Resolve("No attributes here", "I am an attribute")

		assert.Equal(t, "No attributes here", response)
	})
}
