package database

import (
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"

	. "github.com/onsi/gomega"
)

func TestI18n(t *testing.T) {
	RegisterTestingT(t)
	in := testDBSet.I18n

	options := storage.IndexOptions(0)
	lang1 := "en"
	lang2 := "de"

	translations := map[string]string{
		"login":   "login text",
		"welcome": "welcome text",
	}

	// get with no cache
	_, err := in.GetAll(lang1)
	Expect(err).To(Succeed())

	// fallback lang
	Expect(in.PutFallback(lang2)).To(Succeed())
	Expect(in.GetFallback()).To(Equal(lang2))

	// put languages
	Expect(in.PutLang(lang1, true)).To(Succeed())
	Expect(in.PutLang(lang2, false)).To(Succeed())
	Expect(in.HasLang(lang1)).To(Equal(true))
	Expect(in.HasLang("fr")).To(Equal(false))

	// check languages
	enabledLangs, _ := in.GetLangs(true)
	Expect(len(enabledLangs)).To(Equal(1))
	allLangs, _ := in.GetAllLangs()
	Expect(len(allLangs)).To(Equal(2))

	// put translations
	Expect(in.PutAll(lang1, translations)).To(Succeed())
	Expect(in.Put(lang1, "exit", "exit text")).To(Succeed())

	// get translations
	Expect(in.Get(lang1, "welcome")).To(Equal("welcome text"))
	translations["exit"] = "exit text"
	Expect(in.GetAll(lang1)).To(Equal(translations))

	// finds
	Expect(in.Find("login", "", options)).
		To(Equal(map[string]map[string]string{"login": {lang1: "login text"}}))
	Expect(in.Find("", "welcome", options)).
		To(Equal(map[string]map[string]string{"welcome": {lang1: "welcome text"}}))
	Expect(in.Find("welcome", "welcome", options)).
		To(Equal(map[string]map[string]string{"welcome": {lang1: "welcome text"}}))
	_, err = in.Find("", "", options)
	Expect(err).To(Succeed())
}
