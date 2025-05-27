package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Port          string
	CacheDir      string `yaml:"cache_dir"`
	Cdn           []cdnConfig
	SocksProxyUrl string `yaml:"socks_proxy_url"`
}

type cdnConfig struct {
	Name string
	Url  string
}

var Config = &config{}

func InitConfig() error {
	configBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Read config failed, %v\n", err)
		return err
	}

	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		log.Fatalf("Decode config failed: %v\n", err)
	}

	return nil
}
