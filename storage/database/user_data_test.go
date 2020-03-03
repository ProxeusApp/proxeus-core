package database

import (
	"bytes"
	"testing"

	"github.com/ProxeusApp/proxeus-core/sys/model/compatability"

	"github.com/ProxeusApp/proxeus-core/sys/file"

	"github.com/ProxeusApp/proxeus-core/storage"

	. "github.com/onsi/gomega"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func TestUserData(t *testing.T) {
	RegisterTestingT(t)
	ud := testDBSet.UserData
	filesIF := testDBSet.Files

	options := storage.IndexOptions(0)
	fio := ud.NewFile(dummy, file.Meta{Name: "some_file"})
	wf := &model.WorkflowItem{ID: "xy"}
	item := &model.UserDataItem{
		Name:       "test_ud",
		WorkflowID: wf.ID,
		Data:       compatability.CarriedStringMap{},
	}

	Expect(ud.AssetsKey()).ToNot(BeEmpty())

	// simple put and get
	Expect(ud.Put(dummy, item)).To(Succeed())
	gotItem, _ := ud.Get(dummy, item.ID)
	Expect(gotItem).To(equalJSON(item))
	item = gotItem
	Expect(ud.Put(dummy, item)).To(Succeed())

	Expect(ud.List(dummy, "", options, true)).
		To(equalJSON([]*model.UserDataItem{item}))

	// data put and get
	dataObj := map[string]interface{}{"input": map[string]interface{}{
		"some_file_key": fio,
		"text":          "text_value",
	}}
	Expect(ud.PutData(dummy, item.ID, dataObj)).To(Succeed())
	Expect(ud.GetData(dummy, item.ID, "input")).To(Equal(dataObj["input"]))

	// get data input and files
	gotDataInput, _, _ := ud.GetDataAndFiles(dummy, item.ID, "input.text")
	Expect(gotDataInput).To(Equal("text_value"))
	_, gotFiles, _ := ud.GetDataAndFiles(dummy, item.ID, "input")
	Expect(gotFiles).To(Equal([]string{fio.Path()}))

	// GetAllFileInfosOf
	item, _ = ud.Get(dummy, item.ID)
	Expect(ud.GetAllFileInfosOf(item)).To(Equal([]*file.IO{fio}))

	// get by workflow id
	gotItem, _, _ = ud.GetByWorkflow(dummy, wf, false)
	Expect(gotItem).To(Equal(item))

	// GetDataFile
	Expect(ud.GetDataFile(dummy, item.ID, "input.some_file_key")).To(Equal(fio))

	// put with non existing ID
	gotItem.ID = "333"
	Expect(ud.Put(dummy, gotItem)).To(Succeed())

	// delete with files
	buf := bytes.NewBuffer([]byte("xxx"))
	Expect(filesIF.Write(fio.Path(), buf)).To(Succeed())
	Expect(ud.Delete(dummy, filesIF, item.ID)).To(Succeed())
}
