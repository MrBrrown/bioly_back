package config

import (
	"bioly/common/storage"
	"bioly/common/yamlconf"
	"log"
	"time"
)

type HTTP struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Config struct {
	DBInfo storage.DbInfo `yaml:"profile_db"`
	HTTP   HTTP           `yaml:"http"`
}

func New(path string) *Config {
	cfg := &Config{}
	err := yamlconf.Load(path, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
