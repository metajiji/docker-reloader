package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Watch []WatchedFile `yaml:"watch"`
}

type WatchedFile struct {
	Path      string `yaml:"path"`
	Command   string `yaml:"command"`
	Container string `yaml:"container"`
}

func LoadConfig(path string) Config {
	cfgFd, err := os.Open(path)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	defer cfgFd.Close()

	var cfg Config
	decoder := yaml.NewDecoder(cfgFd)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	return cfg
}
