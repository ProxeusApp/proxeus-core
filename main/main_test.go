package main

import (
	cfg "git.proxeus.com/core/central/main/config"
	"testing"
)

func TestCoverage(m *testing.T) {
	cfg.Config.Settings.TestMode = "true"
	cfg.Config.ServiceAddress = ":1323"
	main()
}