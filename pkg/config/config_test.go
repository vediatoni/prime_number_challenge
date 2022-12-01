package config

import (
	"testing"

	"github.com/vediatoni/prime_number_challenge/internal/background"
)

// TestLoadConfigFromFile - tries to load specified config files. It doesn't check if the configuration is correct.
// That is (or should be) done at start.
func TestLoadConfigFromFile(t *testing.T) {
	type ConfigFile struct {
		FilePath   string
		StructType interface{}
	}

	cfgFile := ConfigFile{
		"../../config/dev.background.yaml",
		background.Config{},
	}

	t.Logf("Testing config file: %s", cfgFile.FilePath)

	loadedInf, err := LoadConfigFromFile(cfgFile.FilePath, cfgFile.StructType)
	if err != nil {
		t.Errorf("Error loading config file %s: %s", cfgFile, err)
	}

	cfg := loadedInf.(background.Config)

	if cfg.SelfSvcAddress != "localhost:50051" {
		t.Errorf("SelfSvcAddress is not set correctly. Expected: localhost:50051, got: %s", cfg.SelfSvcAddress)
	}

}
