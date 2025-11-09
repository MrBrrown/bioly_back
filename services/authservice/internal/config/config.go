package config

import (
	"bioly/storage"
	"bioly/yamlconf"
	"log"
	"time"
)

type HTTP struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type JWT struct {
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
	Issuer        string        `yaml:"issuer"`
}

type Config struct {
	DBInfo storage.DbInfo `yaml:"auth_db"`
	HTTP   HTTP           `yaml:"http"`
	JWT    JWT            `yaml:"jwt"`
}

func (c *Config) SetDefaults() {
	if c.JWT.AccessTTL == 0 {
		c.JWT.AccessTTL = 15 * time.Minute
	}
	if c.JWT.RefreshTTL == 0 {
		c.JWT.RefreshTTL = 30 * 24 * time.Hour
	}
	if c.JWT.Issuer == "" {
		c.JWT.Issuer = "auth.bioly.local"
	}
}

func New(path string) *Config {
	cfg := &Config{}
	err := yamlconf.Load(path, cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfg.SetDefaults()

	return cfg
}
