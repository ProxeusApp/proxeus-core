package storm

import (
	"os"
	"testing"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func TestProxeusDB(t *testing.T) {
	testDir := "./myTestDir"
	imex, err := NewImex(&model.User{Role: model.ROOT}, nil, testDir)
	if err != nil {
		t.Error(err)
	}
	err = imex.writeProxeusIdentifier(&ImexMeta{Version: 123})
	if err != nil {
		t.Error(err)
	}
	var imexMeta ImexMeta
	err = imex.readProxeusIdentifier(&imexMeta)
	if err != nil || imexMeta.Version != 123 {
		t.Error(err)
	}
	_ = imex.Close()
	_ = os.RemoveAll(testDir)
}
