package database

import (
	"testing"

	"github.com/onsi/gomega/types"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	. "github.com/onsi/gomega"
)

func TestSignatureRequests(t *testing.T) {
	RegisterTestingT(t)
	si := testDBSet.SignatureRequests

	item := &model.SignatureRequestItem{
		ID:        "1",
		DocId:     "doc1",
		DocPath:   "/some/path",
		Hash:      "2323232",
		Requestor: "0xfefe",
		Signatory: "0xbeef",
	}
	item2 := &model.SignatureRequestItem{
		ID:        "2",
		DocId:     "doc2",
		DocPath:   "/some/path2",
		Hash:      "4532234",
		Requestor: "0xfefe",
		Signatory: "0xbeef",
	}

	array := func(args ...*model.SignatureRequestItem) types.GomegaMatcher {
		var r []*model.SignatureRequestItem
		r = append(r, args...)
		return equalJSON(&r)
	}

	// add
	Expect(si.Add(item)).To(Succeed())
	Expect(si.Add(item2)).To(Succeed())

	// get
	Expect(si.GetByID(item.DocId, item.DocPath)).To(array(item))
	Expect(si.GetBySignatory(item.Signatory)).To(array(item, item2))
	Expect(si.GetByHashAndSigner(item2.Hash, item2.Signatory)).To(array(item2))

	// rejected or revoked
	Expect(si.SetRejected(item.DocId, item.DocPath, item.Signatory)).To(Succeed())
	Expect(si.SetRevoked(item2.DocId, item2.DocPath, item2.Signatory)).To(Succeed())
	gotItems, _ := si.GetByID(item.DocId, item.DocPath)
	Expect((*gotItems)[0].Rejected).To(Equal(true))
	gotItems2, _ := si.GetByID(item2.DocId, item2.DocPath)
	Expect((*gotItems2)[0].Revoked).To(Equal(true))
}
