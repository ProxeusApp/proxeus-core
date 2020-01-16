// +build coverage

package main

import (
	"os"
	"testing"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
)

func TestCoverage(m *testing.T) {
	cfg.Config.Settings.TestMode = "true"

	if sa := os.Getenv("PROXEUS_SERVICE_ADDRESS"); sa != "" {
		cfg.Config.ServiceAddress = sa
	}

	if dd := os.Getenv("PROXEUS_DATA_DIR"); dd != "" {
		cfg.Config.DataDir = dd
	}

	main()
}
