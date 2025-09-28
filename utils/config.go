package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	MonitoredWordpress []string `yaml:"wordpress-exporter"`
}

type Settings struct {
	title          string `json:"title"`
	language       string `json:"language"`
	ping_status    string `json:"ping_status"`
	comment_status string `json:"comment_status"`
}

func (c *ConfigFile) ParseConf(configFilePath string) (*ConfigFile, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func LoadConfig(configFilePath string) ([]string, error) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, err
	}

	config := &ConfigFile{}
	parsedConfig, err := config.ParseConf(configFilePath)
	if err != nil {
		return nil, err
	}

	return parsedConfig.MonitoredWordpress, nil
}
