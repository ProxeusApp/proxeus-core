package test

import (
	"encoding/json"
	"reflect"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func toMap(i interface{}) map[string]interface{} {
	var r map[string]interface{}
	j, _ := json.Marshal(i)
	json.Unmarshal(j, &r)
	return r
}

func removeTimeFields(i map[string]interface{}) map[string]interface{} {
	delete(i, "updated")
	delete(i, "created")
	delete(i, "createdAt")
	return i
}

func arrayContainsMap(t *testing.T, array *httpexpect.Array, expected interface{}) {
	var contains bool
	exMap := removeTimeFields(toMap(expected))
	for _, el := range array.Iter() {
		elMap := removeTimeFields(toMap(el.Object().Raw()))
		if reflect.DeepEqual(exMap, elMap) {
			contains = true
		}
	}
	if !contains {
		t.Errorf("exptect to find object %v in array %v", expected, array.Raw())
	}
}

// Removes the variable data in the PDF like for example the creation date or the checksum
var cleanPDFRegexp = regexp.MustCompile(`(?s)<<\/Creator.+?>>|<<\/Size.+?>>`)

func cleanPDF(src []byte) []byte {
	return cleanPDFRegexp.ReplaceAll(src, []byte{})
}
