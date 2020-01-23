package storage

import (
	"os"
	"path/filepath"
	"testing"
)

type flow struct {
	Test string
}

func TestJSONFile(t *testing.T) {
	fPath := filepath.Join(os.TempDir(), "test.json")
	defer os.Remove(fPath)
	jf := JSONFile{filePath: fPath, perm: 0600}

	// write
	fw := &flow{Test: "my json"}
	err := jf.Put(fw)
	if err != nil {
		t.Error(err)
	}

	// read
	var m flow
	err = jf.Get(&m)
	if err != nil {
		t.Error(err)
	}
	if m.Test != fw.Test {
		t.Errorf("expected %v got %v", fw.Test, m.Test)
	}
}
