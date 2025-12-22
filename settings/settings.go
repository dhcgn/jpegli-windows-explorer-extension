package settings

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

// Settings represents the configuration options for the application
type Settings struct {
	Distance             float64 `yaml:"distance"`
	OverrideOriginalFile bool    `yaml:"override_original_file"`
	SkipUpdateCheck      bool    `yaml:"skip_update_check"`
	NoUserInteraction    bool    `yaml:"no_user_interaction"`
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

func LoadOrDefault() (Settings, string, error) {
	defaultOpts := Settings{Distance: 0.5, OverrideOriginalFile: false, SkipUpdateCheck: false, NoUserInteraction: false}
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
	var opts Settings
	err = yaml.Unmarshal(data, &opts)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, cfgPath, err
	}

	return opts, cfgPath, nil
}

func saveDefaultConfig(cfg Settings) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	_ = os.WriteFile(configFilePath(), data, 0644)
}
