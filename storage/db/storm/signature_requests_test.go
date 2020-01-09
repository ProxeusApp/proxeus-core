package storm

import (
	"fmt"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func TestSigning(t *testing.T) {
	baseDir := "./testDir"
	sigdb, err := NewSignatureDB(DBConfig{Dir: baseDir})
	defer func() { os.RemoveAll(baseDir) }()
	if err != nil {
		t.Error(err)
	}

	sigreq := &model.SignatureRequestItem{}
	sigreq.ID = uuid.NewV4().String()
	sigreq.DocId = "abcd"
	sigreq.DocPath = "1234"
	sigreq.Signatory = "signatory_1"
	sigreq.Requestor = "requestor"
	sigreq.RequestedAt = time.Now()

	sigreq2 := &model.SignatureRequestItem{}
	sigreq2.ID = uuid.NewV4().String()
	sigreq2.DocId = "abcd"
	sigreq2.DocPath = "5678"
	sigreq2.Signatory = "signatory_1"
	sigreq2.Requestor = "requestor"
	sigreq2.RequestedAt = time.Now()
	sigreq2.RejectedAt = time.Now()

	sigreq3 := &model.SignatureRequestItem{}
	sigreq3.ID = uuid.NewV4().String()
	sigreq3.DocId = "abcd"
	sigreq3.DocPath = "1234"
	sigreq3.Signatory = "signatory_2"
	sigreq3.Requestor = "requestor"

	fmt.Println("Start")

	fmt.Println("Add 1")
	err = sigdb.Add(sigreq)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Add 2")
	err = sigdb.Add(sigreq2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Add 3")
	err = sigdb.Add(sigreq3)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("GetBySignatory")
	_, err = sigdb.GetBySignatory("signatory_1")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("SetRejected")
	err = sigdb.SetRejected(sigreq.DocId, sigreq.DocPath, sigreq.Signatory)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("GetByID")
	_, err = sigdb.GetByID(sigreq.DocId, sigreq.DocPath)
	if err != nil {
		t.Error(err)
	}
}
