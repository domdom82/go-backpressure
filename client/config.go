package client

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Protocol          string `yaml:"protocol"`
	ServerAddr        string `yaml:"server_addr"`
	ServerPort        string `yaml:"server_port"`
	PayloadSize       int    `yaml:"payload_size"`
	RequestsPerSecond int    `yaml:"requests_per_second"`
	RequestsTotal     int    `yaml:"requests_total"`
}

func (cfg *Config) NewClient() Client {
	var client Client

	switch cfg.Protocol {
	case "tcp":
		client = &TcpClient{config: cfg}
	case "websocket":
		client = &WsClient{config: cfg}
	default:
		panic("Expected protocol of type tcp or websocket.")
	}

	return client
}

func NewClientConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
