package customNode

import (
	"encoding/json"
	"testing"
)

func TestChangeDataBeforeSend(t *testing.T) {
	dat := map[string]interface{}{"input": map[string]interface{}{"CapitalSource": []interface{}{"Andere"}}}
	newDat := changeDataBeforeSend(dat)
	bts, _ := json.Marshal(newDat)
	if string(bts) != `{"CapitalSource":"[\"Andere\"]"}` {
		t.Error(string(bts))
	}
}
