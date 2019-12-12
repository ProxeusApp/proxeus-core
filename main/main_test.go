// +build coverage

package main

import (
	"testing"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
)

func TestCoverage(m *testing.T) {
	cfg.Config.Settings.TestMode = "true"
	cfg.Config.ServiceAddress = ":1323"
	main()
}
