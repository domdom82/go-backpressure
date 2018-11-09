package server

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Protocol string        `yaml:"protocol"`
	Port     string        `yaml:"port"`
	Bufsize  int           `yaml:"bufsize"`
	Delay    time.Duration `yaml:"delay"`
}

func (cfg *Config) NewServer() Server {

	var server Server

	switch cfg.Protocol {
	case "tcp":
		server = &TcpServer{config: cfg}
	case "websocket":
		server = &WsServer{config: cfg}
	default:
		panic("Expected protocol of type tcp or websocket.")
	}

	return server
}

func NewServerConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Overwrite port if given by ENV
	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port = envPort
	}

	return config, nil
}
