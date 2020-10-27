package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

var config *Conf

func Config() *Conf {
	return config
}

func InitConfig(configFile string) error {
	c, err := loadFromFile(configFile)
	if err != nil {
		return err
	}
	config = c
	return nil
}

func loadFromFile(configFile string) (cfg *Conf, err error) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	cfg = &Conf{}
	err = yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		return nil, err
	}

	return
}

type Conf struct {
	ServerPort   string `json:"server_port"`
	WebAdminPort string `json:"web_admin_port"`
	DataDir      string `json:"data_dir"`
}
