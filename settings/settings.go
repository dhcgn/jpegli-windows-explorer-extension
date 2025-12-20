package settings

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

// Seetings represents the configuration options for the application
type Seetings struct {
	Distance             float64 `yaml:"distance"`
	OverrideOriginalFile bool    `yaml:"override_original_file"`
	SkipUpdateCheck      bool    `yaml:"skip_update_check"`
}

func configFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return configFileName // fallback
	}
	dir := filepath.Dir(exePath)
	return filepath.Join(dir, configFileName)
}

func GetConfigFilePath() string {
	return configFilePath()
}

func CheckForConfigFile() bool {
	cfgPath := configFilePath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func LoadOrDefault() (Seetings, string, error) {
	defaultOpts := Seetings{Distance: 0.5, OverrideOriginalFile: false, SkipUpdateCheck: false}
	cfgPath := configFilePath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		saveDefaultConfig(defaultOpts)
	}
	file, err := os.Open(cfgPath)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, cfgPath, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, cfgPath, err
	}
	var opts Seetings
	err = yaml.Unmarshal(data, &opts)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, cfgPath, err
	}

	return opts, cfgPath, nil
}

func saveDefaultConfig(cfg Seetings) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	_ = os.WriteFile(configFilePath(), data, 0644)
}
