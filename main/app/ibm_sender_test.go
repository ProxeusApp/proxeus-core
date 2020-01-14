package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/file"
)

func TestChangeDataBeforeSend(t *testing.T) {
	dat := map[string]interface{}{"input": map[string]interface{}{"CapitalSource": []interface{}{"Andere"}}}
	newDat := changeDataBeforeSend(dat)
	bts, _ := json.Marshal(newDat)
	if string(bts) != `{"CapitalSource":"[\"Andere\"]"}` {
		t.Error(string(bts))
	}
}

func TestIBMSenderExecute(t *testing.T) {

	tests := []struct {
		title  string
		data   map[string]interface{}
		header string
		body   string

		enabled       string
		expectError   bool
		expectProceed bool
	}{
		{
			title:         "Not enabled",
			data:          map[string]interface{}{"input": map[string]interface{}{"CapitalSource": []interface{}{"Andere"}}},
			expectError:   true,
			expectProceed: false,
		},
		{
			title:         "Enabled",
			enabled:       "true",
			data:          map[string]interface{}{"input": map[string]interface{}{"CapitalSource": []interface{}{"Andere"}}},
			header:        `{"Accept-Encoding":["gzip"],"Clientid":["test-client-id"],"Content-Length":["38"],"Content-Type":["application/json"],"Oauthserverurl":["test-oauth-url"],"Secret":["test-secret"],"Tenantid":["test-tenant-id"],"User-Agent":["Go-http-client/1.1"]}`,
			body:          `{"input":{"CapitalSource":["Andere"]}}`,
			expectError:   false,
			expectProceed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {

			os.Setenv("FF_IBM_SENDER_ENABLED", test.enabled)
			os.Setenv("FF_IBM_SENDER_CLIENT_ID", "test-client-id")
			os.Setenv("FF_IBM_SENDER_TENANT_ID", "test-tenant-id")
			os.Setenv("FF_IBM_SENDER_SECRET", "test-secret")
			os.Setenv("FF_IBM_SENDER_OAUTH_URL", "test-oauth-url")
			os.Setenv("FF_IBM_SENDER_URL", "test-url")

			var header http.Header
			var body string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				header = r.Header
				b, _ := ioutil.ReadAll(r.Body)
				body = string(b)
			}))
			defer ts.Close()

			forwardURL = ts.URL

			sender := &IBMSenderNodeImpl{
				ctx: &DocumentFlowInstance{
					DataCluster: &mockDataManager{
						allData: test.data,
					},
				},
			}

			p, err := sender.Execute(nil)
			if test.expectError && err == nil {
				t.Error("Expected an error")
			}

			if !test.expectError && err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if test.expectProceed != p {
				t.Errorf("Wrong proceed result: %t", p)
			}

			if !p {
				return
			}

			j, err := json.Marshal(header)
			if test.header != string(j) {
				t.Errorf("Expected header:%s\n but got:%s", test.header, string(j))
			}

			if test.body != body {
				t.Errorf("Expected body:%s\n but got:%s", test.body, string(body))
			}
		})
	}
}

// Mock

type mockDataManager struct {
	allData map[string]interface{}
}

func (m *mockDataManager) OnLoad() {}

func (m *mockDataManager) GetDataFile(formID, name string) (*file.IO, error) {
	return nil, nil
}

func (m *mockDataManager) GetData(formID string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *mockDataManager) GetDataByPath(formID, dataPath string) (interface{}, error) {
	return nil, nil
}

func (m *mockDataManager) Clear(formID string) error {
	return nil
}

func (m *mockDataManager) GetAllData() (dat map[string]interface{}, err error) {
	return m.allData, nil
}

func (m *mockDataManager) GetAllDataFilePathNameOnly() (dat map[string]interface{}, files []string) {
	return nil, nil
}

func (m *mockDataManager) PutData(formID string, dat map[string]interface{}) error {
	return nil
}

func (m *mockDataManager) PutDataWithoutMerge(formID string, dat map[string]interface{}) error {
	return nil
}

func (m *mockDataManager) PutDataFile(db storage.FilesIF, formID, name string, f file.Meta, reader io.Reader) error {
	return nil
}

func (m *mockDataManager) Close() (err error) {
	return nil
}
