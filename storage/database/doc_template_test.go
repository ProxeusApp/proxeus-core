package database

import (
	"bytes"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	. "github.com/onsi/gomega"
)

func TestDocTemplate(t *testing.T) {
	RegisterTestingT(t)
	te := testDBSet.Template
	filesIF := testDBSet.Files

	options := storage.IndexOptions(0)
	vars := []string{"value1", "value2"}

	fio := file.New(te.AssetsKey(), file.Meta{})
	item := &model.TemplateItem{
		Name: "item1",
		Data: model.TemplateLangMap{"en": fio},
	}
	item2 := &model.TemplateItem{
		ID:   "2",
		Name: "item2",
		Data: model.TemplateLangMap{"en": fio},
	}

	// simple put and get
	Expect(te.AssetsKey()).ToNot(BeEmpty())
	Expect(te.Put(dummy, item)).To(Succeed())
	Expect(te.Put(dummy, item2)).To(Succeed())
	item2.Detail = "some details"
	Expect(te.Put(dummy, item2)).To(Succeed())
	Expect(te.Get(dummy, item.ID)).To(equalJSON(item))

	Expect(te.PutVars(dummy, item.ID, "en", vars)).To(Succeed())
	Expect(te.Vars(dummy, "", options.WithInclude([]string{item.ID + "en"}))).To(Equal(vars))

	Expect(te.List(dummy, "", options)).
		To(equalJSON([]*model.TemplateItem{item2, item}))

	Expect(te.ProvideFileInfoFor(dummy, item.ID, "en", &file.Meta{})).To(Equal(fio))
	Expect(te.GetTemplate(dummy, item.ID, "en")).To(Equal(fio))

	// deleting template generated document
	buf := bytes.NewBuffer([]byte("xxx"))
	Expect(filesIF.Write(fio.Path(), buf)).To(Succeed())
	Expect(te.DeleteTemplate(dummy, filesIF, item.ID, "en")).To(Succeed())
	_, err := te.GetTemplate(dummy, item.ID, "en")
	Expect(err).To(HaveOccurred())

	// adding template generated document
	fio2, err := te.ProvideFileInfoFor(dummy, item.ID, "de", &file.Meta{})
	Expect(err).To(Succeed())
	Expect(te.GetTemplate(dummy, item.ID, "de")).To(Equal(fio2))

	// renaming a file
	fio2, err = te.ProvideFileInfoFor(dummy, item.ID, "de", &file.Meta{Name: "f2"})
	Expect(err).To(Succeed())

	// listing modified
	item.Data = model.TemplateLangMap{"de": fio2}
	Expect(te.List(dummy, "", options)).
		To(equalJSON([]*model.TemplateItem{item, item2}))

	// deleting template
	Expect(te.Delete(dummy, filesIF, item2.ID)).To(Succeed())
	Expect(te.List(dummy, "", options)).
		To(equalJSON([]*model.TemplateItem{item}))
}
