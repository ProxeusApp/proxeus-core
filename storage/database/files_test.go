package database

import (
	"bytes"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	. "github.com/onsi/gomega"
)

func TestFile(t *testing.T) {
	RegisterTestingT(t)
	fi := testDBSet.Files

	bytes1 := []byte("yyy")
	bytes2 := []byte("zzz")

	buf := &bytes.Buffer{}
	path := "/tmp/some/file"
	path2 := "/tmp/some/file2"

	// file writes
	Expect(fi.Write(path, bytes.NewBuffer(bytes1))).To(Succeed())
	Expect(fi.Write(path2, bytes.NewBuffer(bytes2))).To(Succeed())

	// file reads
	Expect(fi.Read(path, buf)).To(Succeed())
	Expect(buf.Bytes()).To(Equal(bytes1))
	buf.Reset()
	Expect(fi.Read(path2, buf)).To(Succeed())
	Expect(buf.Bytes()).To(Equal(bytes2))

	// read with err
	err := fi.Read("bla_ble", buf)
	Expect(db.NotFound(err)).To(Equal(true))
	Expect(fi.Read("", buf)).To(HaveOccurred())

	// exists or not
	Expect(fi.Exists(path)).To(Equal(true))
	Expect(fi.Exists(path + "/bla")).To(Equal(false))

	// test helper functions
	Expect(storage.FileSize(fi, path)).To(BeEquivalentTo(len(bytes1)))
	Expect(storage.CopyFile(fi, "file3", path)).To(BeEquivalentTo(len(bytes1)))
	buf.Reset()
	Expect(fi.Read("file3", buf)).To(Succeed())
	Expect(buf.Bytes()).To(Equal(bytes1))

	// delete first
	Expect(fi.Delete(path)).To(Succeed())
	Expect(fi.Exists(path)).To(Equal(false))

	// delete second
	Expect(fi.Exists(path2)).To(Equal(true))
	Expect(fi.Delete(path2)).To(Succeed())
	Expect(fi.Exists(path2)).To(Equal(false))
}
