package config

import (
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Rules             Rules           `yaml:"rules" json:"rules" env-required:"true"`
	AutoFix           AutoFix         `yaml:"auto_fix" json:"auto_fix" env-required:"true"`
	SensitiveKeywords []string        `yaml:"sensitive_keywords" json:"sensitive_keywords" env-required:"true"`
	CustomPatterns    []CustomPattern `yaml:"custom_patterns" json:"custom_patterns"`
}

type Rules struct {
	Lowercase         bool `yaml:"lowercase" json:"lowercase"`
	English           bool `yaml:"english" json:"english"`
	SpecialChars      bool `yaml:"special_chars" json:"special_chars"`
	SensitiveKeywords bool `yaml:"sensitive_keywords" json:"sensitive_keywords"`
	CustomPatterns    bool `yaml:"custom_patterns" json:"custom_patterns"`
}

type AutoFix struct {
	Enabled           bool `yaml:"enabled" json:"enabled"`
	Lowercase         bool `yaml:"lowercase" json:"lowercase"`
	SpecialChars      bool `yaml:"special_chars" json:"special_chars"`
	SensitiveKeywords bool `yaml:"sensitive_keywords" json:"sensitive_keywords"`
	CustomPatterns    bool `yaml:"custom_patterns" json:"custom_patterns"`
}

type CustomPattern struct {
	Name        string `yaml:"name" json:"name"`
	Pattern     string `yaml:"pattern" json:"pattern"`
	Message     string `yaml:"message" json:"message"`
	AutoFix     bool   `yaml:"auto_fix" json:"auto_fix"`
	Replacement string `yaml:"replacement" json:"replacement"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("no config file")
	}
	return MustLoadPath(configPath)
}

func MustLoadPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}
	if exePath, err := os.Executable(); err == nil {
		baseDir := filepath.Dir(filepath.Dir(exePath))
		configPath := filepath.Join(baseDir, "config", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return configPath
}
