package www

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	_ "github.com/ProxeusApp/proxeus-core/test"
)

func TestEmbedded_Asset2(t *testing.T) {

	embedded := Embedded{
		Asset: func(name string) (bytes []byte, err error) {
			// Return the same given content
			return []byte(name), nil
		},
	}

	t.Run("when name is empty", func(t *testing.T) {
		_, err := embedded.Asset2("")

		assert.Equal(t, os.ErrNotExist, err, "should return os.ErrNotExist")
	})

	t.Run("when given a name with slash prefixed", func(t *testing.T) {
		b, err := embedded.Asset2("/text")

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "text", string(b), "should take the second part")
	})

	t.Run("when given a name without slash prefixed", func(t *testing.T) {
		b, err := embedded.Asset2("text")

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "text", string(b), "should return the same text")
	})

}

func TestEmbedded_FindAssetWithC(t *testing.T) {

	contentTypeCache = nil

	file, err := os.Open("./test/assets/image.jpg")
	if err != nil {
		panic(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	embedded := Embedded{Asset: func(name string) (bytes []byte, err error) {
		return fileBytes, nil
	}}

	t.Run("when empty name is given", func(t *testing.T) {
		ct := "any"
		_, err := embedded.FindAssetWithContentType("", &ct)

		assert.Equal(t, err, echo.ErrNotFound, "should return echo.ErrNotFound")
	})

	t.Run("when correct data is given", func(t *testing.T) {
		ct := "any"
		b, err := embedded.FindAssetWithContentType("http://blabla.com/image.jpg", &ct)

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, fileBytes, b, "should return same file")
		assert.Equal(t, "image/jpeg", ct, "should change content type to image/jpeg")
	})
}
