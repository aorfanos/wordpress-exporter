package utils

// import (
// 	"io/ioutil"

// 	"gopkg.in/yaml.v2"
// )

// type ConfigFile struct {
// 	MonitoredWordpress []string `yaml:"wordpress-exporter"`
// }

// type Settings struct {
// 	title          string `json:"title"`
// 	language       string `json:"language"`
// 	ping_status    string `json:"ping_status"`
// 	comment_status string `json:"comment_status"`
// }

// func (c *ConfigFile) ParseConf() *ConfigFile {
// 	yamlFile, err := ioutil.ReadFile(*configFile)
// 	errCheck(err)
// 	err = yaml.Unmarshal(yamlFile, c)
// 	errCheck(err)
// 	return c
// }