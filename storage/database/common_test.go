package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

const testEngine = "storm"

type dummyAuth struct {
	role model.Role
}

var dummy model.Auth = dummyAuth{role: model.CREATOR}
var dummySuperAdmin model.Auth = dummyAuth{role: model.SUPERADMIN}

func (d dummyAuth) UserID() string {
	if d.role == model.SUPERADMIN {
		return "124456789"
	}
	return "10"
}
func (d dummyAuth) AccessRights() model.Role {
	return d.role
}

func dummySettings() *model.Settings {
	s := model.NewDefaultSettings()
	s.BlockchainContractAddress = "R"
	s.SparkpostApiKey = "R"
	s.EmailFrom = "R@R.com"
	s.PlatformDomain = "r.com"
	s.InfuraApiKey = "R"
	s.DatabaseEngine = testEngine
	return s
}

var testDBSet *storage.DBSet

func TestMain(m *testing.M) {
	dir := filepath.Join(os.TempDir(), "database_unit_tests")
	os.RemoveAll(dir) // in case previous test run was killed
	var code int
	defer func() { os.Exit(code) }() // needs to be wrapped for code value binding
	defer os.RemoveAll(dir)

	maybeFail := func(err error) {
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			debug.PrintStack()
			code = 1
		}
	}

	se, err := NewSettingsDB(dir)
	maybeFail(err)
	err = se.Put(dummySettings())
	maybeFail(err)

	testDBSet, err = NewDBSet(se, dir)
	defer testDBSet.Close()
	maybeFail(err)

	if code == 0 {
		code = m.Run()
	}
}

func assignTimeField(item interface{}, fieldName string, t time.Time) {
	v := reflect.ValueOf(item)
	vt := reflect.ValueOf(t)
	if v.Kind() == reflect.Slice {
		for k := 0; k < v.Len(); k++ {
			el := v.Index(k)
			sv := reflect.Indirect(el).FieldByName(fieldName)
			if sv.IsValid() {
				sv.Set(vt)
			}
		}
		return
	}
	if v.Kind() == reflect.Ptr && reflect.Indirect(v).Kind() == reflect.Struct {
		sv := reflect.Indirect(v).FieldByName(fieldName)
		if sv.IsValid() {
			sv.Set(vt)
		}
	}
}

func updateTimeFields(item interface{}, t time.Time) {
	for _, f := range []string{"Created", "CreatedAt", "Updated"} {
		assignTimeField(item, f, t)
	}
}

func equalJSON(item interface{}) types.GomegaMatcher {
	ti := time.Now()
	updateTimeFields(item, ti)
	eqFunc := func() types.GomegaMatcher {
		d, _ := json.MarshalIndent(item, "", " ")
		return gomega.MatchJSON(d)
	}
	transformFunc := func(i interface{}) []byte {
		updateTimeFields(i, ti)
		d, _ := json.MarshalIndent(i, "", " ")
		return d
	}
	return gomega.WithTransform(transformFunc, eqFunc())
}
