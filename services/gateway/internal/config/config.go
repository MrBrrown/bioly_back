package config

import (
	"time"
)

type Config struct {
	Proxies []Proxy `yaml:"proxies"`
	Server  Server  `yaml:"server"`
}

type Proxy struct {
	Url      string `yaml:"url"`
	Endpoint string `yaml:"endpint"`
}

type Server struct {
	Address         string        `yaml:"address"`
	ReadTimeout     time.Duration `yaml:"readtimeout"`
	WriteTimeout    time.Duration `yaml:"writetimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdowntime"`
}
