package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/model"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gavv/httpexpect.v2"
)

var serviceURL string
var root *user

type session struct {
	base string
	id   string
	t    *testing.T
	e    *httpexpect.Expect
	c    *http.Client
	s    *model.Settings
	root *user
}

func TestApi(t *testing.T) {
	url := os.Getenv("PROXEUS_URL")
	if !isOnline(url) {
		t.Fatal("Service not online")
	}

	s := newSession(t, url)

	t.Logf("New test api session %s - root %s/%s", s.id, s.root.username, s.root.password)

	initProxeus(s)
	if s.t.Failed() {
		return
	}

	tests := []struct {
		name string
		f    func(s *session)
	}{
		{"User", testUser},
		{"UserProfile", testUserProfile},
		{"APIKey", testApiKey},
		{"Form", testForm},
		{"Payment", testPayment},
		{"UnattendedWorkflow", testUnattendedWorkflow},
		{"UnattendedWorkflowAdvanced", testUnattendedWorkflowAdvanced},
		{"Workflow", testWorkflow},
		{"WorkflowAdvanced", testWorkflowAdvanced},
		{"I18n", testI18n},
		{"I18nAdmin", testI18nAdmin},
		{"ImportExport", testImportExport},
		{"ImportExportAdmin", testImportExportAdmin},
		{"ImportExportRoot", testImportExportRoot},
		{"Template", testTemplate},
		{"WorkflowWithFile", testWorkflowWithFile},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { test.f(cloneSession(t, s)) })
	}
}

func newSession(t *testing.T, serverURL string) *session {
	id := uuid.NewV4().String()

	return &session{
		base: serverURL,
		id:   id,
		t:    t,
		e: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  serverURL,
			Reporter: httpexpect.NewAssertReporter(t),
			Printers: []httpexpect.Printer{
				//httpexpect.NewDebugPrinter(t, true),
				httpexpect.NewCompactPrinter(t),
				//httpexpect.NewCurlPrinter(t),
			},
		}),
		root: &user{
			username: fmt.Sprintf("testroot%s@example.com", id),
			password: id,
		},
	}
}

func cloneSession(t *testing.T, s *session) *session {
	id := uuid.NewV4().String()
	return &session{
		base: s.base,
		id:   id,
		t:    s.t,
		e: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  s.base,
			Reporter: httpexpect.NewAssertReporter(s.t),
			Printers: []httpexpect.Printer{
				//httpexpect.NewDebugPrinter(t, true),
				httpexpect.NewCompactPrinter(s.t),
				//httpexpect.NewCurlPrinter(t),
			},
		}),
		root: s.root,
	}
}

func isOnline(url string) bool {
	for i := 0; i < 10; i++ {
		r, err := http.Get(url)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer r.Body.Close()

		_, err = ioutil.ReadAll(r.Body)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		return true
	}
	return false
}

func initProxeus(s *session) {
	type user struct {
		Email    string     `json:"email"`
		Password string     `json:"password"`
		Role     model.Role `json:"role"`
	}

	type Init struct {
		Settings *model.Settings `json:"settings"`
		User     *user           `json:"user"`
	}

	r := s.e.GET("/api/init").Expect()

	if strings.Contains(r.Body().Raw(), `"configured":false`) { // Get /api/init returned the settings: Proxeus core is not initialised
		var init Init
		err := json.Unmarshal([]byte(r.Body().Raw()), &init)
		if err != nil {
			s.t.Fatal(err)
			return
		}

		init.User = &user{
			Email:    s.root.username,
			Password: s.root.password,
			Role:     100,
		}

		s.e.POST("/api/init").WithJSON(init).Expect().Status(http.StatusOK)
		// We upload components
		login(s, s.root)
		for _, name := range []string{"HC1", "HC2", "HC3", "HC5", "HC7", "HC8", "HC9", "HC10", "HC11", "HC12"} {
			c, err := Asset(fmt.Sprintf("test/assets/components/%s.json", name))
			if err != nil {
				s.t.Errorf("Cannot read asset %s", err)
			}
			s.e.POST("/api/admin/form/component").WithQuery("id", name).WithHeader("Content-Type", "application/json").WithBytes(c).Expect().Status(http.StatusOK)
		}

		logout(s)
	}
}
