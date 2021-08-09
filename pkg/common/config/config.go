package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	configFile string `json:"-"`
	Host       string `json:"host"`
	ServerPort string `json:"server_port"`
	SyncPort   string `json:"sync_port"`
	AdminPort  string `json:"admin_port"`
	DataDir    string `json:"data_dir"`
}

func NewConfig(configFile string) *Config {
	return &Config{
		configFile: configFile,
	}
}

func (c *Config) Init() error {
	yamlFile, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return errors.Wrap(err, "init config: read file error")
	}

	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return errors.Wrap(err, "init config: unmarshal error")
	}

	c.Print()

	return nil
}

func (c *Config) Run() error {
	return nil
}

func (c *Config) Shutdown() {
	return
}

func (c *Config) Print() {
	fmt.Println("config_file", c.configFile)
	fmt.Println("host:", c.Host)
	fmt.Println("server_port:", c.ServerPort)
	fmt.Println("admin_port:", c.AdminPort)
	fmt.Println("sync_port:", c.SyncPort)
	fmt.Println("data_dir:", c.DataDir)
}
