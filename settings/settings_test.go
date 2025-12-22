package settings

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSettingsDefaultValues(t *testing.T) {
	// Create a default Settings instance
	defaultOpts := Settings{Distance: 0.5, OverrideOriginalFile: false}

	// Check default values
	if defaultOpts.Distance != 0.5 {
		t.Errorf("Expected Distance to be 0.5, got %f", defaultOpts.Distance)
	}
	if defaultOpts.OverrideOriginalFile != false {
		t.Errorf("Expected OverrideOriginalFile to be false, got %v", defaultOpts.OverrideOriginalFile)
	}
}

func TestSettingsYAMLSerialization(t *testing.T) {
	// Test YAML marshaling
	opts := Settings{Distance: 1.5, OverrideOriginalFile: true}
	data, err := yaml.Marshal(opts)
	if err != nil {
		t.Fatalf("Failed to marshal Settings: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Settings
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Settings: %v", err)
	}

	// Check values
	if unmarshaled.Distance != opts.Distance {
		t.Errorf("Expected Distance to be %f, got %f", opts.Distance, unmarshaled.Distance)
	}
	if unmarshaled.OverrideOriginalFile != opts.OverrideOriginalFile {
		t.Errorf("Expected OverrideOriginalFile to be %v, got %v", opts.OverrideOriginalFile, unmarshaled.OverrideOriginalFile)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "jpegli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test config
	testConfig := Settings{Distance: 2.5, OverrideOriginalFile: true}
	configPath := filepath.Join(tempDir, "config.yaml")

	// Marshal and save
	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load and verify
	fileData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var loaded Settings
	err = yaml.Unmarshal(fileData, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loaded.Distance != testConfig.Distance {
		t.Errorf("Expected Distance to be %f, got %f", testConfig.Distance, loaded.Distance)
	}
	if loaded.OverrideOriginalFile != testConfig.OverrideOriginalFile {
		t.Errorf("Expected OverrideOriginalFile to be %v, got %v", testConfig.OverrideOriginalFile, loaded.OverrideOriginalFile)
	}
}
