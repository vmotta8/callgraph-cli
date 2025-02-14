package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

const ConfigFileName = ".refactorclirc"

type Config struct {
	Project ProjectConfig `ini:"project"`
	LLM     LLMConfig     `ini:"llm"`
}

type ProjectConfig struct {
	Language    string `ini:"language"`
	EntryFile   string `ini:"entry_file"`
	Endpoint    string `ini:"endpoint"`
	Method      string `ini:"method"`
	LLMProvider string `ini:"llm_provider"`
}

type LLMConfig struct {
	APIKey string `ini:"api_key"`
	Model  string `ini:"model"`
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)

	iniFile, err := ini.Load(ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v. Did you run 'init'?", err)
	}

	err = iniFile.MapTo(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return cfg, nil
}

func CreateLocalConfig() error {
	if _, err := os.Stat(ConfigFileName); err == nil {
		return nil
	}

	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func GenerateConfigFile(language string) error {
	configContent := `[project]
language = "` + language + `"
entry_file = ""
endpoint = ""
method = ""
llm_provider = ""

[llm]
api_key = ""
model = ""
`
	return os.WriteFile(ConfigFileName, []byte(configContent), 0644)
}
