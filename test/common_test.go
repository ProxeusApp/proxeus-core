package test

import (
	"encoding/json"
	"regexp"
)

func toMap(i interface{}) map[string]interface{} {
	var r map[string]interface{}
	j, _ := json.Marshal(i)
	json.Unmarshal(j, &r)
	return r
}

func removeUpdatedField(i map[string]interface{}) map[string]interface{} {
	delete(i, "updated")
	return i
}

// Removes the variable data in the PDF like for example the creation date or the checksum
var cleanPDFRegexp = regexp.MustCompile(`(?s)<<\/Creator.+?>>|<<\/Size.+?>>`)

func cleanPDF(src []byte) []byte {
	return cleanPDFRegexp.ReplaceAll(src, []byte{})
}
