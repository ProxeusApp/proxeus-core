package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLang_Matches(t *testing.T) {
	lang := Lang{
		ID:      "id1",
		Code:    "en",
		Enabled: false,
	}

	t.Run("when given same code", func(t *testing.T) {
		result := lang.Matches("en")

		assert.True(t, result)
	})

	t.Run("when given a different code", func(t *testing.T) {
		result := lang.Matches("it")

		assert.False(t, result)
	})

	t.Run("when given a complex code with underscore", func(t *testing.T) {
		result := lang.Matches("en_UK")

		assert.True(t, result)
	})

	t.Run("when given a complex code with dash", func(t *testing.T) {
		result := lang.Matches("en-UK")

		assert.True(t, result)
	})
}
