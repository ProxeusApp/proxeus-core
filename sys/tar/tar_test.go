package tar

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"

	_ "github.com/ProxeusApp/proxeus-core/test"
)

func TestTar(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), random.String(5, random.Alphanumeric))
	defer os.RemoveAll(tmpDir)

	dirToZip := "./test/assets/"

	err := os.MkdirAll(tmpDir, 0766)
	if err != nil {
		panic(err)
	}

	tmpFile := filepath.Join(tmpDir, "filename.tar.gz")

	f, err := os.Create(tmpFile)
	if err != nil {
		panic(err)
	}

	err = Tar(dirToZip, f)
	assert.Nil(t, err, "should not return error")
	stat, _ := os.Stat(tmpFile)
	assert.Truef(t, stat.Size() > 100, "should be bigger than 100 but is %d", stat.Size())
}

func TestUntar(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), random.String(5, random.Alphanumeric))
	defer os.RemoveAll(tmpDir)

	file, err := os.Open("./test/assets/test-tar.tar.gz")
	if err != nil {
		panic(err)
	}

	err = Untar(tmpDir, file)
	assert.Nil(t, err, "should not return error")

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 2, len(files))
	// Check files have been extracted of the same original size. Function returns them sorted by filename
	assert.Equal(t, int64(60), files[0].Size())
	assert.Equal(t, int64(225561), files[1].Size())
}
