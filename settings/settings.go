package settings

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dhcgn/jpegli-windows-explorer-extension/convert"
	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

func configFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return configFileName // fallback
	}
	dir := filepath.Dir(exePath)
	return filepath.Join(dir, configFileName)
}

func LoadOrDefault() (convert.ConvertOptions, error) {
	defaultOpts := convert.ConvertOptions{Distance: 0.5}
	cfgPath := configFilePath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		saveDefaultConfig(defaultOpts)
	}
	file, err := os.Open(cfgPath)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, err
	}
	var opts convert.ConvertOptions
	err = yaml.Unmarshal(data, &opts)
	if err != nil {
		saveDefaultConfig(defaultOpts)
		return defaultOpts, err
	}
	return opts, nil
}

func saveDefaultConfig(cfg convert.ConvertOptions) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	_ = os.WriteFile(configFilePath(), data, 0644)
}
